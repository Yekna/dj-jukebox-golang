package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"dj-jukebox/internal/config"
	"dj-jukebox/internal/database"
	"dj-jukebox/internal/models"
	"dj-jukebox/internal/utils"
)

func AuthRequired(cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			return fiber.ErrUnauthorized
		}

		token := parts[1]

		claims, err := utils.VerifyJWT(token, cfg.JWTSecret)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		userID, _ := (*claims)["user_id"].(string)

		var user models.User
		err = database.DB.Collection("users").
			FindOne(c.Context(), bson.M{"id": userID}).
			Decode(&user)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		c.Locals("user", user)
		return c.Next()
	}
}

