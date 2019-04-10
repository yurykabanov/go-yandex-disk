package testhelpers

import (
	"net/url"
)

func BuildUrl(base string, params map[string]string) string {
	u, _ := url.Parse(base)

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
