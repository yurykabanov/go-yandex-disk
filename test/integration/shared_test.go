package integration

import (
	"flag"

	"github.com/yurykabanov/go-yandex-disk"
)

var (
	accessToken = flag.String("access-token", "", "Access Token")
)

func YaDiskClient() *yadisk.Client {
	return yadisk.NewFromAccessToken(*accessToken)
}
