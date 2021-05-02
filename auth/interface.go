package auth

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type IController interface {
	GetWellKnowConfig(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	GetJSONWebKey(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	Authorize(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	GetUserInfo(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	GetToken(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}
