// Client creation elements
package vote

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

// HTTPRequestDoer performs HTTP requests.
// The standard http.Client implements this interface.
type HTTPRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
type Client struct {
	httpClient HTTPRequestDoer
	baseURL    string
	urlFormal  *url.URL
	username   string
	password   string
	userAgent  string
	cookieJar  *cookiejar.Jar
}

// ClientOption allows setting custom parameters during construction.
type ClientOption func(*Client) error

// NewClient creates a new AURWebClient.
func NewClient(opts ...ClientOption) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := Client{}

	for _, o := range opts {
		if o == nil {
			continue
		}
		if errOpt := o(&client); errOpt != nil {
			return nil, errOpt
		}
	}

	if client.httpClient == nil {
		httpClient := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}

		client.httpClient = httpClient
	}

	if client.baseURL == "" {
		client.baseURL = defaultURL
	}

	if client.userAgent == "" {
		client.userAgent = defaultUserAgent
	}

	client.cookieJar = jar
	client.urlFormal, err = url.Parse(client.baseURL)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HTTPRequestDoer) ClientOption {
	return func(c *Client) error {
		c.httpClient = doer

		return nil
	}
}

// WithBaseURL allows overriding the default base URL of the client.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		c.baseURL = baseURL

		return nil
	}
}

// WithUserAgent allows overriding the default user agent of the client.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		c.userAgent = userAgent

		return nil
	}
}
