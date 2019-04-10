# Yandex.Disk Client Download

## Building
Build example as following:
```bash
go build -o ./bin/yandex-download ./examples/04_download/main.go
```

## Running
Download file `/some/existing/path/some-existing-filename.txt` and print
its contents.
```bash
./bin/yandex-download                            \
    -access-token="YOUR-ACCESS-TOKEN"            \ 
    -filename="/some/existing/path/some-existing-filename.txt"
```
