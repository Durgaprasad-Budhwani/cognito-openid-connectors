package main

import (
	"context"
	"net/http"

	"cognito-openid-connectors/clever"
	"cognito-openid-connectors/common"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// nolint(gocritic)
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var authController = clever.NewController(clever.NewClever(), clever.NewOpenIDConnect())
	switch req.Resource {
	// open id connect endpoints
	case common.GetStage() + "/auth/clever/token":
		return authController.GetToken(ctx, &req)
	case common.GetStage() + "/auth/clever/authorize":
		return authController.Authorize(ctx, &req)
	case common.GetStage() + "/auth/clever/userinfo":
		return authController.GetUserInfo(ctx, &req)
	case common.GetStage() + "/auth/clever/.well-known/openid-configuration":
		return authController.GetWellKnowConfig(ctx, &req)
	case common.GetStage() + "/auth/clever/.well-known/jwks.json":
		return authController.GetJSONWebKey(ctx, &req)
	}
	return common.ClientError(http.StatusNotFound)
}

func main() {
	lambda.Start(Handler)
}
