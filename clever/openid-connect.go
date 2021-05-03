package clever

import (
	"cognito-openid-connectors/auth"
	"fmt"
)

type openIDConnect struct {
}

func NewOpenIDConnect() openIDConnect {
	return openIDConnect{}
}

// WellKnownConfiguration openIDConnect Connect Discovery
// The well known endpoint an be used to retrieve information for openIDConnect Connect clients.
func (h openIDConnect) WellKnownConfiguration(host string) *auth.WellKnown {
	return &auth.WellKnown{
		Issuer:                            fmt.Sprintf("https://%s/auth/clever", host),
		AuthURL:                           fmt.Sprintf("https://%s/auth/clever/authorize", host),
		TokenURL:                          fmt.Sprintf("https://%s/auth/clever/token", host),
		TokenEndpointAuthMethodsSupported: []string{"client_secret_basic", "private_key_jwt"},
		UserinfoEndpoint:                  fmt.Sprintf("https://%s/auth/clever/userinfo", host),
		JWKsURI:                           fmt.Sprintf("https://%s/auth/clever/.well-known/jwks.json", host),
		ScopesSupported: []string{
			"read:district_admins_basic",
			"read:school_admins_basic",
			"read:students_basic",
			"read:teachers_basic",
			"read:user_id",
		},
		ResponseTypes:                          []string{"code", "code id_token", "id_token", "token id_token"},
		SubjectTypes:                           []string{"public"},
		UserinfoSigningAlgValuesSupported:      []string{"none"},
		IDTokenSigningAlgValuesSupported:       []string{"RS256"},
		RequestObjectSigningAlgValuesSupported: []string{"none"},
		DisplayValuesSupported:                 []string{"page", "popup"},
		ClaimsSupported: []string{
			"sub",
			"name",
			"preferred_username",
			"profile",
			"picture",
			"website",
			"email",
			"email_verified",
			"updated_at",
			"iss",
			"aud",
		},
	}
}

// WillKnownJWKSJSON JSON Web Keys Discovery
// This endpoint returns JSON Web Keys to be used as public keys for verifying OpenID Authorize ID Tokens and,
// if enabled, OAuth 2.0 JWT Access Tokens. This endpoint can be used with client libraries like
func (h openIDConnect) WillKnownJWKSJSON(pubKey []byte, kid string) (*string, error) {
	crypto := auth.NewCrypto()
	return crypto.GetJSONWebKey(pubKey, kid)
}
