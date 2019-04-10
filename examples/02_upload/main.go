package main

import (
	"context"
	"flag"
	"fmt"
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

	// 1. Request link to upload file
	link, err := client.RequestUploadLink(context.TODO(), *fileName, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 2. Upload file
	status, err := client.Upload(context.TODO(), link, os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 3. Successful upload should have either:
	// 201 Created - file successfully uploaded and processed
	// 202 Accepted - file successfully uploaded but not processed yet (eventually it will be processed)
	fmt.Printf("HTTP status code: %d\n", status)
}
