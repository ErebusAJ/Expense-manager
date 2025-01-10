package utils

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/lib/pq"
)

// ConnectDB function to connect to database
// Retrieves the DB URL from .env
// Opens the connection
// Tests by pinging the DB
// Returns DB or err if any
func ConnectDB() (*sql.DB, error){

	//Retrieves database URL
	db_url := os.Getenv("DB_URL")
	if db_url == "" {
		log.Println("error retrieving DB_URL")
		return nil, nil
	}

	//Connect to DB
	DB, err := sql.Open("postgres", db_url)
	if err != nil{
		log.Println("error connecting to DB")
		return nil, err
	}

	//Ping DB
	if err := DB.Ping(); err != nil{
		log.Println("error pinging DB")
		return nil, err
	}

	log.Println("connection with DB established")

	return DB, nil;
}