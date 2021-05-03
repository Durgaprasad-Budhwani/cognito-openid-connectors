package clever

import (
	internal2 "cognito-openid-connectors/providers/clever/internal"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"cognito-openid-connectors/common"

	jsoniter "github.com/json-iterator/go"
)

type clever struct {
}

func NewClever() clever {
	return clever{}
}

func (c clever) GetToken(
	ctx context.Context,
	clientID, clientSecret, code, redirectURI string,
) (*internal2.CleverToken, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", redirectURI) // this is cognito redirect URL
	tokenEndpoint := os.Getenv(internal2.CleverTokenEndpoint)
	httpClient := common.NewClientCredentialsHTTPClient(&clientID, &clientSecret)
	res, err := httpClient.PostForm(ctx, tokenEndpoint, data)
	if err != nil {
		return nil, err
	}
	var token internal2.CleverToken
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

func (c clever) GetAuthorizeURL(clientID, redirectURI, scope, responseType, state string) string {
	authorizationEndpoint := os.Getenv(internal2.CleverAuthorizationEndpoint)
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

func (c clever) GetUserInfo(ctx context.Context, token string) (*internal2.CleverUserInfo, error) {
	clientURL := fmt.Sprintf("%s%s%s", os.Getenv(internal2.CleverAPIEndpoint), internal2.CleverAPIVersion, "/me")
	httpClient := common.NewAccessTokenClient(&token)
	res, err := httpClient.Get(ctx, clientURL)
	if err != nil {
		return nil, err
	}
	var cleverUser internal2.CleverUserInfo
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(res, &cleverUser)
	return &cleverUser, err
}

func (c clever) GetUser(ctx context.Context, token, id string) (*internal2.CleverUser, error) {
	clientURL := fmt.Sprintf("%s%s%s%s", os.Getenv(internal2.CleverAPIEndpoint), internal2.CleverAPIVersion, "/users/", id)
	httpClient := common.NewAccessTokenClient(&token)
	res, err := httpClient.Get(ctx, clientURL)
	if err != nil {
		return nil, err
	}
	var resp internal2.CleverUserResp
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(res, &resp)
	return &resp.Data, err
}
