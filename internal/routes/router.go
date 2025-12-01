package routes

import (
	"dj-jukebox/internal/config"
	"dj-jukebox/internal/handlers"
	"dj-jukebox/internal/middleware"
	"dj-jukebox/internal/websocket"

	"github.com/gofiber/fiber/v2"
	wsfiber "github.com/gofiber/websocket/v2"
)

func Register(app *fiber.App, cfg config.Config) {
	manager := websocket.NewManager()

	// Public API group
	public := app.Group("/api")

	// WebSocket (public)
	public.Get("/ws/:room_pin", wsfiber.New(func(c *wsfiber.Conn) {
		websocket.WSHandler(manager)(c)
	}))

	// Auth
	public.Post("/auth/register", handlers.Register(cfg))
	public.Post("/auth/login", handlers.Login(cfg))

	// Public Songs
	public.Get("/songs/search", handlers.SearchSongs(cfg))
	public.Post("/songs/:song_id/vote", handlers.VoteSong(manager))

	// Public Rooms (read-only)
	public.Get("/rooms/:pin", handlers.GetRoom())
	public.Get("/rooms/:pin/songs", handlers.GetRoomSongs())
	public.Post("/rooms/:pin/songs", handlers.RequestSong(manager))

	// Protected API group
	protected := app.Group("/api", middleware.AuthRequired(cfg))
	protected.Post("/rooms/create", handlers.CreateRoom(cfg))
	protected.Post("/rooms/:pin/close", handlers.CloseRoom(manager))
	protected.Patch("/songs/:song_id/status", handlers.UpdateSongStatus(manager, cfg))
}

// func Register(app *fiber.App, cfg config.Config) {
// 	manager := websocket.NewManager()

// 	api := app.Group("/api")

// 	// Auth
// 	api.Post("/auth/register", handlers.Register(cfg))
// 	api.Post("/auth/login", handlers.Login(cfg))

// 	// Rooms
// 	api.Get("/rooms/:pin", handlers.GetRoom())
// 	api.Get("/rooms/:pin/songs", handlers.GetRoomSongs())
// 	api.Post("/rooms/:pin/songs", handlers.RequestSong(manager))

// 	secured := api.Use(middleware.AuthRequired(cfg))

// 	secured.Post("/rooms/create", handlers.CreateRoom(cfg))
// 	secured.Post("/rooms/:pin/close", handlers.CloseRoom(manager))
// 	secured.Patch("/songs/:song_id/status", handlers.UpdateSongStatus(manager, cfg))

// 	// Songs
// 	api.Post("/songs/:song_id/vote", handlers.VoteSong(manager))
// 	api.Get("/songs/search", handlers.SearchSongs(cfg))

// 	// WebSocket
// 	app.Get("/api/ws/:room_pin", wsfiber.New(func(c *wsfiber.Conn) {
// 		websocket.WSHandler(manager)(c)
// 	}))
// }
