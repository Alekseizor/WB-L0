package main

import (
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

func main() {
	// Настройка параметров подключения к серверу NATS
	natsURL := "nats://localhost:4222" // URL сервера NATS

	// Настройка параметров подключения к кластеру NATS Streaming
	clusterID := "test-cluster" // Идентификатор кластера NATS Streaming
	clientID := "server"        // Идентификатор сервера

	// Создание соединения с кластером NATS Streaming
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatal(err)
	}

	// Отложенное закрытие соединений при завершении работы
	defer sc.Close()

	// Создание NATS Streaming канала
	channel := "my-channel"

	// Публикация сообщений в канале
	go func() {
		for i := 0; ; i++ {
			message := []byte("Сообщение " + string(i+1))
			if err := sc.Publish(channel, message); err != nil {
				log.Printf("Ошибка публикации сообщения: %v\n", err)
			} else {
				log.Printf("Опубликовано сообщение: %s\n", message)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	log.Printf("Сервер ожидает сообщений на канале '%s'...", channel)

	// Ожидание событий
	for {
		time.Sleep(1 * time.Second)
	}
}
