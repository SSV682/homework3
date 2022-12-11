package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
)

type MyCustomClaims struct {
	jwt.RegisteredClaims
	Exp    string `json:"exp"`
	IDUser int64  `json:"id_user"`
}

func getUserByToken(ctx context.Context, tokenString string) {
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJpc3MiOiJ0ZXN0IiwiYXVkIjoic2luZ2xlIn0.QAWg1vGvnqRuCFTMcPkjZljXHh8U3L_qUjszOtQbeaA"
	//
	//token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	//	return []byte("AllYourBase"), nil
	//})
	//
	//if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
	//	fmt.Printf("%v %v", claims.Foo, claims.RegisteredClaims.Issuer)
	//} else {
	//	fmt.Println(err)
	//}

}
