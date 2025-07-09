package server

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var jwtKey = []byte("mySecretKey")

func validateJWT(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["username"].(string), nil
	}

	return "", errors.New("invalid token")
}

func authInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	if info.FullMethod != "/pb.User/Login" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no metadata")
	}

	token := md.Get("authorization")
	if len(token) == 0 {
		return nil, errors.New("no token")
	}

	tokenValue, err := validateJWT(token[0][7:])
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, "username", tokenValue)
	return handler(ctx, req)
}
