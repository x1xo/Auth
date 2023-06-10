package routes

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/x1xo/Auth/src/databases"
	"github.com/x1xo/Auth/src/databases/models"
	"github.com/x1xo/Auth/src/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUser(c *fiber.Ctx) error {
	var tokenHeader string
	header := c.Get("Authorization")
	if header != "" {
		tokenHeader = strings.Split(c.Get("Authorization"), " ")[1]
	}
	tokenCookie := c.Cookies("access_token")

	if tokenHeader == "" && tokenCookie == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": true,
			"code":  "UNAUTHORIZED",
		})
	}

	var tokenString string
	if tokenHeader != "" {
		tokenString = tokenHeader
	} else {
		tokenString = tokenCookie
	}

	key, err := jwk.ParseKey([]byte(*utils.PublicKey), jwk.WithPEM(true))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"code":  "INTERNAL_SERVER_ERROR",
		})
	}

	token, err := jwt.ParseString(tokenString, jwt.WithVerify(jwa.RS256, key))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_TOKEN",
		})
	}

	userId := token.Subject()
	if userId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_TOKEN",
		})
	}

	var user models.UserInfo
	databases.GetMongoDatabase().Collection("users").FindOne(context.Background(), bson.M{"id": userId}).Decode(&user)

	return c.Status(200).JSON(user)

}
