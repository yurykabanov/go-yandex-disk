package yadisk

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yurykabanov/go-yandex-disk/internal/testhelpers"
)

func TestClient_Copy(t *testing.T) {
	tests := []struct {
		name string

		responseStatusCode int
		responseBody       string

		src       string
		dst       string
		overwrite bool

		response *Link
		error    error
	}{
		{
			name: "successfully copied and processed",

			responseStatusCode: 201,
			responseBody:       `{"href":"some_href","method":"GET","templated":false}`,

			src:       "/source/some_file.ext",
			dst:       "/destination/some_file.ext",
			overwrite: false,

			response: &Link{
				Href:      "some_href",
				Method:    "GET",
				Templated: false,
			},

			error: nil,
		},

		{
			name: "successfully copied but not processed",

			responseStatusCode: 202,
			responseBody:       `{"href":"some_href","method":"GET","templated":false}`,

			src:       "/source/some_file.ext",
			dst:       "/destination/some_file.ext",
			overwrite: false,

			response: &Link{
				Href:      "some_href",
				Method:    "GET",
				Templated: false,
			},

			error: nil,
		},

		{
			name: "successfully copied with permanently",

			responseStatusCode: 201,
			responseBody:       `{"href":"some_href","method":"GET","templated":false}`,

			src:       "/source/some_file.ext",
			dst:       "/destination/some_file.ext",
			overwrite: true,

			response: &Link{
				Href:      "some_href",
				Method:    "GET",
				Templated: false,
			},

			error: nil,
		},

		{
			name: "error while copying",

			responseStatusCode: 409,
			responseBody:       `{"message":"Ресурс {path} уже существует","description":"Resource already exists","error":"ResourcesAlreadyExistsError"}`, // NOTE: could differ from actual response

			src:       "/source/some_file.ext",
			dst:       "/destination/some_file.ext",
			overwrite: false,

			response: nil,

			error: ApiError{StatusCode: 409, Message: "Ресурс {path} уже существует", Description: "Resource already exists", ErrorID: "ResourcesAlreadyExistsError"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedUrl := testhelpers.BuildUrl("https://cloud-api.yandex.net/v1/disk/resources/copy", map[string]string{
				"from":        test.src,
				"path":        test.dst,
				"permanently": fmt.Sprintf("%t", test.overwrite),
			})

			client := New(testhelpers.NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, expectedUrl, req.URL.String())

				return &http.Response{
					StatusCode: test.responseStatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(test.responseBody)),
				}
			}))

			link, statusCode, err := client.Copy(context.Background(), test.src, test.dst, test.overwrite)

			assert.Equal(t, test.response, link)
			assert.Equal(t, test.responseStatusCode, statusCode)
			assert.Equal(t, test.error, err)
		})
	}
}

func TestClient_Move(t *testing.T) {
	tests := []struct {
		name string

		responseStatusCode int
		responseBody       string

		src       string
		dst       string
		overwrite bool

		response *Link
		error    error
	}{
		{
			name: "successfully moved and processed",

			responseStatusCode: 201,
			responseBody:       `{"href":"some_href","method":"GET","templated":false}`,

			src:       "/source/some_file.ext",
			dst:       "/destination/some_file.ext",
			overwrite: false,

			response: &Link{
				Href:      "some_href",
				Method:    "GET",
				Templated: false,
			},

			error: nil,
		},

		{
			name: "successfully moved but not processed",

			responseStatusCode: 202,
			responseBody:       `{"href":"some_href","method":"GET","templated":false}`,

			src:       "/source/some_file.ext",
			dst:       "/destination/some_file.ext",
			overwrite: false,

			response: &Link{
				Href:      "some_href",
				Method:    "GET",
				Templated: false,
			},

			error: nil,
		},

		{
			name: "successfully moved with permanently",

			responseStatusCode: 201,
			responseBody:       `{"href":"some_href","method":"GET","templated":false}`,

			src:       "/source/some_file.ext",
			dst:       "/destination/some_file.ext",
			overwrite: true,

			response: &Link{
				Href:      "some_href",
				Method:    "GET",
				Templated: false,
			},

			error: nil,
		},

		{
			name: "error while moving",

			responseStatusCode: 409,
			responseBody:       `{"message":"Ресурс {path} уже существует","description":"Resource already exists","error":"ResourcesAlreadyExistsError"}`, // NOTE: could differ from actual response

			src:       "/source/some_file.ext",
			dst:       "/destination/some_file.ext",
			overwrite: false,

			response: nil,

			error: ApiError{StatusCode: 409, Message: "Ресурс {path} уже существует", Description: "Resource already exists", ErrorID: "ResourcesAlreadyExistsError"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedUrl := testhelpers.BuildUrl("https://cloud-api.yandex.net/v1/disk/resources/move", map[string]string{
				"from":        test.src,
				"path":        test.dst,
				"permanently": fmt.Sprintf("%t", test.overwrite),
			})

			client := New(testhelpers.NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, http.MethodPost, req.Method)
				assert.Equal(t, expectedUrl, req.URL.String())

				return &http.Response{
					StatusCode: test.responseStatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(test.responseBody)),
				}
			}))

			link, statusCode, err := client.Move(context.Background(), test.src, test.dst, test.overwrite)

			assert.Equal(t, test.response, link)
			assert.Equal(t, test.responseStatusCode, statusCode)
			assert.Equal(t, test.error, err)
		})
	}
}

func TestClient_Delete(t *testing.T) {
	tests := []struct {
		name string

		responseStatusCode int
		responseBody       string

		path        string
		permanently bool

		response *Link
		error    error
	}{
		{
			name: "successfully deleted and processed",

			responseStatusCode: 204,
			responseBody:       ``,

			path:        "/some_path/some_file.ext",
			permanently: false,

			response: nil,

			error: nil,
		},

		{
			name: "successfully deleted but not processed",

			responseStatusCode: 202,
			responseBody:       `{"href":"some_href","method":"GET","templated":false}`,

			path:        "/some_path/some_file.ext",
			permanently: false,

			response: &Link{
				Href:      "some_href",
				Method:    "GET",
				Templated: false,
			},

			error: nil,
		},

		{
			name: "successfully deleted permanently",

			responseStatusCode: 204,
			responseBody:       ``,

			path:        "/some_path/some_file.ext",
			permanently: true,

			response: nil,

			error: nil,
		},

		{
			name: "error while deleting",

			responseStatusCode: 409,
			responseBody:       `{"message":"Указанного пути {path} не существует.","description":"Resource doesn't exist","error":"ResourcesDoesntExistError"}`, // NOTE: could differ from actual response

			path:        "/some_path/some_file.ext",
			permanently: false,

			response: nil,

			error: ApiError{StatusCode: 409, Message: "Указанного пути {path} не существует.", Description: "Resource doesn't exist", ErrorID: "ResourcesDoesntExistError"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedUrl := testhelpers.BuildUrl("https://cloud-api.yandex.net/v1/disk/resources", map[string]string{
				"path":        test.path,
				"permanently": fmt.Sprintf("%t", test.permanently),
			})

			client := New(testhelpers.NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, http.MethodDelete, req.Method)
				assert.Equal(t, expectedUrl, req.URL.String())

				return &http.Response{
					StatusCode: test.responseStatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(test.responseBody)),
				}
			}))

			link, statusCode, err := client.Delete(context.Background(), test.path, test.permanently)

			assert.Equal(t, test.response, link)
			assert.Equal(t, test.responseStatusCode, statusCode)
			assert.Equal(t, test.error, err)
		})
	}
}

func TestClient_CreateDirectory(t *testing.T) {
	tests := []struct {
		name string

		responseStatusCode int
		responseBody       string

		path string

		response *Link
		error    error
	}{
		{
			name: "successfully created",

			responseStatusCode: 201,
			responseBody:       `{"href":"some_href","method":"GET","templated":false}`,

			path: "/some_path/some_directory",

			response: &Link{
				Href:      "some_href",
				Method:    "GET",
				Templated: false,
			},

			error: nil,
		},

		{
			name: "error while creating directory",

			responseStatusCode: 409,
			responseBody:       `{"message":"Указанного пути {path} не существует.","description":"Resource doesn't exist","error":"ResourcesDoesntExistError"}`, // NOTE: could differ from actual response

			path: "/some_path/some_directory",

			response: nil,

			error: ApiError{StatusCode: 409, Message: "Указанного пути {path} не существует.", Description: "Resource doesn't exist", ErrorID: "ResourcesDoesntExistError"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedUrl := testhelpers.BuildUrl("https://cloud-api.yandex.net/v1/disk/resources", map[string]string{
				"path": test.path,
			})

			client := New(testhelpers.NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, http.MethodPut, req.Method)
				assert.Equal(t, expectedUrl, req.URL.String())

				return &http.Response{
					StatusCode: test.responseStatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(test.responseBody)),
				}
			}))

			link, err := client.CreateDirectory(context.Background(), test.path)

			assert.Equal(t, test.response, link)
			assert.Equal(t, test.error, err)
		})
	}
}
