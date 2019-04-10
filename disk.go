package yadisk

import (
	"context"
	"net/http"
)

const (
	methodGetDisk = http.MethodGet
	urlGetDisk    = ""
)

// Data about a user's Disk.
//
// Method returns Disk stats or error.
//
// See: https://tech.yandex.com/disk/api/reference/capacity-docpage/
func (c *Client) GetDisk(ctx context.Context) (*Disk, error) {
	var disk Disk

	_, err := c.doRequestAndDecode(ctx, methodGetDisk, urlGetDisk, nil, nil, &disk)
	if err != nil {
		return nil, err
	}

	return &disk, nil
}
