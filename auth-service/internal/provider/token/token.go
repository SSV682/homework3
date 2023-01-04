package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"time"
)

type tokenProvider struct {
	privateSet jwk.Set
	//publicSet  jwk.Set
	//privateKey *rsa.PrivateKey
	//publicKey  *rsa.PublicKey
}

const (
	ErrorParseWithClaims = "parse with claims: %v"
	ErrorInvalidToken    = "invalid token: %v"
)

func NewJWTProvider() *tokenProvider {
	raw, err := rsa.GenerateKey(rand.Reader, 2048)

	hasher := sha1.New()
	hasher.Write([]byte(time.Now().String()))
	kid := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	privateJWK, err := jwk.FromRaw(raw)
	if err != nil {
		fmt.Printf("failed to create asymmetric key: %s\n", err)
	}
	//if _, ok := privateJWK.(jwk.RSAPrivateKey); !ok {
	//	fmt.Printf("expected jwk.SymmetricKey, got %T\n", jwkKey)
	//}

	privateJWK.Set(jwk.KeyIDKey, kid)
	privateJWK.Set(jwk.AlgorithmKey, jwa.RS256)

	privateSet := jwk.NewSet()
	privateSet.AddKey(privateJWK)

	return &tokenProvider{
		privateSet: privateSet,
		//privateKey: raw,
		//publicKey:  &raw.PublicKey,
	}
}

func (t *tokenProvider) CreateToken(userID string) (string, error) {
	tok, err := jwt.NewBuilder().
		Issuer("http://userservice-authservice.userservice.svc.cluster.local").
		IssuedAt(time.Now()).
		NotBefore(time.Now()).
		Expiration(time.Now().Add(24*time.Hour)).
		JwtID(uuid.New().String()).
		Claim("id_user", userID).
		Build()
	if err != nil {
		return "", fmt.Errorf("doesnt build token: %v", err)
	}

	buf, err := json.MarshalIndent(tok, "", " ")
	if err != nil {
		return "", fmt.Errorf("failed to generate JSON: %s", err)
	}

	key, _ := t.privateSet.Key(0)

	headers := jws.NewHeaders()
	headers.Set(jws.TypeKey, "JWT")
	signed, err := jws.Sign(buf, jws.WithKey(jwa.RS256, key, jws.WithProtectedHeaders(headers)))
	if err != nil {
		return "", fmt.Errorf("doesnt sign token: %v", err)
	}
	return string(signed), nil
}

func (t *tokenProvider) ParseToken(tokenString string) (jwt.Token, error) {
	set, err := t.GetKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to verify JWS: %s\n", err)
	}

	verifiedToken, err := jwt.Parse([]byte(tokenString), jwt.WithValidate(true), jwt.WithKeySet(set))
	if err != nil {
		return nil, fmt.Errorf("failed to verify JWS: %s\n", err)
	}

	return verifiedToken, nil
}

func (t *tokenProvider) GetKeys() (jwk.Set, error) {
	set, err := jwk.PublicSetOf(t.privateSet)
	if err != nil {
		return nil, fmt.Errorf("cant get set: %s", err)
	}
	return set, nil
}
