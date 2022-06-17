package vote

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const cookieString = "AURSID=sidexample; HttpOnly; Max-Age=2592000; Path=/; SameSite=strict; Secure"

func TestNewClient(t *testing.T) {
	t.Parallel()
	type args struct {
		httpClient *http.Client
		baseURL    *string
	}

	sampleClient := &http.Client{}
	sampleURL := "http://azert.y"
	sampleUserAgent := "testuseragent"

	tests := []struct {
		name    string
		opts    []ClientOption
		wantErr bool
	}{
		{
			name:    "all nil",
			opts:    []ClientOption{},
			wantErr: false,
		},
		{
			name: "all set",
			opts: []ClientOption{
				WithBaseURL(sampleURL),
				WithHTTPClient(sampleClient),
				WithUserAgent(sampleUserAgent),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.opts...)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.NotEmpty(t, got)
			assert.NotNil(t, got.httpClient)
			assert.NotEmpty(t, got.urlFormal)
			assert.NotEmpty(t, got.userAgent)
			assert.NotEmpty(t, got.cookieJar)

			if len(tt.opts) != 0 {
				assert.Equal(t, sampleURL, got.baseURL)
				assert.Equal(t, sampleClient, got.httpClient)
			} else {
				assert.Equal(t, defaultURL, got.baseURL)
				assert.Equal(t, defaultURL, got.urlFormal.String())
			}
		})
	}
}

func TestMissingCredentials(t *testing.T) {
	t.Parallel()
	client, err := NewClient()
	require.NoError(t, err)

	err = client.Vote(context.Background(), "votar")
	require.Error(t, ErrNoCredentials, err)
}

type MockClient struct {
	t            *testing.T
	wantURL      string
	wantForm     url.Values
	giveResponse *http.Response
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	require.NotNil(m.t, req)
	assert.Equal(m.t, "POST", req.Method)
	assert.Equal(m.t, m.wantURL, req.URL.String())

	//req body to form
	err := req.ParseForm()
	require.NoError(m.t, err)
	assert.Equal(m.t, m.wantForm, req.Form)

	return m.giveResponse, nil
}

func TestLoginRequest(t *testing.T) {
	t.Parallel()
	client, err := NewClient(WithHTTPClient(&MockClient{
		t:       t,
		wantURL: "https://aur.archlinux.org/login",
		wantForm: url.Values{
			"user":        {"bob"},
			"passwd":      {"bob"},
			"next":        {"packages"},
			"referer":     {"https://aur.archlinux.org"},
			"remember_me": {"on"},
		},
		giveResponse: &http.Response{Body: http.NoBody, Status: "307 OK", StatusCode: 307, Header: http.Header{
			"Set-Cookie": []string{cookieString},
		}},
	}))
	require.NoError(t, err)

	client.SetCredentials("bob", "bob")

	err = client.login(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "AURSID=sidexample", client.cookieJar.Cookies(client.urlFormal)[0].String())
}

func TestVote(t *testing.T) {
	t.Parallel()
	client, err := NewClient(WithHTTPClient(&MockClient{
		t:       t,
		wantURL: "https://aur.archlinux.org/pkgbase/votar/vote",
		wantForm: url.Values{
			"do_Vote": {"Vote+for+this+package"},
		},
		giveResponse: &http.Response{Body: http.NoBody, Status: "303 OK", StatusCode: 303},
	}))
	require.NoError(t, err)

	client.SetCredentials("bob", "bob")

	client.cookieJar.SetCookies(client.urlFormal, []*http.Cookie{
		{Name: "AURSID", Value: "sidexample", Path: "/", Expires: time.Now().Add(time.Hour), MaxAge: 3600}})

	assert.Equal(t, "AURSID=sidexample", client.cookieJar.Cookies(client.urlFormal)[0].String())
	err = client.Vote(context.Background(), "votar")
	require.NoError(t, err)
}

func TestUnvote(t *testing.T) {
	t.Parallel()
	client, err := NewClient(WithHTTPClient(&MockClient{
		t:       t,
		wantURL: "https://aur.archlinux.org/pkgbase/votar/unvote",
		wantForm: url.Values{
			"do_UnVote": {"Remove+vote"},
		},
		giveResponse: &http.Response{Body: http.NoBody, Status: "303 OK", StatusCode: 303},
	}))
	require.NoError(t, err)

	client.SetCredentials("bob", "bob")

	client.cookieJar.SetCookies(client.urlFormal, []*http.Cookie{
		{Name: "AURSID", Value: "sidexample", Path: "/", Expires: time.Now().Add(time.Hour), MaxAge: 3600}})

	assert.Equal(t, "AURSID=sidexample", client.cookieJar.Cookies(client.urlFormal)[0].String())
	err = client.Unvote(context.Background(), "votar")
	require.NoError(t, err)
}
