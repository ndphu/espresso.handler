package handler

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
)

type Handler interface {
	HandleMessage(message string)
}

func BindHandler(client mqtt.Client, topic string, hander Handler) {
	channel := make(chan string)
	token := client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		channel <- fmt.Sprintf("%s", msg.Payload())
	})
	if token.Wait() && token.Error() != nil {
		log.Panic("Failed to subcribe to topic", topic)
	}

	go func() {
		for {
			msg, success := <-channel
			if !success {
				break
			}
			hander.HandleMessage(msg)
		}
	}()
}
