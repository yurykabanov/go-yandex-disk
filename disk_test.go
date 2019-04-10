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

func TestClient_GetDisk(t *testing.T) {
	tests := []struct {
		name string

		responseStatusCode int
		responseBody       string

		response *Disk
		error    error
	}{
		{
			name: "successful request",

			responseStatusCode: 200,
			responseBody:       `{"total_space":12345}`, // NOTE: some fields are omitted

			response: &Disk{TotalSpace: 12345},
			error:    nil,
		},

		{
			name: "api error",

			responseStatusCode: 401,
			responseBody:       `{"message":"Не авторизован","description":"Unauthorized","error":"UnauthorizedError"}`,

			response: nil,
			error:    ApiError{StatusCode: 401, Message: "Не авторизован", Description: "Unauthorized", ErrorID: "UnauthorizedError"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := New(testhelpers.NewTestClient(func(req *http.Request) *http.Response {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "https://cloud-api.yandex.net/v1/disk/", req.URL.String())

				return &http.Response{
					StatusCode: test.responseStatusCode,
					Body:       ioutil.NopCloser(bytes.NewBufferString(test.responseBody)),
				}
			}))

			disk, err := client.GetDisk(context.Background())

			assert.Equal(t, test.response, disk)
			assert.Equal(t, test.error, err)
		})
	}
}
