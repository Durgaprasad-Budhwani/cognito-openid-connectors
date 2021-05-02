package auth

import (
	"crypto/x509"
	"encoding/pem"
	"log"

	"github.com/dgrijalva/jwt-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/lestrrat-go/jwx/jwk"
)

type Crypto struct {
}

func NewCrypto() *Crypto {
	return &Crypto{}
}

func (c Crypto) GetJSONWebKey(pubKey []byte, kid string) (*string, error) {
	pubPem, _ := pem.Decode(pubKey)
	pub, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
	if err != nil {
		log.Printf("failed to convert to ParsePKCS1PrivateKey: %s", err)
		log.Printf("failed to convert to JWK: %s", err)
		return nil, err
	}

	set, err := jwk.New(pub)
	if err != nil {
		log.Printf("failed to convert to JWK: %s", err)
		return nil, err
	}

	err = jwk.AssignKeyID(set)
	if err != nil {
		log.Printf("failed to assign kid: %s", err)
		log.Printf("failed to convert to JWK: %s", err)
		return nil, err
	}

	type resp struct {
		Keys []jwk.Key `json:"keys"`
	}

	err = set.Set("kid", kid)
	if err != nil {
		return nil, err
	}
	result, err := jsoniter.MarshalToString(&resp{Keys: []jwk.Key{set}})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c Crypto) GetIDToken(privateKey []byte, claims jwt.Claims, kid string) (*string, error) {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return nil, err
	}
	return &signedToken, err
}
