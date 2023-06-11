package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/x1xo/Auth/src/databases"
	"github.com/x1xo/Auth/src/routes"
	callbackRoutes "github.com/x1xo/Auth/src/routes/callback"
	"github.com/x1xo/Auth/src/utils"
)

func main() {
	godotenv.Load()
	go databases.GetRedis()
	databases.GetMongo()

	app := fiber.New(fiber.Config{
		ProxyHeader: "X-Forwarded-For",
		EnableTrustedProxyCheck: false,
	})
	app.Use(logger.New())

	jwks, err := utils.GenerateJWKS()
	if err != nil {
		log.Println("Couldn't generate JWK Set for public use.")
		panic(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Identity provider by x1xo. All rights reserved.")
	})

	app.Get("/jwks", func(c *fiber.Ctx) error {
		return c.JSON(jwks)
	})

	app.Get("/api/user", routes.GetUser)
	app.Get("/api/user/sessions", routes.GetUserSessions)
	app.Get("/api/user/sessions/:tokenId", routes.GetUserSession)

	app.Get("/login", routes.Login)

	app.Get("/callback/github", callbackRoutes.CallbackGithub)
	app.Get("/callback/discord", callbackRoutes.CallbackDiscord)
	app.Get("/callback/google", callbackRoutes.CallbackGoogle)

	environment := os.Getenv("ENVIRONMENT")
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}
	if environment == "" {
		environment = "production"
	}

	if environment == "production" {
		log.Fatal(app.Listen(fmt.Sprintf("%s:%s", "0.0.0.0", port)))
	} else {
		log.Fatal(app.Listen(fmt.Sprintf("%s:%s", "127.0.0.1", port)))
	}
}
