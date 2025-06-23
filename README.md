# HEM-INTEGRATION

Information system for sending registered letters in PDF format via Hybrid Electronic Mail (HEM) system with the ability to receive delivery statuses.

## Project goal

The project implements a prototype of direct integration with the Kazpost JSC HEM system. The system allows the user to:

1. Download the PDF document containing the registered letter.
2. Fill in the information about the recipient and sender (full name, addresses, index, etc.).
3. Send the generated XML-document (`RPOInfo` structure) to the ERT system via SOAP.
4. Receive the delivery status from the ERT (via incoming SOAP request).

## Project structure
```
GEP-INTEGRATION/
├── cmd/
│   ├── sender-service/
│   │   └── main.go
│   └── status-receiver/
│       └── main.go
│
├── internal/
│   ├── sender/
│   │   └── form.go
│   └── receiver/
│       └── status.go
│
├── myservice/
│   └── gep_client.go
│
├── test-data/
│   ├── soap-request.xml
│   └── test.pdf
│
├── web/
│   ├── form.html
│   └── style.css
│
├── wsdl/
│   ├── HybridMail.wsdl
│   └── status-receiver.wsdl
├── docker-compose.yaml
├── Dockerfile.sender
├── Dockerfile.receiver
├── go.mod
└── go.sum 

```

## Installation

1. Install [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/).
2. Clone the repository:

```bash
git clone https://github.com/yourname/gep-integration.git
cd gep-integration
Build and run the containers:

docker-compose up --build
🔹 Services will be available at:

Sending: http://localhost:8081

Receive status: http://localhost:8082
```

## Utilization
```
Sending a PDF document to GEO
Open the page in your browser: http://localhost:8081/form

Fill in the form fields:

F1 - Full name of the recipient

F2 - Street, house, etc.

F3 - City

F4 - Region / district

F5 - Postal code

F6 - Sender's organization

F7 - Sender's address

F25 - Send type (e.g. R200)

PDF file - registered letter as a PDF file (its contents will be embedded in SOAP).

Click “Send”. The SOAP request will be generated and sent to the GEP.
```

## Manual Test: Receiving Status (HEM Emulation)
```
You can test the receiving of delivery statuses either via curl or using SoapUI.

Option 1: cURL

curl -X POST http://localhost:8082/status \
  -H "Content-Type: text/xml" \
  --data "@test-data/soap-request.xml"

Option 2: SoapUI

To test the /status endpoint using SoapUI:

Create a new SOAP project

Open SoapUI

Go to File → New SOAP Project

Project name: status-receiver

In WSDL path, provide:

file:/C:/path/to/wsdl/status-receiver.wsdl

Click OK

Open the SendStatus request

Expand status-receiver → StatusReceiverBinding → SendStatus

You will see a request template

Paste sample request body:

<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
                  xmlns:stat="http://kazpost.kz/StatusService">
  <soapenv:Header/>
  <soapenv:Body>
    <stat:KazpostStatus>
      <stat:id>123</stat:id>
      <stat:barcode>ABC-456</stat:barcode>
      <stat:date>2025-06-23</stat:date>
      <stat:status>Delivered</stat:status>
      <stat:operator>Ivanov</stat:operator>
      <stat:recipient>Smith P.</stat:recipient>
      <stat:returnReason></stat:returnReason>
    </stat:KazpostStatus>
  </soapenv:Body>
</soapenv:Envelope>

Set the endpoint:

http://localhost:8082/status

Click “Send”You should receive this response:
```

  
## Technologies

Go 1.22+

Gin (Web Framework)

gowsdl (WSDL → Go client generator)

PDFCPU (PDF page count)

Docker + Docker Compose

HTML / CSS (Web UI for sending)

SoapUI / Postman (Manual testing)
