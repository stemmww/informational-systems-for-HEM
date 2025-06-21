package main

import (
	"gep-integration/internal/receiver"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/status", receiver.ReceiveStatus)
	router.Run(":8082")
}
