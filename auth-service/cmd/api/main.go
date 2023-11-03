package main

import (
	data "auth-service/cmd/api/data"
	"auth-service/internal/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	DB     *database.Queries
	Models data.Models
}

func main() {

	var env []string
	env = loadEnv()

	conn, err := sql.Open("postgres", env[1])
	if err != nil {
		log.Fatal("Cannot connect to DB!")
	}

	app := Config{
		DB: database.New(conn),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", env[0]),
		Handler: app.routes(),
	}

	log.Printf("Starting auth service...")

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func loadEnv() []string {
	godotenv.Load(".env")
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	if port == "" {
		log.Fatal("PORT not found")
	}

	if dbURL == "" {
		log.Fatal("DB URL not found")
	}

	var env []string
	env = append(env, port, dbURL)

	return env
}
