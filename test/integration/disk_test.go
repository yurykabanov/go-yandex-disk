// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestYaDiskClient_GetDisk(t *testing.T) {
	client := YaDiskClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	disk, err := client.GetDisk(ctx)

	assert.NotNil(t, disk)
	assert.Nil(t, err)

	assert.True(t, disk.TotalSpace > 0)
}
