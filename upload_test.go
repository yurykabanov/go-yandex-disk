package yadisk

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yurykabanov/go-yandex-disk/internal/testhelpers"
)

func TestClient_RequestUploadLink(t *testing.T) {
	tests := []struct {
		name string

		responseStatusCode int
		responseBody       string

		path      string
		overwrite bool

		response *Link
		error    error
	}{
		{
			name: "successfully created upload link without permanently",

			responseStatusCode: 200,
			responseBody:       `{"href":"some_href","method":"PUT","templated":false}`,

			path:      "/some_path/some_file.ext",
			overwrite: false,

			response: &Link{
				Href:      "some_href",
				Method:    "PUT",
				Templated: false,
			},

			error: nil,
		},

		{
			name: "error path exists without permanently",

			responseStatusCode: 409,
			responseBody:       `{"message":"Ресурс {path} уже существует","description":"Resource already exists","error":"ResourcesAlreadyExistsError"}`, // NOTE: could differ from actual response

			path:      "/some_path/some_file.ext",
			overwrite: false,

			response: nil,

			error: ApiError{StatusCode: 409, Message: "Ресурс {path} уже существует", Description: "Resource already exists", ErrorID: "ResourcesAlreadyExistsError"},
		},

		{
			name: "successfully created upload link with permanently",

			responseStatusCode: 200,
			responseBody:       `{"href":"some_href","method":"PUT","templated":false}`,

			path:      "/some_path/some_file_to_overwrite.ext",
			overwrite: true,

			response: &Link{
				Href:      "some_href",
				Method:    "PUT",
				Templated: false,
			},

			error: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedUrl := testhelpers.BuildUrl("https://cloud-api.yandex.net/v1/disk/resources/upload", map[string]string{
				"path":      test.path,
				"permanently": fmt.Sprintf("%t", test.overwrite),
			})

			client := New(testhelpers.NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, expectedUrl, req.URL.String())

				return &http.Response{
					StatusCode: test.responseStatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(test.responseBody)),
				}
			}))

			link, err := client.RequestUploadLink(context.Background(), test.path, test.overwrite)

			assert.Equal(t, test.response, link)
			assert.Equal(t, test.error, err)
		})
	}
}

func TestClient_Upload(t *testing.T) {
	tests := []struct {
		name string

		link Link
		body string

		responseStatusCode int
		responseBody       string

		statusCode int
		error      error
	}{
		{
			name: "successfully uploaded and processed",

			link: Link{
				Href:      "https://uploader.yandex.net/upload-target/some_href",
				Method:    "PUT",
				Templated: false,
			},

			responseStatusCode: 201,
			responseBody:       ``,

			statusCode: 201,
			error:      nil,
		},

		{
			name: "successfully uploaded but not processed",

			link: Link{
				Href:      "https://uploader.yandex.net/upload-target/some_href",
				Method:    "PUT",
				Templated: false,
			},

			responseStatusCode: 202,
			responseBody:       ``,

			statusCode: 202,
			error:      nil,
		},

		{
			name: "error while uploading",

			link: Link{
				Href:      "https://uploader.yandex.net/upload-target/some_href",
				Method:    "PUT",
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
			bodyReader := strings.NewReader(test.body)

			client := New(testhelpers.NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, test.link.Method, req.Method)
				assert.Equal(t, test.link.Href, req.URL.String())

				body, err := ioutil.ReadAll(req.Body)
				assert.Nil(t, err)
				assert.Equal(t, test.body, string(body))

				return &http.Response{
					StatusCode: test.responseStatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(test.responseBody)),
				}
			}))

			statusCode, err := client.Upload(context.Background(), &test.link, bodyReader)

			assert.Equal(t, test.statusCode, statusCode)
			assert.Equal(t, test.error, err)
		})
	}
}
