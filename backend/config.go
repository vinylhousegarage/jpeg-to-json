package main

import "os"

type Config struct {
	AllowedOrigin string
}

func LoadConfig() *Config {
	origin := os.Getenv("ALLOWED_ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000" 
	}

	return &Config{
		AllowedOrigin: origin,
	}
}
