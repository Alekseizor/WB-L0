package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const driver = "postgres"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("ИЗМЕНИ ЛОГИ")
	}

	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USERNAME")
	dbPort := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", dbUser, dbName, dbPass, dbHost, dbPort)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Println(fmt.Errorf("failed to connect to the db - %s", err.Error()))
		return
	}
	db.SetMaxOpenConns(0)
	err = db.Ping()
	if err != nil {
		log.Println(fmt.Errorf("failed to connect to the db - %s", err.Error()))
		return
	}

}
