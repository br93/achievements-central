package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Config struct{}

func main() {

	godotenv.Load(".env")
	port := os.Getenv("PORT")

	app := Config{}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	if port == "" {
		log.Fatal("PORT not found")
	}

	log.Printf("Starting auth service...")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
