package main

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	ServerAddr  string
}

func LoadConfig() Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	return Config{
		DatabaseURL: dbURL,
		ServerAddr:  addr,
	}
}
