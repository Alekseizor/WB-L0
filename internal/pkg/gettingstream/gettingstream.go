package gettingstream

import (
	"WB-L0/internal/pkg/repository/orders"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

type ClientNatsStreaming struct {
	PostgresOrderRepo orders.OrderRepo
	InMemoryOrderRepo orders.OrderRepo
	DeliveryRepo      orders.OrderRepo
	ItemsRepo         orders.OrderRepo
	PaymentRepo       orders.OrderRepo
}

func (c ClientNatsStreaming) ReceivingOrder() {
	// Настройка параметров подключения
	natsURL := "nats://localhost:4222" // URL сервера NATS
	clusterID := "test-cluster"        // Идентификатор кластера NATS Streaming
	clientID := "client-1"             // Идентификатор клиента

	// Создание соединения с сервером NATS Streaming
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatal(err)
	}

	// Отложенное закрытие соединения при завершении работы
	defer sc.Close()

	// Подписка на канал
	channel := "my-channel"
	var order orders.OrderAllData
	_, err = sc.Subscribe(channel, func(msg *stan.Msg) {
		// Обработка полученного сообщения
		log.Printf("Получено сообщение: %s\n", string(msg.Data))
		err = json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Println("json err: ", err)
		}

		err := c.PostgresRepo.AddOrder(order)
		if err != nil {
			log.Println("PostgresRepo error:", err)
		}
		err = c.InMemoryRepo.AddOrder(order)
		if err != nil {
			log.Println("InMemoryRepo error:", err)
		}
		// Подтверждение получения сообщения
		msg.Ack()
	}, stan.SetManualAckMode())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Подписка на канал '%s'...", channel)

	//Ожидание событий
	for {
		time.Sleep(1 * time.Second)
	}
}
