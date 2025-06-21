package receiver

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Структура входящего SOAP-обёртки
type Envelope struct {
	XMLName xml.Name     `xml:"Envelope"`
	Body    EnvelopeBody `xml:"Body"`
}

type EnvelopeBody struct {
	Status KazpostStatus `xml:"SendStatus"`
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

	if err := c.ShouldBindXML(&env); err != nil {
		fmt.Println("❌ Ошибка разбора XML:", err)
		c.String(http.StatusBadRequest, "Ошибка XML")
		return
	}

	fmt.Println("📥 Получен статус от ГЭП:")
	fmt.Printf("%+v\n", env.Body.Status)

	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, `<ACCEPT>1</ACCEPT>`)
}
