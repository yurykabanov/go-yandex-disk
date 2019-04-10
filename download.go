package yadisk

import (
	"context"
	"net/http"
)

const (
	methodRequestDownloadLink = http.MethodGet
	urlRequestDownloadLink    = "resources/download"
)

// Request download URL for file with given path.
//
// path - The path to the file to download.
//
// Method returns a Link if it has succeeded.
//
// See: https://tech.yandex.com/disk/api/reference/content-docpage/
func (c *Client) RequestDownloadLink(ctx context.Context, path string) (*Link, error) {
	var link Link

	params := map[string]string{
		"path": path,
	}

	_, err := c.doRequestAndDecode(ctx, methodRequestDownloadLink, urlRequestDownloadLink, params, nil, &link)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

// Download file from given link.
//
// link - Previously requested link.
//
// NOTE: this method lacks proper documentation, possible errors are not
// described. It WILL NOT return an error on successful request with 4xx-5xx
// HTTP response codes. The application MUST check response code by itself.
//
// See: https://tech.yandex.com/disk/api/reference/content-docpage/
func (c *Client) Download(ctx context.Context, link *Link) (*http.Response, error) {
	return c.doRawRequest(ctx, link.Method, link.Href, nil)
}
