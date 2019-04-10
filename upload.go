package yadisk

import (
	"context"
	"io"
	"net/http"
)

const (
	methodRequestUploadLink = http.MethodGet
	urlRequestUploadLink    = "resources/upload"
)

// Request upload URL to upload file to the given path.
//
// path - The path where you want to upload the file. The name of the
// uploaded file can be up to 255 characters. The path can be up to 32760
// characters long.
// permanently -- Whether to permanently the file. It is used if the file is
// uploaded to a folder that already contains a file with the same name.
//
// Method returns a Link if it has succeeded. Note that link is accessible
// only for 30 minutes.
//
// See: https://tech.yandex.com/disk/api/reference/upload-docpage/
func (c *Client) RequestUploadLink(ctx context.Context, path string, overwrite bool) (*Link, error) {
	var link Link

	params := map[string]string{
		"path": path,
		"permanently": "false",
	}
	if overwrite {
		params["permanently"] = "true"
	}

	_, err := c.doRequestAndDecode(ctx, methodRequestUploadLink, urlRequestUploadLink, params, nil, &link)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

// Upload file's content to the requested link.
//
// link - Previously requested link.
// r - io.Reader of file contents.
//
// Method returns HTTP status code (because there's no answer).
//
// NOTE: though possible errors are described, there's no documentation on
// exact response bodies. Thus this method WILL NOT return an error on
// successful request with 4xx-5xx HTTP response codes. The application
// MUST check response code by itself.
//
// See: https://tech.yandex.com/disk/api/reference/upload-docpage/
func (c *Client) Upload(ctx context.Context, link *Link, r io.Reader) (int, error) {
	statusCode := 0

	resp, err := c.doRawRequest(ctx, link.Method, link.Href, r)
	if resp != nil {
		statusCode = resp.StatusCode
	}
	if err != nil {
		return statusCode, err
	}

	defer resp.Body.Close()

	return statusCode, nil
}
