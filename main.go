package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/muhammadolammi/tinydb/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("error loading env")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}

	config := Config{
		PORT: port,
	}
	log.Println("starting server")
	server := server.NewServer(":" + config.PORT)
	log.Fatal(server.Start())
}
