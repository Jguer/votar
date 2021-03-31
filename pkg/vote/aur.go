package vote

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"

	"golang.org/x/net/context/ctxhttp"
)

const defaultURL = "https://aur.archlinux.org"

var tokenExpr = regexp.MustCompile(`<input type="hidden" name="token"\s+value="([0-9a-f]+)" />`)

type AURClient struct {
	client    *http.Client
	url       string
	urlFormal *url.URL
	username  string
	password  string
}

func NewClient(httpClient *http.Client, baseURL *string) (*AURClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Jar = jar
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	client := &AURClient{url: defaultURL, client: httpClient}

	if baseURL != nil {
		client.url = *baseURL
	}

	client.urlFormal, err = url.Parse(client.url)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (a *AURClient) SetCredentials(username, password string) {
	a.username = username
	a.password = password
}

func (a *AURClient) login(ctx context.Context) error {
	resp, err := ctxhttp.PostForm(ctx, a.client, a.url+"/login/", url.Values{
		"user":        []string{a.username},
		"passwd":      []string{a.password},
		"remember_me": []string{"on"},
	})
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	a.client.Jar.SetCookies(a.urlFormal, resp.Cookies())

	return nil
}

func (a *AURClient) getToken(ctx context.Context, pkgbase string) (string, error) {
	resp, err := ctxhttp.Get(ctx, a.client, fmt.Sprintf("%s/packages/%s", a.url, pkgbase))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("token status not OK")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	match := tokenExpr.FindStringSubmatch(string(bodyBytes))
	if match == nil {
		return "", errors.New("no match for token")
	}
	return match[1], nil
}

func (a *AURClient) handleVote(ctx context.Context, pkgbase string, vote bool) error {
	if len(a.client.Jar.Cookies(a.urlFormal)) == 0 {
		if err := a.login(ctx); err != nil {
			return err
		}
	}

	token, err := a.getToken(ctx, pkgbase)
	if err != nil {
		return err
	}

	values := url.Values{
		"token": []string{token},
	}

	voteURL := ""
	if vote {
		values.Add("do_Vote", "Vote+for+this+package")
		voteURL = fmt.Sprintf("%s/pkgbase/%s/vote/", a.url, pkgbase)
	} else {
		values.Add("do_Vote", "Remove+vote")
		voteURL = fmt.Sprintf("%s/pkgbase/%s/unvote/", a.url, pkgbase)
	}

	resp, err := ctxhttp.PostForm(ctx, a.client, voteURL, values)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusFound {
		return errors.New("unable to vote")
	}

	resp.Body.Close()
	return nil
}

func (a *AURClient) Vote(ctx context.Context, pkgbase string) error {
	return a.handleVote(ctx, pkgbase, true)
}

func (a *AURClient) Unvote(ctx context.Context, pkgbase string) error {
	return a.handleVote(ctx, pkgbase, false)
}
