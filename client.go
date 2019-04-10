package yadisk

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const (
	defaultBaseUrl = "https://cloud-api.yandex.net/v1/disk/"
)

type Client struct {
	client  *http.Client
	baseUrl *url.URL
}

func New(client *http.Client) *Client {
	baseUrl, _ := url.Parse(defaultBaseUrl)

	return &Client{
		client:  client,
		baseUrl: baseUrl,
	}
}

func NewFromConfigAndToken(config *oauth2.Config, token *oauth2.Token, ctx context.Context) *Client {
	source := config.TokenSource(ctx, token)

	httpClient := &http.Client{
		Transport: &oauth2.Transport{
			Source: source,
		},
	}

	return New(httpClient)
}

func NewFromAccessToken(accessToken string) *Client {
	source := oauth2.StaticTokenSource(&oauth2.Token{
		TokenType:   "OAuth",
		AccessToken: accessToken,
	})

	httpClient := &http.Client{
		Transport: &oauth2.Transport{
			Source: source,
		},
	}

	return New(httpClient)
}

func (c *Client) doRawRequest(
	ctx context.Context,
	method, absoluteUrl string,
	bodyReader io.Reader,
) (*http.Response, error) {
	req, err := http.NewRequest(method, absoluteUrl, bodyReader)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	return c.client.Do(req)
}

func (c *Client) doRequest(
	ctx context.Context,
	method, relativeUrl string,
	params map[string]string,
	body interface{},
) (*http.Response, error) {
	rel, err := url.Parse(relativeUrl)
	if err != nil {
		return nil, err
	}

	requestUrl := c.baseUrl.ResolveReference(rel)

	var bodyReader io.Reader

	if body != nil {
		bodyReader, err = c.encodeBody(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, requestUrl.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	return c.client.Do(req)
}

func (c *Client) doRequestAndDecode(
	ctx context.Context,
	method, relativeUrl string,
	params map[string]string,
	body, result interface{},
) (int, error) {
	code := 0

	resp, err := c.doRequest(ctx, method, relativeUrl, params, body)
	if resp != nil {
		code = resp.StatusCode
	}
	if err != nil {
		return code, err
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	return code, c.decodeResponseOrError(resp.StatusCode, dec, result)
}

func (c *Client) encodeBody(body interface{}) (io.Reader, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *Client) decodeResponseOrError(statusCode int, dec *json.Decoder, target interface{}) error {
	var err error

	// HTTP "204 No content" is not an error, though no response object could be extracted
	if statusCode == http.StatusNoContent {
		return nil
	}

	if statusCode >= http.StatusBadRequest {
		err = c.decodeError(dec, statusCode)
		return err
	}

	err = dec.Decode(target)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) decodeError(dec *json.Decoder, statusCode int) error {
	var apiError ApiError

	err := dec.Decode(&apiError)
	if err != nil {
		return err
	}

	apiError.StatusCode = statusCode

	return apiError
}
