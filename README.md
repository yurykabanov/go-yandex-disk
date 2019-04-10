# Yandex.Disk Client

This is unofficial Yandex.Disk API client in Go.

## Getting started

### Installation

Import it to your project:
```
import "github.com/yurykabanov/go-yandex-disk"
```

### Usage

Instantiate client:
```
client := yadisk.New(&http.Client{ /* ... */ })
// or
client := yadisk.NewFromAccessToken("YOUR-OAUTH-ACCESS-TOKEN")
// or
client := yadisk.NewFromConfigAndToken(yandexOauthConfig, oauthToken, context.TODO())
```

and use it:
```
# Upload
link, err := client.RequestUploadLink(context.TODO(), "/some-path/uploaded-file.txt", false)
status, err := client.Upload(context.TODO(), link, anyIoReader)

# Download
link, err := client.RequestDownloadLink(context.TODO(), "/some-path/existing-file.txt")
resp, err := client.Download(context.TODO(), link)
defer resp.Body.Close()
// resp.Body is io.Reader for requested file

# Actions
client.Copy(context.TODO(), "/some-path/source-file.txt", "/some-path/destination-file.txt", false)
client.Move(context.TODO(), "/some-path/source-file.txt", "/some-path/destination-file.txt", false)
client.Delete(context.TODO(), "/some-path/existing-file.txt", false)
client.Mkdir(context.TODO(), "/some-path/new-directory")
```

More detailed examples could be found in `examples/` directory.

## Supported methods

This client currently support the following methods:
- [x] Upload and download
- [x] Actions: copy, move, delete and create directory
- [x] Disk stats
- [ ] File meta information actions (read/write)
- [ ] Publishing resources and performing actions on them
- [ ] Working with Trash

## Running the tests

```bash
# Unit tests
go test . -v

# Integration tests
go test ./test/integration/ -v -tags=integration -access-token=YOUR-OAUTH-ACCESS-TOKEN
```

## Contributing
Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process
for submitting pull requests.

## Versioning
We use [SemVer](http://semver.org/) for versioning. For the versions available, 
see the tags [on this repository](https://github.com/yurykabanov/go-yandex-disk/tags).

## License
This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
