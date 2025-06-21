package sender

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/tiaguinho/gosoap"
)

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
	endpoint := "http://mock-soap:9999/soap"

	// HTTP-клиент с логированием
	transport := &http.Transport{}
	clientWithLogging := &http.Client{
		Transport: loggingRoundTripper{transport},
		Timeout:   30 * time.Second,
	}

	client, err := gosoap.SoapClient(endpoint, clientWithLogging)
	if err != nil {
		return err
	}

	// Чтение и кодирование pdf
	pdfPath := "test.pdf"
	pdfBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		log.Fatalf("Ошибка чтения PDF: %v", err)
	}
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)

	// Подсчёт страниц
	pageCount, err := api.PageCount(bytes.NewReader(pdfBytes), nil)
	if err != nil {
		log.Fatalf("Ошибка определения количества страниц PDF: %v", err)
	}

	// SOAP-параметры
	params := gosoap.Params{
		"RPOInfo": map[string]interface{}{
			"PackageCode":    "UUID-123",
			"SenderID":       123,
			"SenderPass":     "password",
			"DocumentID":     "DOC-456",
			"F1":             "Иванов Иван",
			"F2":             "ул. Пушкина, д. 10",
			"F3":             "Алматы",
			"F4":             "Алматинская",
			"F5":             "050000",
			"F6":             "ООО Ромашка",
			"F7":             "ул. Ленина, 5",
			"F25":            "R200",
			"PageCount":      pageCount,
			"FileAttachment": pdfBase64,
		},
	}

	// XML-просмотр
	fmt.Printf("📦 Параметры: %+v\n", params)
	xmlPreview, _ := xml.MarshalIndent(params, "", "  ")
	fmt.Println("🧾 XML Body:")
	fmt.Println(string(xmlPreview))

	// Отправка
	res, err := client.Call("SendRPO", params)
	if err != nil {
		return err
	}

	log.Printf("Ответ от сервиса: %+v", res)
	return nil
}

// Логирующий RoundTripper
type loggingRoundTripper struct {
	rt http.RoundTripper
}

func (l loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var buf bytes.Buffer
	req.Body = ioutil.NopCloser(io.TeeReader(req.Body, &buf))

	fmt.Println("📤 RAW SOAP-запрос:")
	fmt.Println(buf.String())

	resp, err := l.rt.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("📥 RAW SOAP-ответ:")
	fmt.Println(string(respBody))

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))
	return resp, nil
}
