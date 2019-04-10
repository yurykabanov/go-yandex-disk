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

// Ensure file is copied by comparing its contents with original and by
// performing successful delete operation for copy and original
func TestYaDiskClient_Copy(t *testing.T) {
	fileContent := "SOME FILE CONTENT"

	client := YaDiskClient()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	originalName := "/" + testhelpers.RandString(50) + ".txt"
	copyName := "/" + testhelpers.RandString(50) + ".txt"

	uploadLink, err := client.RequestUploadLink(ctx, originalName, true)
	assert.Nil(t, err)

	_, err = client.Upload(ctx, uploadLink, strings.NewReader(fileContent))
	assert.Nil(t, err)

	link, status, err := client.Copy(ctx, originalName, copyName, true)
	assert.Nil(t, err)
	assert.NotNil(t, link)
	assert.NotEmpty(t, link.Href)
	assert.Equal(t, http.MethodGet, link.Method)
	assert.Equal(t, http.StatusCreated, status)

	downloadLink, err := client.RequestDownloadLink(ctx, copyName)
	assert.Nil(t, err)
	assert.NotNil(t, downloadLink)

	resp, err := client.Download(ctx, downloadLink)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, fileContent, string(body))

	deleteOriginalLink, status, err := client.Delete(ctx, originalName, true)
	assert.Nil(t, err)
	assert.Nil(t, deleteOriginalLink)
	assert.Equal(t, http.StatusNoContent, status)

	deleteCopyLink, status, err := client.Delete(ctx, copyName, true)
	assert.Nil(t, err)
	assert.Nil(t, deleteCopyLink)
	assert.Equal(t, http.StatusNoContent, status)
}

// Ensure file is moved by comparing its contents with original and by
// performing successful delete operation for moved file and 404 NotFound
// for original file
func TestYaDiskClient_Move(t *testing.T) {
	fileContent := "SOME FILE CONTENT"

	client := YaDiskClient()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	originalName := "/" + testhelpers.RandString(50) + ".txt"
	movedName := "/" + testhelpers.RandString(50) + ".txt"

	uploadLink, err := client.RequestUploadLink(ctx, originalName, true)
	assert.Nil(t, err)

	_, err = client.Upload(ctx, uploadLink, strings.NewReader(fileContent))
	assert.Nil(t, err)

	link, status, err := client.Move(ctx, originalName, movedName, true)
	assert.Nil(t, err)
	assert.NotNil(t, link)
	assert.NotEmpty(t, link.Href)
	assert.Equal(t, http.MethodGet, link.Method)
	assert.Equal(t, http.StatusCreated, status)

	downloadLink, err := client.RequestDownloadLink(ctx, movedName)
	assert.Nil(t, err)
	assert.NotNil(t, downloadLink)

	resp, err := client.Download(ctx, downloadLink)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, fileContent, string(body))

	deleteOriginalLink, status, err := client.Delete(ctx, originalName, true)
	assert.NotNil(t, err)
	assert.Nil(t, deleteOriginalLink)
	assert.Equal(t, http.StatusNotFound, status)

	deleteCopyLink, status, err := client.Delete(ctx, movedName, true)
	assert.Nil(t, err)
	assert.Nil(t, deleteCopyLink)
	assert.Equal(t, http.StatusNoContent, status)
}

// Ensure directory is created by successful delete operation for it
func TestYaDiskClient_CreateDirectory(t *testing.T) {
	client := YaDiskClient()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	directoryName := "/" + testhelpers.RandString(50)

	link, err := client.CreateDirectory(ctx, directoryName)
	assert.Nil(t, err)
	assert.NotNil(t, link)
	assert.NotEmpty(t, link.Href)
	assert.Equal(t, http.MethodGet, link.Method)

	deleteOriginalLink, status, err := client.Delete(ctx, directoryName, true)
	assert.Nil(t, err)
	assert.Nil(t, deleteOriginalLink)
	assert.Equal(t, http.StatusNoContent, status)
}

// Ensure file is deleted by double delete operations: success on first,
// not found on second
func TestYaDiskClient_Delete(t *testing.T) {
	fileContent := "SOME FILE CONTENT"

	client := YaDiskClient()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	randomName := "/" + testhelpers.RandString(50) + ".txt"

	uploadLink, err := client.RequestUploadLink(ctx, randomName, true)
	assert.Nil(t, err)

	status, err := client.Upload(ctx, uploadLink, strings.NewReader(fileContent))
	assert.Nil(t, err)

	// First delete op is successful
	deleteLink, status, err := client.Delete(ctx, randomName, true)
	assert.Nil(t, err)
	assert.Nil(t, deleteLink)
	assert.Equal(t, http.StatusNoContent, status)

	// Seconds delete op should return 404 NotFound (because it's been deleted already)
	deleteOriginalLink, status, err := client.Delete(ctx, randomName, true)
	assert.NotNil(t, err)
	assert.Nil(t, deleteOriginalLink)
	assert.Equal(t, http.StatusNotFound, status)
}
