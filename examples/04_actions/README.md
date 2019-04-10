# Yandex.Disk Client Actions

## Building
Build example as following:
```bash
go build -o ./bin/yandex-ops ./examples/05_actions/main.go
```

## Copy
Copy file form `src` to `dst`:

```bash
./bin/yandex-ops                      \
    -access-token="YOUR-ACCESS-TOKEN" \
    -op copy                          \
    -src=/whatever-file.txt           \
    -dst=/whatever-else.txt
```

## Move
Move file form `src` to `dst`:

```bash
./bin/yandex-ops                      \
    -access-token="YOUR-ACCESS-TOKEN" \
    -op move                          \
    -src=/whatever-file.txt           \
    -dst=/whatever-else.txt
```

## Delete
Delete file in `src`:

```bash
./bin/yandex-ops                      \
    -access-token="YOUR-ACCESS-TOKEN" \
    -op delete                        \
    -src=/whatever-file.txt
```
## Create Directory
Create directory specified in `src`:

```bash
./bin/yandex-ops                      \
    -access-token="YOUR-ACCESS-TOKEN" \
    -op mkdir                         \
    -src=/whatever-directory-name
```
