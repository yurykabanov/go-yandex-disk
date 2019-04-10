package yadisk

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yurykabanov/go-yandex-disk/internal/testhelpers"
)

func TestClient_RequestDownloadLink(t *testing.T) {
	tests := []struct {
		name string

		responseStatusCode int
		responseBody       string

		path      string

		response *Link
		error    error
	}{
		{
			name: "successfully created upload link without permanently",

			responseStatusCode: 200,
			responseBody:       `{"href":"some_href","method":"PUT","templated":false}`,

			path: "/some_path/some_file.ext",

			response: &Link{
				Href:      "some_href",
				Method:    "PUT",
				Templated: false,
			},

			error: nil,
		},

		{
			name: "error resource not found",

			responseStatusCode: 404,
			responseBody:       `{"message":"Не удалось найти запрошенный ресурс.","description":"Resource not found","error":"ResourcesNotFoundError"}`, // NOTE: could differ from actual response

			path:      "/some_path/some_file.ext",

			response: nil,

			error: ApiError{StatusCode: 404, Message: "Не удалось найти запрошенный ресурс.", Description: "Resource not found", ErrorID: "ResourcesNotFoundError"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedUrl := testhelpers.BuildUrl("https://cloud-api.yandex.net/v1/disk/resources/download", map[string]string{
				"path":      test.path,
			})

			client := New(testhelpers.NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, expectedUrl, req.URL.String())

				return &http.Response{
					StatusCode: test.responseStatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(test.responseBody)),
				}
			}))

			link, err := client.RequestDownloadLink(context.Background(), test.path)

			assert.Equal(t, test.response, link)
			assert.Equal(t, test.error, err)
		})
	}
}

func TestClient_Download(t *testing.T) {
	tests := []struct {
		name string

		link Link

		responseStatusCode int
		responseBody       string

		statusCode int
		error      error
	}{
		{
			name: "successfully uploaded and processed",

			link: Link{
				Href:      "https://downloader.yandex.net/disk/some_href",
				Method:    "GET",
				Templated: false,
			},

			responseStatusCode: 200,
			responseBody:       `SOME FILE CONTENT`,

			statusCode: 200,
			error:      nil,
		},

		{
			name: "error while downloading",

			link: Link{
				Href:      "https://downloader.yandex.net/disk/some_href",
				Method:    "GET",
				Templated: false,
			},

			responseStatusCode: 500,
			responseBody:       ``,

			statusCode: 500,
			error:      nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := New(testhelpers.NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, test.link.Method, req.Method)
				assert.Equal(t, test.link.Href, req.URL.String())

				return &http.Response{
					StatusCode: test.responseStatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(test.responseBody)),
				}
			}))

			resp, err := client.Download(context.Background(), &test.link)
			assert.Equal(t, test.error, err)

			body, _ := ioutil.ReadAll(resp.Body)
			assert.Equal(t, test.responseBody, string(body))
		})
	}
}
