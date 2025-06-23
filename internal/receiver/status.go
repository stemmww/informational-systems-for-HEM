package receiver

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Структура входящего SOAP-обёртки
type Envelope struct {
	XMLName xml.Name     `xml:"Envelope"`
	Body    EnvelopeBody `xml:"Body"`
}

type EnvelopeBody struct {
	Status KazpostStatus `xml:"KazpostStatus"`
}

// Структура статуса от ГЭП
type KazpostStatus struct {
	XMLName      xml.Name `xml:"KazpostStatus"`
	ID           string   `xml:"id"`
	Barcode      string   `xml:"barcode"`
	Date         string   `xml:"date"`
	Status       string   `xml:"status"`
	Operator     string   `xml:"operator"`
	Recipient    string   `xml:"recipient"`
	ReturnReason string   `xml:"returnReason"`
}

// Обработчик запроса от ГЭП
func ReceiveStatus(c *gin.Context) {
	var env Envelope

	// Парсит входящий SOAP XML
	if err := c.ShouldBindXML(&env); err != nil {
		fmt.Println(" Ошибка разбора XML:", err)
		c.String(http.StatusBadRequest, "Ошибка XML")
		return
	}

	status := env.Body.Status

	// Валидация обязательных полей
	if status.ID == "" || status.Barcode == "" || status.Date == "" || status.Status == "" {
		fmt.Println(" Ошибка валидации: обязательные поля пусты")
		c.String(http.StatusBadRequest, "Ошибка: обязательные поля пусты (id, barcode, date, status)")
		return
	}

	// Успех
	fmt.Println(" Получен статус от ГЭП:")
	fmt.Printf("%+v\n", status)

	c.Header("Content-Type", "application/xml")

	log.Printf(" Получен статус от ГЭП: ID=%s, Barcode=%s, Date=%s, Status=%s, Recipient=%s\n",
		status.ID, status.Barcode, status.Date, status.Status, status.Recipient)

	c.String(http.StatusOK, `<ACCEPT>1</ACCEPT>`)
}
