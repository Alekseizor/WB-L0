package main

import (
	"WB-L0/internal/pkg/handlers"
	"WB-L0/internal/pkg/sendingjson"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"net/http"
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
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Println(fmt.Errorf("couldn't create a new logger - %s", err.Error()))
		return
	}
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	serviceSend := sendingjson.NewServiceSend(logger)

	orderHandler := &handlers.OrdersHandler{
		Logger: logger,
		Send:   serviceSend,
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/orders/{ID}", orderHandler.GetOrderByID).Methods("GET")
	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		logger.Errorf("couldn't start listening - %s", err.Error())
	}
}
