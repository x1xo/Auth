package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lestrrat-go/jwx/jwa"
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

	token, err := jwt.ParseString(tokenString, jwt.WithVerify(jwa.RS256, *utils.PublicJWTKey))
	if err != nil {
		fmt.Println("Invalid token here:", err)
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

func GetUserSessions(c *fiber.Ctx) error {
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

	token, err := jwt.ParseString(tokenString, jwt.WithVerify(jwa.RS256, *utils.PublicJWTKey))
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

	keys := databases.GetRedis().Keys(context.Background(), userId+"*").Val()

	var sessions []models.UserSession

	jsonSessions, err := databases.GetRedis().MGet(context.Background(), keys...).Result();
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"code":  "INTERNAL_SERVER_ERROR",
		})
	}

	for _, jsonSession := range jsonSessions {
		var session models.UserSession
		err = json.Unmarshal([]byte(jsonSession.(string)), &session)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": true,
				"code":  "INTERNAL_SERVER_ERROR",
			})
		}

		sessions = append(sessions, session)
	}

	return c.Status(200).JSON(sessions)
}



func GetUserSession(c *fiber.Ctx) error {
	tokenId := c.Params("tokenId");
	if tokenId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_TOKEN_ID",
		})
	}

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

	token, err := jwt.ParseString(tokenString, jwt.WithVerify(jwa.RS256, *utils.PublicJWTKey))
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

	sessionJSON := databases.GetRedis().Get(context.Background(), userId+"_"+tokenId).Val();
	if sessionJSON == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_TOKEN",
		})
	}

	var session models.UserSession;
	err = json.Unmarshal([]byte(sessionJSON), &session);
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"code":  "INTERNAL_SERVER_ERROR",
		})
	}

	return c.Status(200).JSON(session);

}
