package yadisk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yurykabanov/go-yandex-disk/internal/testhelpers"
)

func TestClient_doRequest(t *testing.T) {
	errorStringType := fmt.Errorf("")

	tests := []struct {
		name string

		method       string
		url          string
		expectedUrl  string
		requestBody  interface{}
		responseBody []byte

		isError   bool
		errorType interface{}
	}{
		{
			name: "successful GET",

			method:       http.MethodGet,
			url:          "whatever/url/part",
			expectedUrl:  "https://cloud-api.yandex.net/v1/disk/whatever/url/part",
			requestBody:  nil,
			responseBody: []byte("whatever"),

			isError:   false,
			errorType: nil,
		},

		{
			name: "successful POST",

			method:       http.MethodPost,
			url:          "whatever/url/part",
			expectedUrl:  "https://cloud-api.yandex.net/v1/disk/whatever/url/part",
			requestBody:  struct{ One, Two string }{One: "111", Two: "222"},
			responseBody: []byte{},

			isError:   false,
			errorType: nil,
		},

		{
			name: "bad relative URL",

			method:       http.MethodGet,
			url:          "%",
			expectedUrl:  "",
			requestBody:  nil,
			responseBody: nil,

			isError:   true,
			errorType: &url.Error{},
		},

		{
			name: "bad POST object",

			method:       http.MethodPost,
			url:          "whatever/url/part",
			expectedUrl:  "https://cloud-api.yandex.net/v1/disk/whatever/url/part",
			requestBody:  make(chan struct{}),
			responseBody: []byte{},

			isError:   true,
			errorType: &json.UnsupportedTypeError{},
		},

		{
			name: "unable to create request",

			method:       "BAD METHOD",
			url:          "whatever/url/part",
			expectedUrl:  "",
			requestBody:  nil,
			responseBody: nil,

			isError:   true,
			errorType: errorStringType,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := New(testhelpers.NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, test.method, req.Method)
				assert.Equal(t, test.expectedUrl, req.URL.String())

				assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
				assert.Equal(t, "application/json", req.Header.Get("Accept"))

				if test.requestBody != nil {
					expectedBody, _ := json.Marshal(&test.requestBody)
					body, _ := ioutil.ReadAll(req.Body)
					assert.Equal(t, string(expectedBody)+"\n", string(body))
				}

				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBuffer(test.responseBody)),
				}
			}))

			resp, err := client.doRequest(context.Background(), test.method, test.url, nil, test.requestBody)
			assert.Equal(t, test.isError, err != nil)
			assert.IsType(t, test.errorType, err)

			if err == nil && test.requestBody != nil {
				responseBody, _ := ioutil.ReadAll(resp.Body)
				resp.Body.Close()

				assert.Equal(t, test.responseBody, responseBody)
			}
		})
	}
}

func TestClient_doRequestAndDecode(t *testing.T) {
	tests := []struct {
		name string

		method       string
		url          string
		expectedUrl  string
		requestBody  interface{}
		responseBody []byte

		responseStatusCode int

		isError   bool
		errorType interface{}
	}{
		{
			name: "successful request",

			method:       http.MethodGet,
			url:          "whatever/url/part",
			expectedUrl:  "https://cloud-api.yandex.net/v1/disk/whatever/url/part",
			requestBody:  nil,
			responseBody: []byte("{}"),

			responseStatusCode: 200,

			isError:   false,
			errorType: nil,
		},

		{
			name: "successful request with bad body",

			method:       http.MethodGet,
			url:          "whatever/url/part",
			expectedUrl:  "https://cloud-api.yandex.net/v1/disk/whatever/url/part",
			requestBody:  nil,
			responseBody: []byte("invalid json"),

			responseStatusCode: 200,

			isError:   true,
			errorType: &json.SyntaxError{},
		},

		{
			name: "error response",

			method:       http.MethodGet,
			url:          "whatever/url/part",
			expectedUrl:  "https://cloud-api.yandex.net/v1/disk/whatever/url/part",
			requestBody:  nil,
			responseBody: []byte("{}"),

			responseStatusCode: 429,

			isError:   true,
			errorType: ApiError{},
		},

		{
			name: "error response with bad body",

			method:       http.MethodGet,
			url:          "whatever/url/part",
			expectedUrl:  "https://cloud-api.yandex.net/v1/disk/whatever/url/part",
			requestBody:  nil,
			responseBody: []byte("invalid json"),

			responseStatusCode: 429,

			isError:   true,
			errorType: &json.SyntaxError{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := New(testhelpers.NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, test.method, req.Method)
				assert.Equal(t, test.expectedUrl, req.URL.String())

				assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
				assert.Equal(t, "application/json", req.Header.Get("Accept"))

				if test.requestBody != nil {
					expectedBody, _ := json.Marshal(&test.requestBody)
					body, _ := ioutil.ReadAll(req.Body)
					assert.Equal(t, string(expectedBody)+"\n", string(body))
				}

				return &http.Response{
					StatusCode: test.responseStatusCode,
					Body:       ioutil.NopCloser(bytes.NewBuffer(test.responseBody)),
				}
			}))

			var resp interface{}

			_, err := client.doRequestAndDecode(context.Background(), test.method, test.url, nil, test.requestBody, &resp)
			assert.Equal(t, test.isError, err != nil)
			assert.IsType(t, test.errorType, err)
		})
	}
}
