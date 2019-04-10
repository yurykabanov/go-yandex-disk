# Yandex.Disk Client Upload

## Building
Build example as following:
```bash
go build -o ./bin/yandex-upload ./examples/03_upload/main.go

```

## Running
Read contents from stdin and upload them as file 
`/some/existing/path/filename.txt`.

```bash
echo 123 | ./bin/yandex-upload                   \
    -access-token="YOUR-ACCESS-TOKEN"            \ 
    -filename="/some/existing/path/filename.txt"
```
