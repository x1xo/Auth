package routes

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/x1xo/Auth/src/databases"
	"github.com/x1xo/Auth/src/databases/models"
	"github.com/x1xo/Auth/src/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUser(c *fiber.Ctx) error {
	var tokenHeader string
	header := c.Get("Authorization")
	if header != "" {
		split := strings.Split(c.Get("Authorization"), " ")
		if len(split) < 2 {
			tokenHeader = ""
		} else {
			tokenHeader = strings.Split(c.Get("Authorization"), " ")[1]
		}
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

	token, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_TOKEN",
		})
	}

	var user models.UserInfo
	databases.GetMongoDatabase().Collection("users").FindOne(context.Background(), bson.M{"id": (*token).Subject()}).Decode(&user)

	return c.Status(200).JSON(user)
}

func GetUserSessions(c *fiber.Ctx) error {
	var tokenHeader string
	header := c.Get("Authorization")
	if header != "" {
		split := strings.Split(c.Get("Authorization"), " ")
		if len(split) < 2 {
			tokenHeader = ""
		} else {
			tokenHeader = strings.Split(c.Get("Authorization"), " ")[1]
		}
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

	token, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_TOKEN",
		})
	}

	keys := databases.GetRedis().Keys(context.Background(), (*token).Subject()+"*").Val()

	var sessions []models.UserSession

	jsonSessions, err := databases.GetRedis().MGet(context.Background(), keys...).Result()
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
	tokenId := c.Params("tokenId")
	if tokenId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_TOKEN_ID",
		})
	}

	var tokenHeader string
	header := c.Get("Authorization")
	if header != "" {
		split := strings.Split(c.Get("Authorization"), " ")
		if len(split) < 2 {
			tokenHeader = ""
		} else {
			tokenHeader = strings.Split(c.Get("Authorization"), " ")[1]
		}
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

	token, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_TOKEN",
		})
	}

	sessionJSON := databases.GetRedis().Get(context.Background(), (*token).Subject()+"_"+tokenId).Val()
	if sessionJSON == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_TOKEN",
		})
	}

	var session models.UserSession
	err = json.Unmarshal([]byte(sessionJSON), &session)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"code":  "INTERNAL_SERVER_ERROR",
		})
	}

	return c.Status(200).JSON(session)
}

func InvalidateSession(c *fiber.Ctx) error {
	tokenId := c.Params("tokenId")
	if tokenId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_TOKEN_ID",
		})
	}

	var tokenHeader string
	header := c.Get("Authorization")
	if header != "" {
		split := strings.Split(c.Get("Authorization"), " ")
		if len(split) < 2 {
			tokenHeader = ""
		} else {
			tokenHeader = strings.Split(c.Get("Authorization"), " ")[1]
		}
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

	token, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_TOKEN",
		})
	}

	err = databases.GetRedis().Del(context.Background(), (*token).Subject()+"_"+tokenId).Err()
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": true,
			"code":  "NOT_FOUND",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error": false,
		"code":  "OK",
	})
}

func InvalidateAllSessions(c *fiber.Ctx) error {
	var tokenHeader string
	header := c.Get("Authorization")
	if header != "" {
		split := strings.Split(c.Get("Authorization"), " ")
		if len(split) < 2 {
			tokenHeader = ""
		} else {
			tokenHeader = strings.Split(c.Get("Authorization"), " ")[1]
		}
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

	token, err := utils.ValidateToken(tokenString)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_TOKEN",
		})
	}

	keys := databases.GetRedis().Keys(context.Background(), (*token).Subject()+"*").Val()

	err = databases.GetRedis().Del(context.Background(), keys...).Err()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"code":  "INTERNAL_SERVER_ERROR",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error": false,
		"code":  "OK",
	})
}
