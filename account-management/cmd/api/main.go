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

	app := Config{}
	godotenv.Load()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.getEnv("PORT")),
		Handler: app.routes(),
	}

	log.Printf("Starting account management...")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *Config) getEnv(key string) string {

	value := os.Getenv(key)

	if value == "" {
		log.Fatal(key + " not found")
	}

	return value

}
