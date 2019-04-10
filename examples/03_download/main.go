package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/yurykabanov/go-yandex-disk"
)

var (
	accessToken = flag.String("access-token", "", "Access Token")

	fileName = flag.String("filename", "/whatever-file.txt", "File to download")
)

func main() {
	flag.Parse()

	client := yadisk.NewFromAccessToken(*accessToken)

	// 1. Request download link
	link, err := client.RequestDownloadLink(context.TODO(), *fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 2. Download file
	resp, err := client.Download(context.TODO(), link)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// 3. Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("File contents:\n%s", body)
}
