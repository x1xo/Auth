package callbackRoutes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/x1xo/Auth/src/databases"
	"github.com/x1xo/Auth/src/databases/models"
	"github.com/x1xo/Auth/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CallbackGoogle(c *fiber.Ctx) error {
	state := c.Query("state", "")
	code := c.Query("code", "")

	result, err := databases.GetRedis().Get(context.Background(), state).Result()
	if err != nil {
		log.Println("[Error] Couldn't get state from redis: \n", err)
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_STATE",
		})
	}
	if result == "" || result != "google" {
		return c.Status(400).JSON(fiber.Map{
			"error": true,
			"code":  "INVALID_STATE",
		})
	}

	databases.GetRedis().Del(context.Background(), state)

	response, err := getGoogleResponse(code)

	if err != nil {
		log.Println("[Error] Couldn't get response from google: \n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"code":  "INTERNAL_SERVER_ERROR",
		})
	}

	userInfo, err := getGoogleUserInfo(response)
	if err != nil {
		log.Println("[Error] Couldn't get user info from google: \n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"code":  "INTERNAL_SERVER_ERROR",
		})
	}

	db := databases.GetMongoDatabase()

	var user models.UserInfo
	err = db.Collection("users").FindOne(context.Background(), bson.M{"email": userInfo.Email}).Decode(&user)
	if err != nil && err == mongo.ErrNoDocuments {
		user = models.UserInfo{
			Id:        uuid.New().String(),
			Email:     userInfo.Email,
			Username:  userInfo.Username,
			AvatarURL: userInfo.AvatarURL,
			Google:    *userInfo,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		db.Collection("users").InsertOne(context.Background(), &user)
	}

	user.Google = *userInfo
	user.UpdatedAt = time.Now()

	go func() { db.Collection("users").ReplaceOne(context.Background(), bson.M{"id": user.Id}, user) }()

	token, tokenId, err := utils.GenerateToken(user.Id)
	if err != nil {
		log.Println("[Error] Couldn't generate token: \n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"code":  "INTERNAL_SERVER_ERROR",
		})
	}

	tokenJson, _ := json.Marshal(fiber.Map{
		"ip":        c.IP(),
		"userAgent": string(c.Context().UserAgent()),
		"createdAt": time.Now().UnixMilli(),
		"expiresAt": time.Now().Add(time.Hour * 3).UnixMilli(),
	})

	err = databases.GetRedis().Set(context.Background(), tokenId, tokenJson, time.Hour*3).Err()

	if err != nil {
		log.Println("[Error] Couldn't set token in redis: \n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": true,
			"code":  "INTERNAL_SERVER_ERROR",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 3),
		HTTPOnly: true,
		Secure:   true,
	})

	c.Set("Authorization", "Bearer "+token)

	return c.Redirect(os.Getenv("REDIRECT_URL"))
}

type GoogleAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Error       string `json:"error"`
}

func getGoogleUserInfo(data *GoogleAccessTokenResponse) (*models.GoogleUser, error) {

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + data.AccessToken)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var userInfo models.GoogleUser
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func getGoogleResponse(code string) (*GoogleAccessTokenResponse, error) {
	body := url.Values{
		"client_id":     {os.Getenv("GOOGLE_CLIENT_ID")},
		"client_secret": {os.Getenv("GOOGLE_CLIENT_SECRET")},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {os.Getenv("CALLBACK_URL") + "/callback/google"},
	}.Encode()

	request, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", bytes.NewBuffer([]byte(body)))

	if err != nil {
		log.Println("[Error] Couldn't create request for google callback.", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Println("[Error] Couldn't make request for google callback.", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Println("[Error] Response status from google callback is not 200.")
		return nil, errors.New("response status from google callback is not 200")
	}

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseMap GoogleAccessTokenResponse
	err = json.Unmarshal(responseBody, &responseMap)
	if err != nil {
		log.Println("[Error] Couldn't unmarshal json body for github callback.", err)
		return nil, err
	}

	if responseMap.Error != "" {
		log.Println("[Error] Invalid code was passed to the request somehow.")
		return nil, errors.New("invalid code was passed to the request somehow")
	}

	return &responseMap, nil
}
