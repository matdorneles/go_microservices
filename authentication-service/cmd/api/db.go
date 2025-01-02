package main

import (
	"database/sql"
	"log"
	"os"
	"time"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	counts := 0

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres !!")
			return connection
		}

		if counts > 20 {
			log.Println(err)
			return nil
		}

		log.Println("Waiting for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
