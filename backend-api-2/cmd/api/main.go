package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	Username   string `json:"preferred_username"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

type ProfileResponse struct {
	Username   string `json:"preferred_username"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

// "http://keycloak.default.svc.cluster.local"
// "http://backend-api-1.default.svc.cluster.local"
// "http://backend-api-2.default.svc.cluster.local"

var KeycloakHost = os.Getenv("KEYCLOAK_HOST")
var Backend1Host = os.Getenv("BACKEND_API_1_HOST")
var Backend2Host = os.Getenv("BACKEND_API_2_HOST")
var ClientID = os.Getenv("CLIENT_ID")
var ClientSecret = os.Getenv("CLIENT_SECRET")

func main() {

	app := fiber.New()

	prometheus := fiberprometheus.New("my-service-name")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	validator, err := NewKeycloakJWTValidator(KeycloakHost+"/realms/master", ClientID)
	if err != nil {
		println(err)
		panic(err)
	}

	protected := app.Group("protected", keyauth.New(keyauth.Config{
		Validator: validator,
	}))
	protected.Get("profile/name", func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(*Claims)
		data, _ := json.MarshalIndent(
			ProfileResponse{
				Username:   claims.Username,
				Name:       claims.Name,
				GivenName:  claims.GivenName,
				FamilyName: claims.FamilyName,
			}, "", "  ")
		return c.SendString(string(data))
	})
	protected.Get("/pets/list", petHandler)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from backend2!")
	})
	app.Get("/400", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusBadRequest).SendString("Hello, BadRequest!")
	})
	app.Get("/500", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusInternalServerError).SendString("Hello, InternalServerError!")
	})
	app.Get("/501", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotImplemented).SendString("Hello, StatusNotImplemented!")
	})
	app.Get("/502", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusBadGateway).SendString("Hello, BadGateway!")
	})

	app.Listen(":3000")
}

func petHandler(c *fiber.Ctx) error {
	petlist := PetList{}
	petlist.PetsSlice = pets
	data, _ := json.Marshal(petlist)
	return c.Status(fiber.StatusOK).SendString(string(data))
}

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
		// fmt.Fprintln(os.Stdout, "User Context:", ctx)
		_, err := verifier.Verify(ctx, key)
		//  fmt.Fprintln(os.Stdout, "VERIFIER VERIFY", vvv, " ERR=", err)
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
