package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func NewKeycloakJWTValidator(issuerUrl, clientId string) (func(*fiber.Ctx, string) (bool, error), error) {
	fmt.Fprintln(os.Stdout, "issuerUrl: ", issuerUrl, " clientId:", clientId)
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, issuerUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		return nil, err
	}
	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientId,
	})
	fmt.Fprintln(os.Stdout, "VERIFIER OUT:", verifier)

	return func(c *fiber.Ctx, key string) (bool, error) {
		var ctx = c.UserContext()
		_, err := verifier.Verify(ctx, key)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			return false, err
		}
		token, _ := jwt.ParseWithClaims(key, &Claims{},
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v",
						token.Header["alg"])
				}
				return key, nil
			})
		c.Locals("claims", token.Claims)
		return true, nil
	}, nil
}
