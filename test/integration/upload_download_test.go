// +build integration

package integration

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/yurykabanov/go-yandex-disk/internal/testhelpers"
)

func TestYaDiskClient_Upload_Download(t *testing.T) {
	fileContent := "SOME FILE CONTENT"

	client := YaDiskClient()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	randomName := "/" + testhelpers.RandString(50) + ".txt"

	uploadLink, err := client.RequestUploadLink(ctx, randomName, true)
	assert.Nil(t, err)
	assert.NotNil(t, uploadLink)
	assert.NotEmpty(t, uploadLink.Href)
	assert.Equal(t, http.MethodPut, uploadLink.Method)

	status, err := client.Upload(ctx, uploadLink, strings.NewReader(fileContent))
	assert.Nil(t, err)
	assert.True(t, 201 <= status && status <= 202)

	downloadLink, err := client.RequestDownloadLink(ctx, randomName)
	assert.Nil(t, err)
	assert.NotNil(t, downloadLink)
	assert.NotEmpty(t, downloadLink.Href)
	assert.Equal(t, http.MethodGet, downloadLink.Method)

	resp, err := client.Download(ctx, downloadLink)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, fileContent, string(body))

	deleteLink, status, err := client.Delete(ctx, randomName, true)
	assert.Nil(t, err)
	assert.Nil(t, deleteLink) // It should be nil, because it's a file
	assert.Equal(t, http.StatusNoContent, status)
}
