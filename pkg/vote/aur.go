package vote

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
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

func NewClient() (*AURClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &AURClient{url: defaultURL, client: &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}}

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

func (a *AURClient) login() error {
	resp, err := a.client.PostForm(a.url+"/login/", url.Values{
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

func (a *AURClient) getToken(pkgbase string) (string, error) {
	resp, err := a.client.Get(fmt.Sprintf("%s/packages/%s", a.url, pkgbase))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Token status not OK")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	match := tokenExpr.FindStringSubmatch(string(bodyBytes))
	if match == nil {
		return "", errors.New("No match for token")
	}
	return match[1], nil
}

func (a *AURClient) Vote(pkgbase string) error {
	if len(a.client.Jar.Cookies(a.urlFormal)) == 0 {
		if err := a.login(); err != nil {
			return err
		}
	}
	token, err := a.getToken(pkgbase)
	if err != nil {
		return err
	}

	resp, err := a.client.PostForm(fmt.Sprintf("%s/pkgbase/%s/vote/", a.url, pkgbase), url.Values{
		"token":   []string{token},
		"do_Vote": []string{"Vote+for+this+package"},
	})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusFound {
		return errors.New("unable to vote")
	}

	resp.Body.Close()
	return nil
}

func (a *AURClient) Unvote(pkgbase string) error {
	if len(a.client.Jar.Cookies(a.urlFormal)) == 0 {
		if err := a.login(); err != nil {
			return err
		}
	}

	token, err := a.getToken(pkgbase)
	if err != nil {
		return err
	}

	resp, err := a.client.PostForm(fmt.Sprintf("%s/pkgbase/%s/unvote/", a.url, pkgbase), url.Values{
		"token":     []string{token},
		"do_UnVote": []string{"Remove+vote"},
	})
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusFound {
		return errors.New("unable to vote")
	}

	resp.Body.Close()
	return nil
}
