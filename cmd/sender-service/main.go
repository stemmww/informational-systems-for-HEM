package main

import (
	"gep-integration/internal/sender"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Папка со статическими файлами (например, CSS)
	r.Static("/static", "./web")

	// Отображение HTML формы
	r.GET("/form", sender.ShowForm)

	// Обработка формы и отправка SOAP
	r.POST("/submit-form", sender.HandleForm)

	log.Println("🚀 sender-service слушает на :8081")
	r.Run(":8081")
}
