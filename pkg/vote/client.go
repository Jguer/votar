package vote

import (
	"net/http"
)

// HTTPRequestDoer performs HTTP requests.
// The standard http.Client implements this interface.
type HTTPRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HTTPRequestDoer) ClientOption {
	return func(c *AURWebClient) error {
		c.httpClient = doer

		return nil
	}
}

// WithBaseURL allows overriding the default base URL of the client.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *AURWebClient) error {
		c.baseURL = baseURL

		return nil
	}
}

// WithUserAgent allows overriding the default user agent of the client.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *AURWebClient) error {
		c.userAgent = userAgent

		return nil
	}
}
