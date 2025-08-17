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

	//创建Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("open channel failed: %s", err)
	}
	defer ch.Close()

	//声明Exchange。如果Exchange不存在会创建它；如果Exchange已存在，Server会检查声明的参数和Exchange的真实参数是否一致。
	err = ch.ExchangeDeclare(
		rabbitmq.ExchangeName3,
		"topic", //type
		true,    //durable
		false,   //auto delete
		false,   //internal
		false,   //no-wait
		nil,     //arguments
	)
	if err != nil {
		log.Panicf("declare exchange failed: %s", err)
	}

	Produce("hello zgw", ch, rabbitmq.ExchangeName3, "machine1.debug")
	Produce("hello world", ch, rabbitmq.ExchangeName3, "machine2.debug")
	Produce("hello golang", ch, rabbitmq.ExchangeName3, "machine1.info")
	Produce("hello zzzz", ch, rabbitmq.ExchangeName3, "machine2.info")
}

// go run .\mq\rabbitmq\producer\
