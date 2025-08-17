package main

import (
	"fmt"
	"go/frame/mq/rabbitmq"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main5() {
	//连接RabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmq.User, rabbitmq.Pass, rabbitmq.Host, rabbitmq.Port))
	if err != nil {
		log.Panicf("connect to RabbitMQ failed: %s", err)
	}
	defer conn.Close()

	log.Printf("waiting for messages, to exit press CTRL+C")
	go SubscribeByKey(conn, 1, rabbitmq.ExchangeName3, "*.info")
	go SubscribeByKey(conn, 2, rabbitmq.ExchangeName3, "machine1.*")
	go SubscribeByKey(conn, 3, rabbitmq.ExchangeName3, "*.*")
	select {}
}

// go run .\mq\rabbitmq\consumer\
