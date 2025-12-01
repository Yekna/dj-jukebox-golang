package handlers

import (
	"time"

	"dj-jukebox/internal/config"
	"dj-jukebox/internal/database"
	"dj-jukebox/internal/models"
	"dj-jukebox/internal/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.UserRegister
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		// Check if exists
		count, _ := database.DB.Collection("users").
			CountDocuments(c.Context(), bson.M{"email": req.Email})
		if count > 0 {
			return fiber.NewError(fiber.StatusBadRequest, "Email already registered")
		}

		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		user := models.User{
			ID:           utils.NewUUID(),
			Email:        req.Email,
			PasswordHash: hash,
			CreatedAt:    time.Now().UTC(),
		}

		_, err = database.DB.Collection("users").InsertOne(c.Context(), user)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		token, err := utils.CreateJWT(user.ID, user.Email, cfg.JWTSecret)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":    user.ID,
				"email": user.Email,
			},
		})
	}
}

func Login(cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.UserLogin
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		var user models.User
		err := database.DB.Collection("users").
			FindOne(c.Context(), bson.M{"email": req.Email}).
			Decode(&user)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		if !utils.VerifyPassword(req.Password, user.PasswordHash) {
			return fiber.ErrUnauthorized
		}

		token, err := utils.CreateJWT(user.ID, user.Email, cfg.JWTSecret)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":    user.ID,
				"email": user.Email,
			},
		})
	}
}
