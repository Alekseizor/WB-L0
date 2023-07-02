package gettingstream

import (
	"WB-L0/internal/pkg/repository/delivery"
	"WB-L0/internal/pkg/repository/items"
	"WB-L0/internal/pkg/repository/orders"
	"WB-L0/internal/pkg/repository/payment"
	"context"
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

type ClientNatsStreaming struct {
	PostgresOrderRepo orders.OrderRepo
	InMemoryOrderRepo orders.OrderInMemoryRepo
	DeliveryRepo      delivery.DeliveryRepo
	ItemsRepo         items.ItemRepo
	PaymentRepo       payment.PaymentRepo
}

func (c ClientNatsStreaming) ReceivingOrder() {
	// Настройка параметров подключения
	natsURL := "nats://localhost:4222" // URL сервера NATS
	clusterID := "test-cluster"        // Идентификатор кластера NATS Streaming
	clientID := "client-1"             // Идентификатор клиента

	// Создание соединения с сервером NATS Streaming
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Println(err)
	}

	// Отложенное закрытие соединения при завершении работы
	defer sc.Close()

	// Подписка на канал
	channel := "my-channel"
	var orderAllData orders.OrderAllData
	var order orders.Order
	ctx := context.Background()
	_, err = sc.Subscribe(channel, func(msg *stan.Msg) {
		// Обработка полученного сообщения

		log.Printf("Получено сообщение: %s\n", string(msg.Data))
		err = json.Unmarshal(msg.Data, &orderAllData)
		if err != nil {
			log.Println("json err: ", err)
			return
		}
		deliveryUUID, err := c.DeliveryRepo.AddDelivery(ctx, orderAllData.Delivery)
		if err != nil {
			log.Println("DeliveryRepo error:", err)
			return
		}
		paymentUUID, err := c.PaymentRepo.AddPayment(ctx, orderAllData.Payment)
		if err != nil {
			log.Println("PaymentRepo error:", err)
			return
		}
		itemsUUID, err := c.ItemsRepo.AddItems(ctx, orderAllData.Items)
		if err != nil {
			log.Println("PaymentRepo error:", err)
			return
		}
		order.OrderUID = orderAllData.OrderUID
		order.TrackNumber = orderAllData.TrackNumber
		order.Entry = orderAllData.Entry
		order.Locale = orderAllData.Locale
		order.InternalSignature = orderAllData.InternalSignature
		order.CustomerID = orderAllData.CustomerID
		order.DeliveryService = orderAllData.DeliveryService
		order.Shardkey = orderAllData.Shardkey
		order.SmID = orderAllData.SmID
		order.DateCreated = orderAllData.DateCreated
		order.OofShard = orderAllData.OofShard
		order.Delivery = *deliveryUUID
		order.Payment = *paymentUUID
		order.Items = itemsUUID
		err = c.PostgresOrderRepo.AddOrder(ctx, order)
		if err != nil {
			log.Println("PostgresOrderRepo error:", err)
			return
		}
		err = c.InMemoryOrderRepo.AddOrder(ctx, orderAllData)
		if err != nil {
			log.Println("InMemoryOrderRepo error:", err)
			return
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
