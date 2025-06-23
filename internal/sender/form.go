package sender

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/tiaguinho/gosoap"
)

// Обработчик формы
func HandleForm(c *gin.Context) {
	// Чтение полей
	f1 := c.PostForm("F1")
	f2 := c.PostForm("F2")
	f3 := c.PostForm("F3")
	f4 := c.PostForm("F4")
	f5 := c.PostForm("F5")
	f6 := c.PostForm("F6")
	f7 := c.PostForm("F7")
	f25 := c.PostForm("F25")

	// Получает PDF
	file, _, err := c.Request.FormFile("pdf")
	if err != nil {
		c.String(http.StatusBadRequest, "Ошибка чтения PDF: %v", err)
		return
	}
	defer file.Close()

	// Чтение файла в память
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		c.String(http.StatusInternalServerError, "Ошибка буфера: %v", err)
		return
	}
	pdfBytes := buf.Bytes()

	pageCount, err := api.PageCount(bytes.NewReader(pdfBytes), nil)
	if err != nil {
		c.String(http.StatusBadRequest, "Ошибка определения страниц: %v", err)
		return
	}

	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)

	// SOAP отправка
	data := RPOInfo{
		PackageCode:    "UUID-123",
		SenderID:       123,
		SenderPass:     "password",
		DocumentID:     "DOC-456",
		F1:             f1,
		F2:             f2,
		F3:             f3,
		F4:             f4,
		F5:             f5,
		F6:             f6,
		F7:             f7,
		F25:            f25,
		PageCount:      pageCount,
		FileAttachment: pdfBase64,
	}

	xmlPreview, _ := xml.MarshalIndent(data, "", "  ")
	fmt.Println(" XML Body:")
	fmt.Println(string(xmlPreview))

	endpoint := "http://localhost:8081/soap"
	client, err := gosoap.SoapClient(endpoint, http.DefaultClient)
	if err != nil {
		c.String(http.StatusInternalServerError, "SOAP client error: %v", err)
		return
	}

	params := gosoap.Params{"RPOInfo": data}
	res, err := client.Call("SendRPO", params)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка отправки SOAP-запроса: %v", err)
		return
	}

	log.Printf(" Ответ от сервиса: %+v", res)
	c.String(http.StatusOK, "Сообщение успешно отправлено")
}
