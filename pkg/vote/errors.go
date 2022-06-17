// Errors
package vote

import (
	"fmt"
	"net/http"
)

var ErrNoCredentials = fmt.Errorf("no credentials provided")

type ErrLoginFailed struct {
	status int
	body   string
}

func (e *ErrLoginFailed) Error() string {
	return fmt.Sprintf("login failed with status %d. body: %s", e.status, e.body)
}

type ErrVoteFailed struct {
	status  int
	cookies []*http.Cookie
	body    string
}

func (e *ErrVoteFailed) Error() string {
	cookieString := ""
	for _, c := range e.cookies {
		cookieString += fmt.Sprintf("%s=%s; ", c.Name, c.Value)
	}

	return fmt.Sprintf("body: %s. Vote failed with status %d. Cookies: %s. ", e.body, e.status, cookieString)
}
