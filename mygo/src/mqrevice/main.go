package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

const (
	queueName = "push.msg.q"
	exchange  = "t.msg.ex"
	mqurl     = "amqp://rabbituser:dtsgx0591@47.95.145.191:8802"
)

func main() {
	conn, err := amqp.Dial(mqurl)
	failOnErr(err, "failed to connect tp rabbitmq")

	channel, err := conn.Channel()
	failOnErr(err, "failed to open a channel")

	msgs, err := channel.Consume(queueName, "", true, false, false, false, nil)
	failOnErr(err, "")

	forever := make(chan bool)

	go func() {
		//fmt.Println(*msgs)
		for d := range msgs {
			fmt.Printf("receve msg is :%s -- \n", string(d.Body))
		}
	}()

	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever

}

func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}
