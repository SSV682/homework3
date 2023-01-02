package token

type tokenProvider struct {
	//verifyKey *rsa.PublicKey
	//secretKey string
}

const (
	ErrorParseWithClaims = "parse with claims: %v"
	ErrorInvalidToken    = "invalid token: %v"
)

func NewJWTProvider() *tokenProvider {
	return &tokenProvider{}
}

//func getKey(token *jwt.Token) (interface{}, error) {
//	keyID, ok := token.Header["kid"].(string)
//	if !ok {
//		return nil, errors.New("expecting JWT header to have string kid")
//	}
//
//	if key, found := JWKS[keyID]; found {
//		return key, nil
//	}
//
//	return nil, fmt.Errorf("unable to find key %q", keyID)
//}

//func (t *tokenProvider) ParseToken(tokenString string) (jwt.MapClaims, error) {
//
//	key, err := jwt.Parse(tokenString, getKey)
//	if err != nil {
//		return nil, fmt.Errorf("validate parse key: %w", err)
//	}
//
//	token, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
//		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
//			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
//		}
//
//		return key, nil
//	})
//	if err != nil {
//		return nil, fmt.Errorf("validate: %w", err)
//	}
//
//	log.Info(token.Claims)
//	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//		return claims, nil
//	} else {
//		return nil, fmt.Errorf(ErrorInvalidToken, err)
//	}
//}
