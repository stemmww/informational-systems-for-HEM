package main

import (
	"io"
	"log"
	"os"

	"gep-integration/internal/sender"

	"github.com/gin-gonic/gin"
)

func main() {
	// Открытие лог-файла
	f, err := os.OpenFile("sender.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка открытия лог-файла: %v", err)
	}

	// Лог — и в файл, и в терминал
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	r := gin.Default()

	// Статические файлы
	r.Static("/static", "./web")

	// Форма отправки
	r.GET("/form", sender.ShowForm)

	// Иконка браузера (чтобы убрать 404 на /favicon.ico)
	r.StaticFile("/favicon.ico", "./web/favicon.ico")

	// Обработка формы
	r.POST("/submit-form", sender.HandleForm)

	log.Println("sender-service слушает на :8081")
	r.Run(":8081")
}
