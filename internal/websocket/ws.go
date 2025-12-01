package websocket

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WSHandler(manager *Manager) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		room := c.Params("room_pin")

		manager.Connect(room, c)
		defer manager.Disconnect(room, c)

		// initial message
		c.WriteJSON(fiber.Map{
			"type":     "connected",
			"room_pin": room,
		})

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				return
			}

			var data map[string]interface{}
			json.Unmarshal(msg, &data)

			if data["type"] == "user_joined" {
				manager.Broadcast(room, fiber.Map{
					"type": "user_joined",
					"user": data["user"],
				})
			}
		}
	}
}
