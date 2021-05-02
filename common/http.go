package common

import (
	"context"
	"net/http"
	netUrl "net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

var isDebug = os.Getenv(os.Getenv(Debug)) == TrueString
var isOffline = os.Getenv(os.Getenv(ServerlessOffline)) == TrueString

const (
	timeout    = 10
	backOffLow = 10
	retry      = 2
)

type HTTPClient struct {
	clientID     *string
	clientSecret *string
	accessToken  *string
	client       http.Client
}

type IHttpClient interface {
	Get(ctx context.Context, url string) ([]byte, error)
	Post(ctx context.Context, url, body string) ([]byte, error)
	PostForm(ctx context.Context, url string, data netUrl.Values) ([]byte, error)
}

func NewClientCredentialsHTTPClient(clientID, clientSecret *string) IHttpClient {
	return &HTTPClient{clientID: clientID, clientSecret: clientSecret, client: http.Client{}}
}

func NewAccessTokenClient(accessToken *string) IHttpClient {
	return &HTTPClient{accessToken: accessToken, client: http.Client{}}
}

func (c *HTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	httpClient := c.getHTTPClient()
	httpClient.SetHeader("Content-Type", "application/json")
	response, err := httpClient.R().SetContext(ctx).Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make a request to server")
	}

	if response.RawResponse.StatusCode != http.StatusOK {
		return nil, errors.New(string(response.Body()))
	}
	return response.Body(), nil
}

func (c *HTTPClient) Post(ctx context.Context, url, body string) ([]byte, error) {
	httpClient := c.getHTTPClient()
	httpClient.SetHeader("Content-Type", "application/json")
	response, err := httpClient.R().SetContext(ctx).SetBody(body).Post(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make a request to server")
	}

	if response.RawResponse.StatusCode != http.StatusOK {
		return nil, errors.New(string(response.Body()))
	}
	return response.Body(), nil
}

func (c *HTTPClient) PostForm(ctx context.Context, url string, data netUrl.Values) ([]byte, error) {
	httpClient := c.getHTTPClient()
	response, err := httpClient.R().SetContext(ctx).SetFormDataFromValues(data).Post(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make a request to server")
	}

	if response.RawResponse.StatusCode != http.StatusOK {
		return nil, errors.New(string(response.Body()))
	}
	return response.Body(), nil
}

func (c *HTTPClient) getHTTPClient() *resty.Client {
	client := resty.NewWithClient(xray.Client(nil))
	if isOffline {
		client = resty.NewWithClient(&http.Client{})
	}
	client.
		SetTimeout(timeout * time.Second).
		SetDebug(isDebug).
		SetRetryCount(retry).
		SetRetryWaitTime(backOffLow * time.Second)

	client.SetHeader("accept", "application/json")
	client.SetHeader("cache-control", "no-store, max-age=0")
	if c.clientID != nil && c.clientSecret != nil {
		client.SetBasicAuth(aws.StringValue(c.clientID), aws.StringValue(c.clientSecret))
	}
	if c.accessToken != nil {
		client.SetAuthToken(aws.StringValue(c.accessToken))
	}
	return client
}
