package internal

import (
	"fmt"

	"cognito-openid-connectors/auth"
)

type OpenIDConnect struct {
}

func NewOpenIDConnect() OpenIDConnect {
	return OpenIDConnect{}
}

// WellKnownConfiguration OpenIDConnect Connect Discovery
// The well known endpoint an be used to retrieve information for OpenIDConnect Connect clients.
func (h OpenIDConnect) WellKnownConfiguration(host string) *auth.WellKnown {
	return &auth.WellKnown{
		Issuer:                            fmt.Sprintf("https://%s/auth/Clever", host),
		AuthURL:                           fmt.Sprintf("https://%s/auth/Clever/authorize", host),
		TokenURL:                          fmt.Sprintf("https://%s/auth/Clever/token", host),
		TokenEndpointAuthMethodsSupported: []string{"client_secret_basic", "private_key_jwt"},
		UserinfoEndpoint:                  fmt.Sprintf("https://%s/auth/Clever/userinfo", host),
		JWKsURI:                           fmt.Sprintf("https://%s/auth/Clever/.well-known/jwks.json", host),
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
func (h OpenIDConnect) WillKnownJWKSJSON(pubKey []byte, kid string) (*string, error) {
	crypto := auth.NewCrypto()
	return crypto.GetJSONWebKey(pubKey, kid)
}
