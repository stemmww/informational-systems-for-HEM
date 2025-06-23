# HEM-INTEGRATION

Информационная система для отправки заказных писем в формате PDF через систему Гибридной Электронной Почты (ГЭП) с возможностью получения статусов доставки.

## Цель проекта

Проект реализует прототип прямой интеграции с системой ГЭП АО "Казпочта". Система позволяет пользователю:

1. Загрузить PDF-документ, содержащий заказное письмо.
2. Заполнить сведения о получателе и отправителе (ФИО, адреса, индекс и т.д.).
3. Отправить сформированный XML-документ (структура `RPOInfo`) в систему ГЭП через SOAP.
4. Получить статус доставки от ГЭП (через входящий SOAP-запрос).

## Структура проекта
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
│   │   ├── client.go
│   │   └── form.go
│   └── receiver/
│       └── status.go
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

## Установка

1. Установите [Docker](https://www.docker.com/) и [Docker Compose](https://docs.docker.com/compose/).
2. Клонируйте репозиторий:

```bash
git clone https://github.com/yourname/gep-integration.git
cd gep-integration
Постройте и запустите контейнеры:

docker-compose up --build
🔹 Сервисы будут доступны по адресам:

Отправка: http://localhost:8081

Получение статуса: http://localhost:8082
```

## Использование
```
Отправка PDF-документа в ГЭП
Откройте страницу в браузере: http://localhost:8081/form

Заполните поля формы:

F1 – ФИО получателя

F2 – Улица, дом и прочее

F3 – Город

F4 – Область / район

F5 – Почтовый индекс

F6 – Организация отправителя

F7 – Адрес отправителя

F25 – Тип отправки (например, R200)

PDF-файл – заказное письмо в виде PDF (его содержимое будет встроено в SOAP)

Нажмите «Отправить». SOAP-запрос будет сгенерирован и отправлен в ГЭП.

Приём статуса доставки (эмуляция ГЭП)
Вы можете вручную протестировать приём статусов через:

curl -X POST http://localhost:8082/status \
  -H "Content-Type: text/xml" \
  --data "@test-data/soap-request.xml"
```
  
## Технологии

Go 1.22+

Gin (Web Framework)

gosoap

PDFCPU

Docker + Docker Compose

HTML / CSS (форма)

SoapUI / Postman (для ручного тестирования)
