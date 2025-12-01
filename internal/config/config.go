package config

import "os"

type Config struct {
	MongoURL    string
	DBName      string
	JWTSecret   string
	CORSOrigins string
	YouTubeKey  string
}

func Load() Config {
	return Config{
		MongoURL:    os.Getenv("MONGO_URL"),
		DBName:      os.Getenv("DB_NAME"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		CORSOrigins: os.Getenv("CORS_ORIGINS"),
		YouTubeKey:  os.Getenv("YOUTUBE_API_KEY"),
	}
}

