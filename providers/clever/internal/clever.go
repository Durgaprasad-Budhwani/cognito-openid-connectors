package internal

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"cognito-openid-connectors/common"

	jsoniter "github.com/json-iterator/go"
)

type Clever struct {
}

func NewClever() Clever {
	return Clever{}
}

func (c Clever) GetToken(
	ctx context.Context,
	clientID, clientSecret, code, redirectURI string,
) (*CleverToken, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", redirectURI) // this is cognito redirect URL
	tokenEndpoint := os.Getenv(CleverTokenEndpoint)
	httpClient := common.NewClientCredentialsHTTPClient(&clientID, &clientSecret)
	res, err := httpClient.PostForm(ctx, tokenEndpoint, data)
	if err != nil {
		return nil, err
	}
	var token CleverToken
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(res, &token)
	if err != nil {
		return nil, err
	}

	userInfo, err := c.GetUserInfo(ctx, token.AccessToken)
	if err != nil {
		return nil, err
	}
	token.Subject = userInfo.Data.ID
	return &token, err
}

func (c Clever) GetAuthorizeURL(clientID, redirectURI, scope, responseType, state string) string {
	authorizationEndpoint := os.Getenv(CleverAuthorizationEndpoint)
	URL, err := url.Parse(authorizationEndpoint)
	if err != nil {
		panic("error parsing url")
	}

	parameters := url.Values{}
	parameters.Add("client_id", clientID)
	parameters.Add("response_type", responseType)
	parameters.Add("scope", scope)
	parameters.Add("redirect_uri", redirectURI)
	parameters.Add("state", state)

	URL.RawQuery = parameters.Encode()

	log.Printf("Encoded URL is %q\n", URL.String())
	return URL.String()
}

func (c Clever) GetUserInfo(ctx context.Context, token string) (*CleverUserInfo, error) {
	clientURL := fmt.Sprintf("%s%s%s", os.Getenv(CleverAPIEndpoint), CleverAPIVersion, "/me")
	httpClient := common.NewAccessTokenClient(&token)
	res, err := httpClient.Get(ctx, clientURL)
	if err != nil {
		return nil, err
	}
	var cleverUser CleverUserInfo
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(res, &cleverUser)
	return &cleverUser, err
}

func (c Clever) GetUser(ctx context.Context, token, id string) (*CleverUser, error) {
	clientURL := fmt.Sprintf("%s%s%s%s", os.Getenv(CleverAPIEndpoint), CleverAPIVersion, "/users/", id)
	httpClient := common.NewAccessTokenClient(&token)
	res, err := httpClient.Get(ctx, clientURL)
	if err != nil {
		return nil, err
	}
	var resp CleverUserResp
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(res, &resp)
	return &resp.Data, err
}
