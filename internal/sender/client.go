package sender

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/tiaguinho/gosoap"
)

// Структура строго по WSDL
type RPOInfo struct {
	XMLName        xml.Name `xml:"RPOInfo"`
	PackageCode    string   `xml:"PackageCode"`
	SenderID       int      `xml:"SenderID"`
	SenderPass     string   `xml:"SenderPass"`
	DocumentID     string   `xml:"DocumentID"`
	F1             string   `xml:"F1"`
	F2             string   `xml:"F2"`
	F3             string   `xml:"F3"`
	F4             string   `xml:"F4"`
	F5             string   `xml:"F5"`
	F6             string   `xml:"F6"`
	F7             string   `xml:"F7,omitempty"`
	F25            string   `xml:"F25,omitempty"`
	PageCount      int      `xml:"PageCount"`
	FileAttachment string   `xml:"FileAttachment"`
}

func SendTestMessage() error {
	endpoint := "http://localhost:8082/status" // временный SOAP

	// Создает SOAP-клиент
	client, err := gosoap.SoapClient(endpoint, http.DefaultClient)
	if err != nil {
		return fmt.Errorf("ошибка создания SOAP клиента: %v", err)
	}

	// Читает PDF
	pdfBytes, err := os.ReadFile("test.pdf")
	if err != nil {
		return fmt.Errorf("ошибка чтения PDF: %v", err)
	}
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)

	pageCount, err := api.PageCount(bytes.NewReader(pdfBytes), nil)
	if err != nil {
		return fmt.Errorf("ошибка определения количества страниц PDF: %v", err)
	}

	// Формирует структуру RPOInfo
	data := RPOInfo{
		PackageCode:    "UUID-123",
		SenderID:       123,
		SenderPass:     "password",
		DocumentID:     "DOC-456",
		F1:             "Иванов Иван",
		F2:             "ул. Пушкина, д. 10",
		F3:             "Алматы",
		F4:             "Алматинская",
		F5:             "050000",
		F6:             "ООО Ромашка",
		F7:             "ул. Ленина, 5",
		F25:            "R200",
		PageCount:      pageCount,
		FileAttachment: pdfBase64,
	}

	// Показывает XML
	xmlPreview, _ := xml.MarshalIndent(data, "", "  ")
	fmt.Println(" XML Body:")
	fmt.Println(string(xmlPreview))

	// Оборачивает структуру в Params
	params := gosoap.Params{
		"RPOInfo": data,
	}

	// Отправка SOAP
	res, err := client.Call("SendRPO", params)
	if err != nil {
		return fmt.Errorf("ошибка при отправке SOAP: %v", err)
	}

	log.Printf(" Ответ от сервиса: %+v", res)
	return nil
}
