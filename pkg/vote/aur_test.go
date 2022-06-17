package vote

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
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
	client, err := NewClient()
	require.NoError(t, err)

	err = client.Vote(context.Background(), "votar")
	require.Error(t, ErrNoCredentials, err)
}

type MockClient struct {
	t       *testing.T
	wantURL string
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	require.NotNil(m.t, req)
	assert.Equal(m.t, "POST", req.Method)
	assert.Equal(m.t, m.wantURL, req.URL.String())
	return &http.Response{Body: http.NoBody}, nil
}

func TestLoginRequest(t *testing.T) {
	client, err := NewClient(WithHTTPClient(&MockClient{
		t:       t,
		wantURL: "https://aur.archlinux.org/login?next=packages"}))
	require.NoError(t, err)

	client.SetCredentials("bob", "bob")

	err = client.login(context.Background())
	require.Error(t, ErrNoCredentials, err)
}
