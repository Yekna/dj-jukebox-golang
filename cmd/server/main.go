package main

import (
	"dj-jukebox/internal/config"
	"dj-jukebox/internal/database"
	"dj-jukebox/internal/routes"
	"embed"
	"io/fs"
	"log"
	"mime"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

//go:embed dist
var embeddedFiles embed.FS

func main() {
	godotenv.Load(".env")

	app := fiber.New()

	// load config
	cfg := config.Load()

	// connect database
	if err := database.Connect(cfg.MongoURL, cfg.DBName); err != nil {
		log.Fatal("Mongo connect error:", err)
	}

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowCredentials: true,
		AllowMethods:     "GET,POST,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Authorization,Content-Type",
	}))

	// register routes
	routes.Register(app, cfg)

	fsys, _ := fs.Sub(embeddedFiles, "dist")

	app.All("/*", func(c *fiber.Ctx) error {
		path := c.Path()[1:]
		if path == "" {
			path = "index.html"
		}

		f, err := fsys.Open(path)
		if err != nil {
			// fallback to index.html for SPA routes
			f, err = fsys.Open("index.html")
			if err != nil {
				return fiber.ErrNotFound
			}
			path = "index.html"
		}
		defer f.Close()

		// Set MIME type based on file extension
		ext := filepath.Ext(path)
		if ext != "" {
			mimeType := mime.TypeByExtension(ext)
			if mimeType != "" {
				c.Set("Content-Type", mimeType)
			}
		}

		return c.SendStream(f)
	})

	log.Println("Server running on :8080")
	log.Fatal(app.Listen(":8080"))
}
