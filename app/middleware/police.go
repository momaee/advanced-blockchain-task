package middleware

import (
	"backend_task/conf"
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

func AuthenticateToken(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")

	if err != nil {
		return nil, err
	}

	ok, claims := ValidateToken(token, conf.GetAppConfig().JWT.Secret)

	fmt.Println("token", ok, claims)
	if !ok {
		return nil, fmt.Errorf("token is not valid.")
	}

	grpc_tags.Extract(ctx).Set("auth.sub", "I AM SOMEBODY")
	// WARNING: in production define your own type to avoid context collisions
	newCtx := context.WithValue(ctx, "tokenInfo", "something I added")

	return newCtx, nil
}

// ValidateToken valid token
func ValidateToken(tokenString, secret string) (bool, jwt.MapClaims) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

	if err != nil {
		return false, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		// fmt.Println(claims["foo"], claims["nbf"])
		return true, claims
	}

	return false, nil
}
