# Inventory Management System -Go

A RESTful API for managing inventory with automatic stock updates and email notifications.

![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat\&logo=go)
![Postgres](https://img.shields.io/badge/PostgreSQL-13+-316192?style=flat\&logo=postgresql)
![SMTP](https://img.shields.io/badge/Email-SMTP-EA4335?style=flat\&logo=gmail)

---

## Features

* CRUD for Categories, Products, and Orders
* Automatic stock deduction on orders
* Email notifications on create operations
* Foreign key protection for data integrity
* Asynchronous processing using goroutines

---

## Architecture Overview

```
Client (Postman / Frontend)
        |
        v
   Gin Router (HTTP)
        |
        v
   Usecase Layer (Business Logic)
        |
   ---------------------------
   |                         |
Repository (GORM)     Email Service (SMTP)
        |                         |
        v                         v
 PostgreSQL                  Mail Server
```

---

## Database Schema

```
CATEGORIES
  └── id (PK)
  └── name

PRODUCTS
  └── id (PK)
  └── product_name
  └── quantity
  └── category_id (FK → categories.id)

ORDERS
  └── id (PK)
  └── product_id (FK → products.id)
  └── quantity
```

---

## Tech Stack

* Go (Gin Framework)
* PostgreSQL
* GORM ORM
* SMTP (Gmail / SendGrid)
* godotenv

---

## Getting Started

### Prerequisites

* Go installed
* PostgreSQL running
* SMTP credentials
* Postman or curl

---

### Setup

```bash
git clone <repo-url>
go mod tidy
```

Create `.env` file:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=inventory_user
DB_PASS=password
DB_NAME=inventory_db

SERVER_PORT=8080

SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SENDER_EMAIL=your@gmail.com
SENDER_PASSWORD=app-password
RECEIVER_EMAIL=receiver@gmail.com
```

Run application:

```bash
go run cmd/main.go
```

---

## Email Notification Flow

```
Create Request
      |
      v
 Save to Database
      |
      +----> Spawn Goroutine ----> Send Email (SMTP)
      |
   Return API Response
```

Emails are sent for:

* Category creation
* Product creation
* Order creation

---

## API Endpoints

Base URL:

```
http://localhost:8080/api/v1
```

### Categories

* POST `/categories`
* GET `/categories`
* PUT `/categories/{id}`
* DELETE `/categories/{id}`

### Products

* POST `/products`
* GET `/products`
* PUT `/products/{id}`
* DELETE `/products/{id}`

### Orders

* POST `/orders`
* GET `/orders`

### Aggregated

* GET `/get-tasks`
* POST `/get-by-tablename`

---

## Key Logic

### Stock Management

```
Order Request
     |
Check Available Stock
     |
Enough? ---- No ----> Reject
     |
    Yes
     |
Reduce Product Quantity
     |
Create Order
```

---

## Project Structure

```
inventory-api/
├── cmd/main.go
├── database/database.go
├── service/email.go
├── internal/
│   ├── domain/entities.go
│   ├── repository/repository.go
│   ├── usecase/usecase.go
│   └── delivery/http/
│       ├── handler.go
│       └── routes.go
├── .env
├── go.mod
└── README.md
```

---

## Testing

* Create category → email sent
* Create product → email sent
* Create order → stock reduced + email sent
* Order more than stock → error
* Delete category with products → error

---

## Troubleshooting

* Email not received → check spam and SMTP credentials
* Auth failed → use Gmail App Password
* FK error → delete child records first
* DB error → verify `.env`

---

## License

MIT


