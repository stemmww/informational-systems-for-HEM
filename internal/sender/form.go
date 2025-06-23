package sender

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"gep-integration/myservice"

	"github.com/gin-gonic/gin"
	"github.com/hooklift/gowsdl/soap"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func ShowForm(c *gin.Context) {
	c.File("web/form.html")
}

func HandleForm(c *gin.Context) {
	f1 := c.PostForm("F1")
	f2 := c.PostForm("F2")
	f3 := c.PostForm("F3")
	f4 := c.PostForm("F4")
	f5 := c.PostForm("F5")
	f6 := c.PostForm("F6")
	f7 := c.PostForm("F7")
	f25 := c.PostForm("F25")

	file, err := c.FormFile("pdf")
	if err != nil {
		c.String(http.StatusBadRequest, "Ошибка загрузки файла: %v", err)
		return
	}

	src, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка открытия файла: %v", err)
		return
	}
	defer src.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(src)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка чтения файла: %v", err)
		return
	}

	pdfBytes := buf.Bytes()
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)

	tmpPath := os.TempDir() + "/temp.pdf"
	if err := os.WriteFile(tmpPath, pdfBytes, 0644); err != nil {
		c.String(http.StatusInternalServerError, "Ошибка временной записи файла: %v", err)
		return
	}

	pageCount, err := api.PageCountFile(tmpPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка определения количества страниц: %v", err)
		return
	}

	client := soap.NewClient("http://localhost:8082/soap")
	hybrid := myservice.NewHybridMail(client)

	rpo := &myservice.RPOInfo{
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
		PageCount:      int32(pageCount),
		FileAttachment: pdfBase64,
	}

	req := &myservice.SendRPORequest{
		RPOInfo: rpo,
	}

	resp, err := hybrid.SendRPO(req)
	if err != nil {
		log.Printf("Ошибка при отправке SOAP-запроса: %v", err)
		c.String(http.StatusInternalServerError, "Ошибка отправки SOAP-запроса: %v", err)
		return
	}

	info := resp.ResponseInfo
	if info == nil {
		log.Println("Пустой ответ от сервиса")
		c.String(http.StatusInternalServerError, "Ошибка: пустой ответ от сервиса")
		return
	}

	log.Printf("Ответ от сервиса:\n- Msg: %s\n- KPST ID: %s\n- Time: %s",
		info.ResponseMsg, info.ResponseKpstID, info.ResponseTime)

	log.Printf("Сохранено в лог: KPST_ID=%s, Документ=%s, Получатель=%s",
		info.ResponseKpstID, rpo.DocumentID, rpo.F1)

	if info.ResponseMsg != "Документ принят" {
		c.String(http.StatusOK, "Ответ: %s\nID: %s", info.ResponseMsg, info.ResponseKpstID)
	} else {
		c.String(http.StatusOK, "Успешно отправлено!\nОтвет: %s\nID: %s\nВремя: %s",
			info.ResponseMsg, info.ResponseKpstID, info.ResponseTime)
	}

}
