package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/jackc/pgx/v4"
	"github.com/matdorneles/go_microservices/authentication-service/data"
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres !!")
	}

	// set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// run web server
	app.routes()
}
