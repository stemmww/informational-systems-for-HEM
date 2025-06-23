package main

import (
	"gep-integration/internal/sender"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Подключаю статику для стилей
	r.Static("/static", "./web")

	// Возврат HTML формы из файла
	r.GET("/form", func(c *gin.Context) {
		c.File("./web/form.html")
	})

	// Обработка формы
	r.POST("/submit-form", sender.HandleForm)

	log.Println(" sender-service слушает на :8081")
	r.Run(":8081")
}
