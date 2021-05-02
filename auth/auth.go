package auth

import (
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	jsoniter "github.com/json-iterator/go"
)

func stripBearerPrefixFromTokenString(tok string) string {
	// Should be a bearer token
	if len(tok) > 6 && strings.EqualFold(tok[0:7], "BEARER ") {
		return tok[7:]
	}
	return tok
}

func GetBearerToken(req *events.APIGatewayProxyRequest) (string, error) {
	if authHeader, ok := req.Headers["Authorization"]; ok {
		return stripBearerPrefixFromTokenString(authHeader), nil
	}
	if accessToken, ok := req.QueryStringParameters["access_token"]; ok {
		return accessToken, nil
	}
	if _, ok := req.Headers["application/x-www-form-urlencoded"]; ok && req.Body != "" {
		var token Token
		err := jsoniter.UnmarshalFromString(req.Body, &token)
		if err != nil {
			return "", err
		}
		return token.AccessToken, nil
	}
	return "", errors.New("no token specified in request")
}
