package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ansrivas/fiberprometheus/v2"
	// "github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/golang-jwt/jwt/v4"
)

var KeycloakHost = os.Getenv("KEYCLOAK_HOST")
var Backend1Host = os.Getenv("BACKEND_API_1_HOST")
var Backend2Host = os.Getenv("BACKEND_API_2_HOST")
var ClientID = os.Getenv("CLIENT_ID")
var ClientSecret = os.Getenv("CLIENT_SECRET")

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

type Token struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	Scope            string `json:"scope"`
}

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

	profile := app.Group("profile", keyauth.New(keyauth.Config{
		Validator: validator,
	}))
	profile.Get("name", func(c *fiber.Ctx) error {
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

	backend2 := app.Group("backend2", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	backend2.Get("/pets/list", listPets)
	backend2.Get("/profile/name", getProfileName)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from backend1!")
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

func listPets(c *fiber.Ctx) error {
	token, err := getAccessToken(c)
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusInternalServerError).SendString("An error ocurred, check logs")
	}
	reponse, err := getBackendApi2(c, token, "/protected/pets/list")
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusInternalServerError).SendString("An error ocurred, check logs")
	}
	return c.Status(fiber.StatusOK).SendString(string(reponse))
}

func getProfileName(c *fiber.Ctx) error {
	token, err := getAccessToken(c)
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusInternalServerError).SendString("An error ocurred, check logs")
	}
	reponse, err := getBackendApi2(c, token, "/protected/profile/name")
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusInternalServerError).SendString("An error ocurred, check logs")
	}
	return c.Status(fiber.StatusOK).SendString(string(reponse))
}

func getAccessToken(c *fiber.Ctx) (string, error) {

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", ClientID)
	data.Set("client_secret", ClientSecret)

	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, KeycloakHost+"/realms/master/protocol/openid-connect/token", strings.NewReader(data.Encode())) // URL-encoded payload
	// r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(r)
	fmt.Println(resp.Status)
	reponseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Could not get body")
	}

	var ResponseToken Token
	err113 := json.Unmarshal(reponseData, &ResponseToken)
	if err113 != nil {
		return "", errors.New("Could not unmarshal body")
	}

	return ResponseToken.AccessToken, nil
}

func getBackendApi2(c *fiber.Ctx, token string, path string) (string, error) {
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodGet, Backend2Host+path, bytes.NewReader([]byte(``)))
	r.Header.Add("Authorization", "Bearer "+token)
	resp, _ := client.Do(r)
	log.Println(resp.Status, resp.Request.URL)
	responseData, err := ioutil.ReadAll(resp.Body)
	log.Println(string(responseData))
	if err != nil {
		log.Fatal(err)

		return "", errors.New("Could not get a body")
	}
	return string(responseData), nil
}
