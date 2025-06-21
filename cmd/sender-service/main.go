package main

import (
	"gep-integration/internal/sender"
	"log"
)

func main() {
	err := sender.SendTestMessage()
	if err != nil {
		log.Fatalf("Ошибка при отправке сообщения: %v", err)
	}
}
