package clever

import (
	internal2 "cognito-openid-connectors/providers/clever/internal"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"cognito-openid-connectors/auth"
	"cognito-openid-connectors/common"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dgrijalva/jwt-go"
	jsoniter "github.com/json-iterator/go"
)

type controller struct {
	client        clever
	openIDConnect openIDConnect
}

func NewController(client clever, openIDConnect openIDConnect) auth.IController {
	return &controller{client: client, openIDConnect: openIDConnect}
}

func (c controller) GetWellKnowConfig(
	_ context.Context,
	req *events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	host, ok := req.Headers["Host"]
	if !ok {
		return common.ClientError(http.StatusBadRequest)
	}
	if req.RequestContext.Stage != "" {
		host = fmt.Sprintf("%s/%s", host, req.RequestContext.Stage)
	}
	resp, err := jsoniter.MarshalToString(c.openIDConnect.WellKnownConfiguration(host))
	if err != nil {
		return common.ServerError(err)
	}
	return common.Success(resp)
}

func (c controller) GetJSONWebKey(
	_ context.Context,
	_ *events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return common.ServerError(err)
	}
	pubKey, err := ioutil.ReadFile(filepath.Join(cwd, "resources", "jwtRS256.key.pub"))
	if err != nil {
		return common.ServerError(err)
	}
	resp, err := c.openIDConnect.WillKnownJWKSJSON(pubKey, os.Getenv(internal2.CleverAppKid))
	if err != nil {
		return common.ServerError(err)
	}
	return common.Success(*resp)
}

func (c controller) GetUserInfo(
	ctx context.Context,
	req *events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	token, err := auth.GetBearerToken(req)
	if err != nil {
		return common.ServerError(err)
	}
	userInfo, err := c.client.GetUserInfo(ctx, token)
	if err != nil {
		return common.ServerError(err)
	}

	user, err := c.client.GetUser(ctx, token, userInfo.Data.ID)
	if err != nil {
		return common.ServerError(err)
	}

	claim := internal2.Claim{
		Name:       fmt.Sprintf("%s %s", user.Name.First, user.Name.Last),
		FistName:   user.Name.First,
		LastName:   user.Name.Last,
		Email:      user.Email,
		DistrictID: user.District,
		UserName:   user.ID,
	}

	resp, err := jsoniter.MarshalToString(&claim)
	if err != nil {
		return common.ServerError(err)
	}
	return common.Success(resp)
}

func (c controller) Authorize(
	_ context.Context,
	req *events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	clientID, _ := common.GetStringKey(req.QueryStringParameters, "client_id")
	state, _ := common.GetStringKey(req.QueryStringParameters, "state")
	responseType, _ := common.GetStringKey(req.QueryStringParameters, "response_type")
	redirectURI, _ := common.GetStringKey(req.QueryStringParameters, "redirect_uri")
	scope, _ := common.GetStringKey(req.QueryStringParameters, "redirect_uri")
	authorizeURL := c.client.GetAuthorizeURL(clientID, redirectURI, scope, responseType, state)
	return common.Redirect(authorizeURL)

}

func (c *controller) GetToken(
	ctx context.Context,
	req *events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	query, err := url.ParseQuery(req.Body)
	if err != nil {
		return common.ServerError(err)
	}
	code := query.Get("code")
	redirectURI := query.Get("redirect_uri")
	scope := query.Get("scope")
	clientID := query.Get("client_id")
	clientSecret := query.Get("client_secret")
	cleverError := query.Get("error")
	if cleverError != "" {
		return common.ServerError(errors.New(cleverError))
	}
	token, err := c.client.GetToken(ctx, clientID, clientSecret, code, redirectURI)
	if err != nil {
		return common.ServerError(errors.New(cleverError))
	}

	cwd, err := os.Getwd()
	if err != nil {
		return common.ServerError(err)
	}
	privateKey, err := ioutil.ReadFile(filepath.Join(cwd, "resources", "jwtRS256.key"))
	if err != nil {
		return common.ServerError(err)
	}
	authAPIUrl := os.Getenv(internal2.CleverAuthApiUrl)
	crypt := auth.NewCrypto()
	claims := jwt.StandardClaims{
		Subject:   token.Subject,
		ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
		Issuer:    authAPIUrl,
		Audience:  os.Getenv(internal2.CleverClientID),
		IssuedAt:  time.Now().Unix(),
	}
	idToken, err := crypt.GetIDToken(privateKey, claims, os.Getenv(os.Getenv(internal2.CleverAppKid)))
	if err != nil {
		return common.ServerError(errors.New(cleverError))
	}

	cognitoToken := map[string]string{
		"id_token":      *idToken,
		"access_token":  token.AccessToken,
		"scope":         scope,
		"refresh_token": "",
		"token_type":    "Bearer",
	}

	resp, err := jsoniter.MarshalToString(&cognitoToken)
	if err != nil {
		return common.ServerError(errors.New(cleverError))
	}
	return common.Success(resp)
}
