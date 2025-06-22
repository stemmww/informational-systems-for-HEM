package main

import (
	"gep-integration/internal/receiver"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/status", receiver.ReceiveStatus)

	f, err := os.OpenFile("status.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть лог-файл: %v", err)
	}
	log.SetOutput(f)

	router.Run(":8082")
}
