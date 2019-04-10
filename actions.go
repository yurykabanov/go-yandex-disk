package yadisk

import (
	"context"
	"net/http"
)

const (
	methodActionCopy = http.MethodPost
	urlActionCopy    = "resources/copy"

	methodActionMove = http.MethodPost
	urlActionMove    = "resources/move"

	methodActionDelete = http.MethodDelete
	urlActionDelete    = "resources"

	methodActionCreateDirectory = http.MethodPut
	urlActionCreateDirectory    = "resources"
)

// Copy file or directory.
//
// src - The path to the resource to copy.
// dst - The path to the copy of the resource that is being created. The name
// of the file can be up to 255 characters. The path can be up to 32760
// characters long.
// permanently -- Whether to permanently the file. It is used if the resource is
// copied to a folder that already contains a resource with the same name.
//
// Method returns Link to created resource, status code and error (if any).
//
// NOTE: for files and empty directories status code is "201 Created" and for
// non-empty directories it is "202 Accepted" which means the operation has
// been started, but hasn't been finished yet. The application MUST track
// the status of operation by itself.
//
// See: https://tech.yandex.com/disk/api/reference/copy-docpage/
func (c *Client) Copy(ctx context.Context, src, dst string, overwrite bool) (*Link, int, error) {
	var link Link

	params := map[string]string{
		"from":      src,
		"path":      dst,
		"permanently": "false",
	}

	if overwrite {
		params["permanently"] = "true"
	}

	statusCode, err := c.doRequestAndDecode(ctx, methodActionCopy, urlActionCopy, params, nil, &link)
	if err != nil {
		return nil, statusCode, err
	}

	return &link, statusCode, nil
}

// Move file or directory.
//
// src - The path to the resource to move.
// dst - The path to the new location of the resource. The name of the file can
// be up to 255 characters. The path can be up to 32760 characters long.
// permanently - Whether to permanently files. It is used if the resource is moved
// to a folder that already contains a resource with the same name.
//
// Method returns Link to created resource, status code and error (if any).
//
// NOTE: for files and empty directories status code is "201 Created" and for
// non-empty directories it is "202 Accepted" which means the operation has
// been started, but hasn't been finished yet. The application MUST track
// the status of operation by itself.
//
// See: https://tech.yandex.com/disk/api/reference/move-docpage/
func (c *Client) Move(ctx context.Context, src, dst string, overwrite bool) (*Link, int, error) {
	var link Link

	params := map[string]string{
		"from":      src,
		"path":      dst,
		"permanently": "false",
	}

	if overwrite {
		params["permanently"] = "true"
	}

	statusCode, err := c.doRequestAndDecode(ctx, methodActionMove, urlActionMove, params, nil, &link)
	if err != nil {
		return nil, statusCode, err
	}

	return &link, statusCode, nil
}

// Delete file or directory.
//
// path - The path to the resource to delete.
// permanently - The flag for permanent deletion. False means file will be
// moved into the Trash.
//
// Method returns Link only if operation hasn't been completed yet. It also
// returns status code and error (if any).
//
// NOTE: for files and empty directories status code is "204 No content" and for
// non-empty directories it is "202 Accepted" which means the operation has
// been started, but hasn't been finished yet. The application MUST track
// the status of operation by itself.
//
// See: https://tech.yandex.com/disk/api/reference/delete-docpage/
func (c *Client) Delete(ctx context.Context, path string, permanently bool) (*Link, int, error) {
	var link Link

	params := map[string]string{
		"path":        path,
		"permanently": "false",
	}

	if permanently {
		params["permanently"] = "true"
	}

	statusCode, err := c.doRequestAndDecode(ctx, methodActionDelete, urlActionDelete, params, nil, &link)
	if err != nil {
		return nil, statusCode, err
	}

	if statusCode == http.StatusNoContent {
		return nil, statusCode, nil
	}

	return &link, statusCode, nil
}

// Create directory.
//
// path - The path to the folder being created. The maximum length of the
// folder name is 255 characters; the maximum length of the path is 32760
// characters.
//
// Method returns Link to created resource or error.
//
// See: https://tech.yandex.com/disk/api/reference/create-folder-docpage/
func (c *Client) CreateDirectory(ctx context.Context, path string) (*Link, error) {
	var link Link

	params := map[string]string{
		"path": path,
	}

	_, err := c.doRequestAndDecode(ctx, methodActionCreateDirectory, urlActionCreateDirectory, params, nil, &link)
	if err != nil {
		return nil, err
	}

	return &link, nil
}
