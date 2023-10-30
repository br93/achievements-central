package main

import (
	models "auth-service/cmd/api/data/models"
	"database/sql"
)

const port = "8080"

type Models models.Models

type Config struct {
	DB     *sql.DB
	Models Models
}

func main() {

}
