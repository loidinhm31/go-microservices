package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"github.com/loidinhm31/go-microservice/common"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var counts int64

type Config struct {
	Repo   data.Repository
	Client *http.Client
}

func main() {
	log.Printf("Starting authentication service on port %s\n", common.AuthPort)

	// Connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Cannot connect to database")
	}

	db := data.NewPostgresRepository(conn)

	// set up config
	app := Config{
		Client: &http.Client{},
		Repo:   db,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", common.AuthPort),
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

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

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Database not yet ready...")
			counts++
		} else {
			log.Println("Connected to database")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
