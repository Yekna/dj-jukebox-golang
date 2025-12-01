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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RequestSong(manager *websocket.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pin := c.Params("pin")

		// get room
		var room models.Room
		err := database.DB.Collection("rooms").
			FindOne(c.Context(), bson.M{"pin": pin, "active": true}).
			Decode(&room)

		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Room not found or closed")
		}

		var req models.SongRequestCreate
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		song := models.SongRequest{
			ID:            utils.NewUUID(),
			RoomID:        room.ID,
			YouTubeID:     req.YouTubeID,
			Title:         req.Title,
			Thumbnail:     req.Thumbnail,
			URL:           req.URL,
			SubmitterName: req.SubmitterName,
			SubmitterType: req.SubmitterType,
			Votes:         0,
			VotedBy:       []string{},
			Status:        "pending",
			CreatedAt:     time.Now().UTC(),
		}

		_, err = database.DB.Collection("song_requests").
			InsertOne(c.Context(), song)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		manager.Broadcast(pin, fiber.Map{
			"type": "song_requested",
			"song": song.Title,
			"user": song.SubmitterName,
		})

		return c.JSON(song)
	}
}

func GetRoomSongs() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pin := c.Params("pin")

		var room models.Room
		err := database.DB.Collection("rooms").
			FindOne(c.Context(), bson.M{"pin": pin, "active": true}).
			Decode(&room)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Room not found or closed")
		}

		cursor, err := database.DB.Collection("song_requests").
			Find(c.Context(),
				bson.M{"room_id": room.ID},
				options.Find().SetSort(bson.M{"created_at": 1}),
			)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		var songs []models.SongRequest
		if err := cursor.All(c.Context(), &songs); err != nil {
			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{"songs": songs})
	}
}

func VoteSong(manager *websocket.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		songID := c.Params("song_id")

		var req models.VoteRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		var song models.SongRequest
		err := database.DB.Collection("song_requests").
			FindOne(c.Context(), bson.M{"id": songID}).
			Decode(&song)

		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Song not found")
		}

		roomID := song.RoomID
		var room models.Room
		_ = database.DB.Collection("rooms").
			FindOne(c.Context(), bson.M{"id": roomID}).
			Decode(&room)

		// toggle vote
		already := false
		for _, v := range song.VotedBy {
			if v == req.SessionID {
				already = true
				break
			}
		}

		update := bson.M{}

		if already {
			update = bson.M{
				"$inc":  bson.M{"votes": -1},
				"$pull": bson.M{"voted_by": req.SessionID},
			}
		} else {
			update = bson.M{
				"$inc":  bson.M{"votes": 1},
				"$push": bson.M{"voted_by": req.SessionID},
			}
		}

		_, err = database.DB.Collection("song_requests").
			UpdateOne(c.Context(), bson.M{"id": songID}, update)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		// return updated doc
		_ = database.DB.Collection("song_requests").
			FindOne(c.Context(), bson.M{"id": songID}).
			Decode(&song)

		// broadcast
		if room.Pin != "" {
			manager.Broadcast(room.Pin, fiber.Map{
				"type": "song_voted",
				"song": song,
			})
		}

		return c.JSON(song)
	}
}

func UpdateSongStatus(manager *websocket.Manager, cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		songID := c.Params("song_id")
		user := c.Locals("user").(models.User)

		var song models.SongRequest
		err := database.DB.Collection("song_requests").
			FindOne(c.Context(), bson.M{"id": songID}).
			Decode(&song)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Song not found")
		}

		var room models.Room
		err = database.DB.Collection("rooms").
			FindOne(c.Context(), bson.M{"id": song.RoomID}).
			Decode(&room)
		if err != nil {
			return fiber.ErrForbidden
		}

		if room.DJID != user.ID {
			return fiber.ErrForbidden
		}

		var req models.StatusUpdate
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		_, err = database.DB.Collection("song_requests").
			UpdateOne(c.Context(), bson.M{"id": songID},
				bson.M{"$set": bson.M{"status": req.Status}},
			)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		_ = database.DB.Collection("song_requests").
			FindOne(c.Context(), bson.M{"id": songID}).
			Decode(&song)

		manager.Broadcast(room.Pin, fiber.Map{
			"type": "song_status_changed",
			"song": song,
		})

		return c.JSON(song)
	}
}

func SearchSongs(cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		q := c.Query("q")
		max := c.QueryInt("max_results", 10)

		results, err := utils.SearchYouTube(q, cfg.YouTubeKey, max)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Search failed")
		}

		return c.JSON(fiber.Map{"results": results})
	}
}
