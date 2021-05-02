package common

import (
	"github.com/davecgh/go-spew/spew"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	jsoniter "github.com/json-iterator/go"
)

// ServerError Add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func ServerError(err error) (events.APIGatewayProxyResponse, error) {
	spew.Dump(err)
	return ResponseError(http.StatusInternalServerError, err.Error())
}

func ResponseError(status int, data string) (events.APIGatewayProxyResponse, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	body, _ := json.MarshalToString(map[string]string{"message": data})
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       body,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Cache-Control":                    "no-cache",
		},
	}, nil
}

// ClientError add a helper for send responses relating to client errors.
func ClientError(status int) (events.APIGatewayProxyResponse, error) {
	return ResponseError(status, http.StatusText(status))
}

// Success add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func Success(resp string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       resp,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Credentials": "true",
			"Cache-Control":                    "no-cache",
		},
	}, nil
}

// Redirect add a helper for redirecting return to URL. it format the data that the AWS API
// Gateway understands.
func Redirect(url string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusTemporaryRedirect,
		Body:       "",
		Headers: map[string]string{
			"Location":      url,
			"Authorization": "",
		},
	}, nil
}
