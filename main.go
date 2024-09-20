package main

import (
	"github.com/joho/godotenv"
	"github.com/nickduskey/go-daisy/internal/server"
)

func main() {
	_ = godotenv.Load()
	server.Run()
}
