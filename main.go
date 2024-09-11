package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/muhammadolammi/tinydb/server"
)

func main() {
	// writing this as a single node client-server will add peers later (multi node cluster)
	err := godotenv.Load()
	if err != nil {
		log.Panic("error loading env")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "6379"
	}

	config := Config{
		PORT: port,
	}
	log.Println("starting server")
	server := server.NewServer(":" + config.PORT)
	log.Fatal(server.Start())
}
