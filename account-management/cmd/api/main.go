package main

import (
	"account-management/cmd/api/services"
	"account-management/internal/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Services Services
}

func main() {

	godotenv.Load()

	conn, err := sql.Open("postgres", getEnv("DB_URL"))
	if err != nil {
		log.Fatal("Cannot connect to DB!")
	}

	rep := database.New(conn)

	app := Config{
		Services{
			User:        services.NewUserService(rep),
			UserProfile: services.NewUserProfileService(rep),
		}}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", getEnv("PORT")),
		Handler: app.routes(),
	}

	log.Printf("Starting account management...")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func getEnv(key string) string {

	value := os.Getenv(key)

	if value == "" {
		log.Fatal(key + " not found")
	}

	return value

}
