package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/yurykabanov/go-yandex-disk"
)

var (
	accessToken = flag.String("access-token", "", "Access Token")

	op = flag.String("op", "", "Operation: copy, move, delete or mkdir")

	src = flag.String("src", "", "Source file or directory")
	dst = flag.String("dst", "", "Destination path")
)

func main() {
	flag.Parse()

	client := yadisk.NewFromAccessToken(*accessToken)

	switch *op {
	case "copy":
		fmt.Printf("Copy '%s' -> '%s'\n", *src, *dst)
		link, status, err := client.Copy(context.TODO(), *src, *dst, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("HTTP Status code: %d\nLink: %+v\n", status, link)

	case "move":
		fmt.Printf("Move '%s' -> '%s'\n", *src, *dst)
		link, status, err := client.Move(context.TODO(), *src, *dst, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("HTTP Status code: %d\nLink: %+v\n", status, link)

	case "delete":
		fmt.Printf("Delete '%s'\n", *src)
		link, status, err := client.Delete(context.TODO(), *src, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("HTTP Status code: %d\nLink: %+v\n", status, link)

	case "mkdir":
		fmt.Printf("Mkdir '%s'\n", *src)
		link, err := client.CreateDirectory(context.TODO(), *src)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Link: %+v\n", link)

	default:
		fmt.Println("You should specify operation using -op=<...> with one of the following values: 'copy', 'move', 'delete' or 'mkdir'")
	}
}
