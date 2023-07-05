package main

import (
	"WB-L0/internal/pkg/repository/delivery"
	"WB-L0/internal/pkg/repository/items"
	"WB-L0/internal/pkg/repository/orders"
	"WB-L0/internal/pkg/repository/payment"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

const layout = "2006-01-02T15:04:05Z"

func main() {
	// Настройка параметров подключения к серверу NATS
	natsURL := "nats://localhost:4222" // URL сервера NATS

	// Настройка параметров подключения к кластеру NATS Streaming
	clusterID := "test-cluster" // Идентификатор кластера NATS Streaming
	clientID := "server"        // Идентификатор сервера

	// Создание соединения с кластером NATS Streaming
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Println(err)
	}

	// Отложенное закрытие соединений при завершении работы
	defer sc.Close()

	// Создание NATS Streaming канала
	channel := "my-channel"

	dateStr := "2021-11-26T06:22:19Z"

	dateOrder, err := time.Parse(layout, dateStr)
	if err != nil {
		log.Println("Ошибка при разборе даты:", err)
		return
	}
	order := orders.OrderAllData{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: delivery.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: payment.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []items.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
			{
				ChrtID:      12312321,
				TrackNumber: "sdsdk",
				Price:       532,
				RID:         "asdasd",
				Name:        "Mascaras",
				Sale:        33,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vsdso",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       dateOrder,
		OofShard:          "1",
	}
	// Публикация сообщений в канале
	go func() {
		for i := 0; ; i++ {

			message, err := json.Marshal(order)
			if err != nil {
				log.Println(err)
				continue
			}
			if err := sc.Publish(channel, message); err != nil {
				log.Printf("Ошибка публикации сообщения: %v\n", err)
			} else {
				log.Printf("Опубликовано сообщение: %s\n", message)
			}

			time.Sleep(10 * time.Second)
		}
	}()

	log.Printf("Сервер ожидает сообщений на канале '%s'...", channel)

	// Ожидание событий
	for {
		time.Sleep(1 * time.Second)
	}
}
