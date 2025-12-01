package handlers

import (
	"time"

	"dj-jukebox/internal/config"
	"dj-jukebox/internal/database"
	"dj-jukebox/internal/models"
	"dj-jukebox/internal/utils"
	"dj-jukebox/internal/websocket"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateRoom(cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.User)

		// Existing active room?
		var existing models.Room
		err := database.DB.Collection("rooms").
			FindOne(c.Context(), bson.M{"dj_id": user.ID, "active": true}).
			Decode(&existing)

		if err == nil {
			return c.JSON(existing)
		}

		// Generate unique PIN
		pin := utils.GenerateRoomPin()
		for {
			count, _ := database.DB.Collection("rooms").
				CountDocuments(c.Context(), bson.M{"pin": pin, "active": true})
			if count == 0 {
				break
			}
			pin = utils.GenerateRoomPin()
		}

		room := models.Room{
			ID:        utils.NewUUID(),
			Pin:       pin,
			DJID:      user.ID,
			DJEmail:   user.Email,
			Active:    true,
			CreatedAt: time.Now().UTC(),
		}

		_, err = database.DB.Collection("rooms").InsertOne(c.Context(), room)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(room)
	}
}

func GetRoom() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pin := c.Params("pin")

		var room models.Room
		err := database.DB.Collection("rooms").
			FindOne(c.Context(), bson.M{"pin": pin, "active": true}).
			Decode(&room)

		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Room not found or closed")
		}

		return c.JSON(room)
	}
}

func CloseRoom(manager *websocket.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pin := c.Params("pin")
		user := c.Locals("user").(models.User)

		var room models.Room
		err := database.DB.Collection("rooms").
			FindOne(c.Context(), bson.M{"pin": pin, "active": true}).
			Decode(&room)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Room not found")
		}

		if room.DJID != user.ID {
			return fiber.ErrForbidden
		}

		_, err = database.DB.Collection("rooms").
			UpdateOne(c.Context(), bson.M{"pin": pin}, bson.M{"$set": bson.M{"active": false}})
		if err != nil {
			return fiber.ErrInternalServerError
		}

		manager.Broadcast(pin, fiber.Map{
			"type":    "room_closed",
			"message": "DJ has closed the room",
		})

		return c.JSON(fiber.Map{"message": "Room closed successfully"})
	}
}
