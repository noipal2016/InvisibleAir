package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)
const (
	queueName = "push.msg.q"
	exchange  = "t.msg.ex"
	mqurl ="amqp://rabbituser:dtsgx0591@47.95.145.191:8802"
)

func main() {
	conn, err := amqp.Dial(mqurl)
	failOnErr(err, "failed to connect tp rabbitmq")
	defer conn.Close()
	channel, err := conn.Channel()
	failOnErr(err, "failed to open a channel")
	defer channel.Close()



	err = channel.ExchangeDeclare(
		exchange,
		amqp.ExchangeHeaders,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnErr(err, "Failed to declare a Exchange")
	_, err = channel.QueueDeclare(
		queueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnErr(err, "Failed to declare a queue")

	err = channel.QueueBind(queueName,queueName,exchange,true,nil)

	failOnErr(err, "Failed to bind a queue")

	go func() {
		count:=0
		msgContent := "hello world!"
		for  {
			channel.Publish(exchange, queueName, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msgContent),
			})
			count++
			log.Println(count)
			time.Sleep(3*time.Second)
		}

	}()

	time.Sleep(10000 * time.Hour)
}

func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}
