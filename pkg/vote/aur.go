// Methods interacting with the AUR website.
package vote

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultURL       = "https://aur.archlinux.org"
	defaultUserAgent = "votar/1.0.0"
)

// SetCredentials sets the username and password for the client.
func (a *Client) SetCredentials(username, password string) {
	a.username = username
	a.password = password
}

func (a *Client) login(ctx context.Context) error {
	if a.username == "" || a.password == "" {
		return ErrNoCredentials
	}

	loginURL := fmt.Sprintf("%s/login", a.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, strings.NewReader(url.Values{
		"user":        []string{a.username},
		"passwd":      []string{a.password},
		"referer":     []string{a.baseURL},
		"remember_me": []string{"on"},
		"next":        []string{"packages"},
	}.Encode()))
	if err != nil {
		return err
	}

	a.setHeaders(req, loginURL)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return &ErrLoginFailed{status: resp.StatusCode, body: string(bodyBytes)}
	}

	a.cookieJar.SetCookies(a.urlFormal, resp.Cookies())

	return nil
}

func (a *Client) setHeaders(req *http.Request, refererURL string) {
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", refererURL)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Origin", a.baseURL)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", a.userAgent)
	for _, cookie := range a.cookieJar.Cookies(a.urlFormal) {
		req.AddCookie(cookie)
	}
}

func (a *Client) handleVote(ctx context.Context, pkgbase string, vote bool) error {
	if len(a.cookieJar.Cookies(a.urlFormal)) == 0 {
		if err := a.login(ctx); err != nil {
			return err
		}
	}

	values := url.Values{}
	packageURL := fmt.Sprintf("%s/pkgbase/%s", a.baseURL, pkgbase)
	voteURL := ""
	if vote {
		values.Add("do_Vote", "Vote+for+this+package")
		voteURL = fmt.Sprintf("%s/vote", packageURL)
	} else {
		values.Add("do_UnVote", "Remove+vote")
		voteURL = fmt.Sprintf("%s/unvote", packageURL)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, voteURL, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}

	a.setHeaders(req, packageURL)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return &ErrVoteFailed{status: resp.StatusCode, cookies: a.cookieJar.Cookies(a.urlFormal), body: string(bodyBytes)}
	}

	return nil
}

// Vote votes for the given package.
func (a *Client) Vote(ctx context.Context, pkgbase string) error {
	return a.handleVote(ctx, pkgbase, true)
}

// Unvote unvotes for the given package.
func (a *Client) Unvote(ctx context.Context, pkgbase string) error {
	return a.handleVote(ctx, pkgbase, false)
}
