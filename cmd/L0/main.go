package main

import (
	"WB-L0/internal/pkg/gettingstream"
	"WB-L0/internal/pkg/handlers"
	"WB-L0/internal/pkg/repository/delivery"
	"WB-L0/internal/pkg/repository/items"
	"WB-L0/internal/pkg/repository/orders"
	"WB-L0/internal/pkg/repository/payment"
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

	serviceSend := sendingjson.NewServiceSendJSON(logger)
	inMemoryOrderRepo, err := orders.NewRepoOrderInMemory()
	if err != nil {
		logger.Errorf("failed to create NewRepoOrderInMemory - %v", err)
	}
	orderHandler := &handlers.OrdersHandler{
		OrderRepo: inMemoryOrderRepo,
		Logger:    logger,
		Send:      serviceSend,
	}

	deliveryRepo, err := delivery.NewRepoDeliveryPostgres(db)
	if err != nil {
		logger.Errorf("failed to create NewRepoDeliveryPostgres - %v", err)
	}
	itemsRepo, err := items.NewRepoItemsPostgres(db)
	if err != nil {
		logger.Errorf("failed to create NewRepoItemsPostgres - %v", err)
	}
	orderRepo, err := orders.NewRepoOrderPostgres(db)
	if err != nil {
		logger.Errorf("failed to create NewRepoOrderPostgres - %v", err)
	}
	paymentRepo, err := payment.NewRepoPaymentPostgres(db)
	if err != nil {
		logger.Errorf("failed to create NewRepoOrderPostgres - %v", err)
	}
	clientNats := &gettingstream.ClientNatsStreaming{
		PostgresOrderRepo: orderRepo,
		InMemoryOrderRepo: inMemoryOrderRepo,
		DeliveryRepo:      deliveryRepo,
		ItemsRepo:         itemsRepo,
		PaymentRepo:       paymentRepo,
	}
	go clientNats.ReceivingOrder()
	r := mux.NewRouter()
	staticFiles := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(staticFiles)
	//r.HandleFunc("/orders", orderHandler.GetOrderByID).Methods("GET")
	r.HandleFunc("/api/orders/{ID}", orderHandler.GetOrderByID).Methods("GET")
	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		logger.Errorf("couldn't start listening - %v", err)
	}
}
