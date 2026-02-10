# Inventory Management System API

A RESTful API service for managing inventory operations with automatic stock management, real-time notifications, and email alerts.

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-316192?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![SMTP](https://img.shields.io/badge/Email-SMTP-EA4335?style=flat&logo=gmail)](https://support.google.com/mail/)

---

## Table of Contents

- [What This API Does](#what-this-api-does)
- [System Architecture](#system-architecture)
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
- [Email Notification System](#email-notification-system)
- [API Usage](#api-usage)
- [Project Structure](#project-structure)
- [Key Features Explained](#key-features-explained)
- [Testing Guide](#testing-guide)
- [Troubleshooting](#troubleshooting)

---

## What This API Does

This inventory management system helps you:

- **Organize Products** into categories (Electronics, Clothing, etc.)
- **Track Stock Levels** in real-time
- **Process Orders** with automatic inventory updates
- **Send Email Notifications** for all major operations
- **Maintain Data Integrity** through foreign key relationships
- **Handle Concurrent Operations** using Go goroutines

### System Flow Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLIENT REQUEST                            │
│                    (Postman, cURL, Frontend)                     │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      GIN WEB FRAMEWORK                           │
│                     (HTTP Router/Handler)                        │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      BUSINESS LOGIC LAYER                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  Categories  │  │   Products   │  │    Orders    │          │
│  │   UseCase    │  │   UseCase    │  │   UseCase    │          │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘          │
│         │                  │                  │                  │
│         │    Validation    │    Stock Check   │                 │
│         │    Goroutines    │    Email Send    │                 │
└─────────┼──────────────────┼──────────────────┼─────────────────┘
          │                  │                  │
          ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                    REPOSITORY LAYER (GORM)                       │
│                     Database Abstraction                         │
└────────────────────────────┬────────────────────────────────────┘
                             │
          ┌──────────────────┴──────────────────┐
          │                                     │
          ▼                                     ▼
┌──────────────────────┐            ┌──────────────────────┐
│  POSTGRESQL DATABASE │            │   EMAIL SERVICE      │
│  ┌────────────┐      │            │   (SMTP Gateway)     │
│  │ Categories │      │            │                      │
│  │  Products  │      │            │  Gmail/SendGrid/etc  │
│  │   Orders   │      │            │                      │
│  └────────────┘      │            └──────────────────────┘
└──────────────────────┘
```

---

## System Architecture

### Database Schema

```
┌─────────────────────┐
│    CATEGORIES       │
├─────────────────────┤
│ id (UUID) PK        │◄─────────┐
│ name (VARCHAR)      │          │
│ created_at          │          │
│ updated_at          │          │
└─────────────────────┘          │
                                 │ Foreign Key
                                 │ ON DELETE RESTRICT
┌─────────────────────┐          │
│     PRODUCTS        │          │
├─────────────────────┤          │
│ id (UUID) PK        │◄─────────┼─────────┐
│ product_name        │          │         │
│ description         │          │         │
│ price (DECIMAL)     │          │         │
│ quantity (INT)      │          │         │ Foreign Key
│ category_id FK ─────┼──────────┘         │ ON DELETE RESTRICT
│ is_active (BOOL)    │                    │
│ created_at          │                    │
│ updated_at          │                    │
└─────────────────────┘                    │
                                           │
┌─────────────────────┐                    │
│      ORDERS         │                    │
├─────────────────────┤                    │
│ id (UUID) PK        │                    │
│ product_id FK ──────┼────────────────────┘
│ quantity (INT)      │
│ order_date          │
│ created_at          │
└─────────────────────┘
```

### Complete Request Flow with Email Notifications

```
POST /api/v1/orders
│
├─ Step 1: Validate Request
│  └─ Check JSON format
│     Check required fields
│
├─ Step 2: Find Product
│  └─ Query database for ProductID
│     Fetch product name for email
│     Return 404 if not found
│
├─ Step 3: Check Stock
│  └─ Compare order quantity with product quantity
│     Return 400 if insufficient
│
├─ Step 4: Update Stock (Transaction)
│  └─ Decrement product.quantity
│     Create order record
│     Commit both or rollback
│
├─ Step 5: Send Email Notification (Goroutine)
│  └─ Spawn async process ───────────────────┐
│     Continue to Step 6                     │
│                                            │
├─ Step 6: Return Response                   │
│  └─ 201 Created with order details         │
│     (User receives instant response)       │
│                                            │
│  ┌─────────────────────────────────────────┘
│  │
│  └─► Background Email Process:
│      ├─ Connect to SMTP server
│      ├─ Format email with product name & quantity
│      ├─ Send to configured recipient
│      ├─ Log success/failure
│      └─ Complete (non-blocking)
```

---

## Technology Stack

```
┌─────────────────────────────────────────────────┐
│              Application Layer                   │
│  ┌───────────────────────────────────────┐     │
│  │  Go 1.19+                              │     │
│  │  - Gin Web Framework (HTTP)            │     │
│  │  - net/smtp (Email)                    │     │
│  │  - Goroutines (Async)                  │     │
│  │  - godotenv (Config)                   │     │
│  └───────────────────────────────────────┘     │
└─────────────────────────────────────────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────┐
│               Data Access Layer                  │
│  ┌───────────────────────────────────────┐     │
│  │  GORM (ORM)                            │     │
│  │  - Query Builder                       │     │
│  │  - Relationship Preloading             │     │
│  │  - Transaction Management              │     │
│  └───────────────────────────────────────┘     │
└─────────────────────────────────────────────────┘
                      │
          ┌───────────┴───────────┐
          ▼                       ▼
┌──────────────────┐    ┌──────────────────┐
│  Database Layer  │    │  Email Service   │
│  ┌────────────┐  │    │  ┌────────────┐  │
│  │PostgreSQL  │  │    │  │ SMTP Server│  │
│  │  13+       │  │    │  │ (Gmail/etc)│  │
│  │- UUID Ext  │  │    │  │- Port 587  │  │
│  │- FK Const  │  │    │  │- TLS/SSL   │  │
│  └────────────┘  │    │  └────────────┘  │
└──────────────────┘    └──────────────────┘
```

---

## Getting Started

### Prerequisites Checklist

- [ ] Go 1.19 or higher installed
- [ ] PostgreSQL 13 or higher installed
- [ ] pgAdmin 4 (for database management)
- [ ] Gmail account or SMTP server credentials
- [ ] Postman or cURL (for testing)

### Step 1: Clone and Install Dependencies

```bash
# Clone the repository
git clone <repository-url>
cd inventory-api

# Initialize Go module
go mod init inventory-api

# Install dependencies
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/google/uuid
go get -u github.com/joho/godotenv
go get -u github.com/gin-gonic/gin

# Download all dependencies
go mod download
go mod tidy
```

### Step 2: Database Setup

#### Quick Setup SQL Script

```sql
-- Create database
CREATE DATABASE inventory_db;

-- Connect to the database
\c inventory_db

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create products table
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_name VARCHAR(200) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    category_id UUID NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
);

-- Create orders table
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
);

-- Create indexes for better performance
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_active ON products(is_active);
CREATE INDEX idx_orders_product ON orders(product_id);
CREATE INDEX idx_orders_date ON orders(order_date);

-- Create application user
CREATE USER inventory_user WITH PASSWORD 'your_secure_password';

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE inventory_db TO inventory_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO inventory_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO inventory_user;
```

### Step 3: Email Configuration Setup

#### For Gmail Users (Recommended for Testing)

1. **Enable 2-Factor Authentication** on your Google account
2. **Generate App Password:**
   - Go to Google Account Settings
   - Security → 2-Step Verification → App passwords
   - Generate password for "Mail" on "Other (Custom name)"
   - Copy the 16-character password

#### For Other SMTP Providers

| Provider | SMTP Host | Port | Notes |
|----------|-----------|------|-------|
| Gmail | smtp.gmail.com | 587 | Requires App Password |
| Outlook | smtp-mail.outlook.com | 587 | Use account password |
| SendGrid | smtp.sendgrid.net | 587 | Use API key as password |
| AWS SES | email-smtp.region.amazonaws.com | 587 | Use SMTP credentials |

### Step 4: Environment Configuration

Create `.env` file in project root:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=inventory_user
DB_PASS=your_secure_password
DB_NAME=inventory_db

# Server Configuration
SERVER_PORT=8080

# Email Configuration (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SENDER_EMAIL=your-email@gmail.com
SENDER_PASSWORD=your-app-password-here
RECEIVER_EMAIL=recipient-email@gmail.com
```

**Security Notes:**
- Never commit `.env` file to version control
- Add `.env` to your `.gitignore` file
- Use environment variables in production
- Rotate SMTP credentials regularly

### Step 5: Run the Application

```bash
# Run directly
go run cmd/main.go

# Expected output:
# Database connected successfully
# Email service initialized
# SMTP configured: smtp.gmail.com:587
# Server running on port 8080...
```

---

## Email Notification System

### What It Does

The email notification system sends automatic alerts when key actions are performed:

```
┌─────────────────────────────────────────────────────────┐
│               EMAIL NOTIFICATION TRIGGERS                │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ✉ Category Created                                     │
│    → "New category created: Electronics"                │
│                                                          │
│  ✉ Product Created                                      │
│    → "New product added: USB Cable (Electronics)"       │
│    → "Quantity: 100, Price: $15.99"                     │
│                                                          │
│  ✉ Order Placed                                         │
│    → "Order created for: USB Cable"                     │
│    → "Quantity ordered: 5"                              │
│    → "Remaining stock: 95"                              │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

###  How It Works

```
API Request → Handler → Business Logic
                           │
                           ├─► Database Operation (synchronous)
                           │   └─ Save data immediately
                           │
                           └─► Email Notification (asynchronous)
                               └─ Goroutine spawned
                                  ├─ Fetch additional data (product name)
                                  ├─ Format email message
                                  ├─ Connect to SMTP server
                                  ├─ Send email
                                  └─ Log result
                                     (runs in background)
```

### Email Flow Diagram

```
POST /api/v1/products
│
├─ Main Thread (< 50ms):
│  ├─ Validate request body
│  ├─ Insert product into database
│  ├─ Spawn email goroutine ────────────┐
│  └─ Return 201 response to client     │
│                                        │
│  ┌─────────────────────────────────────┘
│  │
│  └─► Background Email Goroutine (500-2000ms):
│      ├─ Connect to SMTP server (smtp.gmail.com:587)
│      ├─ Authenticate with credentials
│      ├─ Format email:
│      │  Subject: New Product Created
│      │  Body: Product: USB Cable
│      │        Category: Electronics
│      │        Price: $15.99
│      │        Quantity: 100
│      ├─ Send email via TLS
│      ├─ Log: "Email sent successfully"
│      └─ Goroutine exits
```

###  Email Service Architecture

```
service/
├── email.go
│   ├── EmailConfig struct
│   │   ├── SMTPHost
│   │   ├── SMTPPort
│   │   ├── SenderEmail
│   │   ├── SenderPassword
│   │   └── ReceiverEmail
│   │
│   ├── NewEmailService()
│   │   └── Loads config from environment
│   │
│   └── SendEmail(subject, body)
│       ├── Builds MIME message
│       ├── Authenticates with SMTP
│       └── Sends via smtp.SendMail()
│
└── Integration Points:
    ├── usecase/usecase.go (business logic)
    └── delivery/http/handler.go (HTTP handlers)
```

### Email Templates

#### Category Created Email
```
Subject: Category Created

Body:
New Category Created: Electronics

```

#### Product Created Email
```
Subject: Product Created

Body:
New Product Created: iPhone
```

#### Order Created Email
```
Subject: Order Created

Body:
Order Created
Product: iPhone
Quantity: 1
```

###  Testing Email Notifications

```bash
# Test 1: Category Creation
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{"Name":"Electronics"}'

# Expected:
# - API returns 201 immediately
# - Email arrives within 5-10 seconds
# - Check inbox for "New Category Created"

# Test 2: Product Creation
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "ProductName":"USB Cable",
    "Description":"High-speed cable",
    "Price":15.99,
    "Quantity":100,
    "CategoryID":"<category-uuid-here>"
  }'

# Expected:
# - API returns 201 immediately
# - Email shows product details
# - Check for "New Product Added"

# Test 3: Order Creation
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "ProductID":"<product-uuid-here>",
    "Quantity":5
  }'

# Expected:
# - API returns 201 immediately
# - Email shows product name (not UUID)
# - Email shows quantity ordered
# - Check for "New Order Placed"
```

### Design Benefits

```
┌────────────────────────────────────────────────────────┐
│         WHY ASYNCHRONOUS EMAIL SENDING?                 │
├────────────────────────────────────────────────────────┤
│                                                         │
│  Performance:                                           │
│  ├─ API response: ~10-50ms (without email)            │
│  ├─ Email sending: ~500-2000ms (SMTP latency)         │
│  └─ With goroutines: User doesn't wait               │
│                                                         │
│  Reliability:                                           │
│  ├─ Email failure doesn't break API response          │
│  ├─ Database commit happens regardless                │
│  └─ Retry logic possible in background                │
│                                                         │
│  Scalability:                                           │
│  ├─ Server handles 1000s of concurrent requests       │
│  ├─ Email queue processed independently                │
│  └─ No request blocking                                │
│                                                         │
│  Separation of Concerns:                                │
│  ├─ Business logic isolated from infrastructure       │
│  ├─ Email service can be swapped (SMTP → SES → etc)  │
│  └─ Testable without email server                     │
│                                                         │
└────────────────────────────────────────────────────────┘
```

---

## API Usage

### Base URL
```
http://localhost:8080/api/v1
```

### Complete Endpoint Reference

```
┌─────────────────────────────────────────────────────────────┐
│                      API ENDPOINTS                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  CATEGORIES                                                  │
│  ├─ GET    /categories          List all + Email on CREATE  │
│  ├─ POST   /categories          Create new ✉                │
│  ├─ PUT    /categories/{id}     Update                      │
│  └─ DELETE /categories/{id}     Delete                      │
│                                                              │
│  PRODUCTS                                                    │
│  ├─ GET    /products            List all + Email on CREATE  │
│  ├─ POST   /products            Create new ✉                │
│  ├─ PUT    /products/{id}       Update                      │
│  └─ DELETE /products/{id}       Delete                      │
│                                                              │
│  ORDERS                                                      │
│  ├─ GET    /orders              List all + Email on CREATE  │
│  └─ POST   /orders              Create ✉ (auto-deduct)      │
│                                                              │
│  AGGREGATED                                                  │
│  ├─ GET    /get-tasks           All data at once            │
│  └─ POST   /get-by-tablename    Get by table & ID           │
│                                                              │
│  ✉ = Email notification sent asynchronously                 │
└─────────────────────────────────────────────────────────────┘
```


### cURL Command Reference

```bash
# ============================================
# CATEGORY OPERATIONS
# ============================================

# Create category (triggers email)
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{"Name":"Electronics"}'

# Get all categories
curl http://localhost:8080/api/v1/categories

# Update category
curl -X PUT http://localhost:8080/api/v1/categories/{id} \
  -H "Content-Type: application/json" \
  -d '{"Name":"Updated Electronics"}'

# Delete category
curl -X DELETE http://localhost:8080/api/v1/categories/{id}


# ============================================
# PRODUCT OPERATIONS
# ============================================

# Create product (triggers email)
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "ProductName":"USB Cable",
    "Description":"High-speed USB Type-C cable",
    "Price":15.99,
    "Quantity":100,
    "CategoryID":"550e8400-e29b-41d4-a716-446655440000",
    "IsActive":true
  }'

# Get all products with categories
curl http://localhost:8080/api/v1/products

# Update product
curl -X PUT http://localhost:8080/api/v1/products/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "ProductName":"Premium USB Cable",
    "Price":19.99,
    "Quantity":150
  }'

# Delete product
curl -X DELETE http://localhost:8080/api/v1/products/{id}


# ============================================
# ORDER OPERATIONS
# ============================================

# Create order (triggers email with product name)
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "ProductID":"660e8400-e29b-41d4-a716-446655440111",
    "Quantity":5
  }'

# Get all orders
curl http://localhost:8080/api/v1/orders


# ============================================
# AGGREGATED DATA
# ============================================

# Get everything (categories, products, orders)
curl http://localhost:8080/api/v1/get-tasks

# Get specific record by table and ID
curl -X POST http://localhost:8080/api/v1/get-by-tablename \
  -H "Content-Type: application/json" \
  -d '{
    "TableName":"products",
    "ID":"660e8400-e29b-41d4-a716-446655440111"
  }'
```

---

## Project Structure

```
inventory-api/
│
├── cmd/
│   └── main.go                    # Application entry point
│                                  # - Database initialization
│                                  # - Email service setup
│                                  # - Router configuration
│
├── database/
│   └── database.go                # Database connection
│                                  # - PostgreSQL connection pool
│                                  # - GORM configuration
│
├── service/
│   └── email.go                   # Email notification service
│                                  # - SMTP configuration
│                                  # - Email sending logic
│                                  # - Template formatting
│
├── internal/                      # Internal packages
│   │
│   ├── domain/
│   │   └── entities.go            # Data models
│   │                              #    - Category
│   │                              #    - Product
│   │                              #    - Order
│   │
│   ├── repository/
│   │   └── repository.go          # Database operations
│   │                              #    - CRUD functions
│   │                              #    - Query logic
│   │                              #    - Preloading
│   │
│   ├── usecase/
│   │   └── usecase.go             # Business logic
│   │                              #    - Validation
│   │                              #    - Stock management
│   │                              #    - Email triggers
│   │                              #    - Goroutines
│   │
│   └── delivery/
│       └── http/
│           ├── handler.go         # HTTP handlers
│           │                      #    - Request parsing
│           │                      #    - Response formatting
│           │                      #    - Email coordination
│           │
│           └── routes.go          # Route definitions
│                                  #    - API endpoints
│                                  #    - Middleware
│
├── .env                           # Environment variables
│                                  # - Database config
│                                  # - SMTP config
│                                  # - Server config
│
├── .gitignore                     # Git ignore rules
│                                  # - .env file
│                                  # - Binary files
│
├── go.mod                         # Go dependencies
├── go.sum                         # Dependency checksums
└── README.md                      # This file
```

### Detailed Architecture Layers

```
┌──────────────────────────────────────────────────────────┐
│                  DELIVERY LAYER                           │
│  (HTTP Handlers, Routes, JSON Serialization)             │
│                                                           │
│  Responsibilities:                                        │
│  • Parse HTTP requests (JSON → Go structs)               │
│  • Call use case functions                               │
│  • Return JSON responses                                 │
│  • Handle HTTP status codes                              │
│  • Coordinate email notifications                        │
└─────────────────┬────────────────────────────────────────┘
                  │
                  ▼
┌──────────────────────────────────────────────────────────┐
│                 USE CASE LAYER                            │
│  (Business Logic, Validation, Orchestration)             │
│                                                           │
│  Responsibilities:                                        │
│  • Validate business rules                               │
│  • Check stock availability                              │
│  • Spawn email goroutines                                │
│  • Coordinate repository operations                      │
│  • Manage transactions                                   │
└─────────────────┬────────────────────────────────────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
        ▼                   ▼
┌────────────────┐   ┌────────────────┐
│  REPOSITORY    │   │  EMAIL SERVICE │
│     LAYER      │   │                │
│                │   │  • SMTP config │
│ • SQL queries  │   │  • Send mail   │
│ • Transactions │   │  • Templates   │
│ • Relationships│   │  • Async send  │
└────────┬───────┘   └────────────────┘
         │
         ▼
┌──────────────────────────────────────┐
│         DATABASE LAYER                │
│  (PostgreSQL, Tables, Constraints)   │
└──────────────────────────────────────┘
```

---

## Key Features Explained

### 1. Automatic Stock Management

**How It Works:**

```
Order Creation Process:
│
├─ User submits order for 5 USB cables
│
├─ System checks: Product quantity = 100
│  └─ 100 >= 5? Yes, proceed
│
├─ Transaction starts:
│  ├─ Update products: quantity = 100 - 5 = 95
│  └─ Insert order: record with quantity = 5
│
├─ Transaction commits (both operations succeed)
│
├─ Email notification spawned asynchronously
│  └─ Fetch product name from database
│     Format email with details
│     Send to configured recipient
│
└─ Response: Order created, stock now 95
```

**Business Rules:**
- Orders cannot exceed available stock
- Stock updates are atomic (all or nothing)
- Real-time inventory tracking
- Email confirmation for every order

### 2. Email Notification Architecture

**Separation of Concerns:**

```
┌────────────────────────────────────────────────┐
│          APPLICATION ARCHITECTURE               │
├────────────────────────────────────────────────┤
│                                                 │
│  HTTP Handler Layer                            │
│  ├─ Receives request                           │
│  ├─ Validates input                            │
│  └─ Calls UseCase                              │
│                                                 │
│  ↓                                              │
│                                                 │
│  UseCase Layer (Business Logic)                │
│  ├─ Performs database operation                │
│  ├─ Triggers email service                     │
│  └─ Returns result                             │
│                                                 │
│  ↓                                              │
│                                                 │
│  Email Service (Infrastructure)                │
│  ├─ Independent module                         │
│  ├─ Loaded from environment                    │
│  ├─ SMTP connection handling                   │
│  └─ Can be mocked for testing                  │
│                                                 │
└────────────────────────────────────────────────┘
```

**Why This Design:**

```
✓ Business logic doesn't know about SMTP details
✓ Email service can be replaced (SMTP → SES → SendGrid)
✓ Testing doesn't require email server
✓ Email failures don't crash the application
✓ Async execution improves performance
```

### 3. Goroutine Concurrency

**Email Sending Pattern:**

```go
// In usecase layer
func (uc *UseCase) CreateProduct(product *domain.Product) error {
    // Synchronous: Save to database
    if err := uc.repo.CreateProduct(product); err != nil {
        return err
    }
    
    // Asynchronous: Send email
    go func() {
        subject := "New Product Created"
        body := fmt.Sprintf(
            "Product: %s\nCategory: %s\nPrice: $%.2f\nQuantity: %d",
            product.ProductName,
            product.Category.Name,
            product.Price,
            product.Quantity,
        )
        
        if err := uc.emailService.SendEmail(subject, body); err != nil {
            log.Printf("Email send failed: %v", err)
            // Error logged but doesn't affect API response
        }
    }()
    
    return nil // API responds immediately
}
```

**Performance Comparison:**

```
Without Goroutines (Synchronous):
├─ Validate request: 2ms
├─ Database insert: 10ms
├─ Send email: 1500ms ◄── Blocking!
└─ Total response time: 1512ms

With Goroutines (Asynchronous):
├─ Validate request: 2ms
├─ Database insert: 10ms
├─ Spawn goroutine: 0.1ms
└─ Total response time: 12.1ms ◄── 125x faster!

Background goroutine:
└─ Send email: 1500ms (doesn't block response)
```

### 4. Foreign Key Protection

**Relationship Constraints:**

```
DELETE Category (Electronics)
    │
    ├─ System checks: Are there products?
    │  └─ USB Cable exists in Electronics
    │
    ├─ REJECT: Cannot delete
    │  └─ Error: "cannot delete category: products exist"
    │
    └─ User must:
       1. Delete/move all products first
       2. Then delete category
```

**Why This Matters:**
- Prevents orphaned data (products without categories)
- Maintains referential integrity
- Enforces logical business constraints
- Protects data consistency

### 5. Smart Email Content

**Product Name Resolution:**

```
Order Request:
{
  "ProductID": "660e8400-e29b-41d4-a716-446655440111",
  "Quantity": 5
}

❌ Bad Email (just IDs):
"Order created for product 660e8400-e29b-41d4-a716-446655440111"

✓ Good Email (meaningful info):
"Order created for: USB Cable
 Quantity: 5 units
 Remaining stock: 95 units"

How we achieve this:
1. Query database for product details
2. Fetch product name and current quantity
3. Calculate remaining stock
4. Format user-friendly message
```

---

## Testing Guide

### Complete Testing Workflow

```
┌─────────────────────────────────────────────────────────┐
│              COMPREHENSIVE TEST PLAN                     │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Phase 1: Setup Verification                            │
│  ├─ [ ] PostgreSQL running                              │
│  ├─ [ ] Database created                                │
│  ├─ [ ] Tables created                                  │
│  ├─ [ ] .env file configured                            │
│  ├─ [ ] Email credentials valid                         │
│  └─ [ ] Application starts without errors               │
│                                                          │
│  Phase 2: Basic CRUD Tests                              │
│  ├─ [ ] Create category → Check email                   │
│  ├─ [ ] Create product → Check email                    │
│  ├─ [ ] Create order → Check email with product name    │
│  ├─ [ ] Verify stock reduction                          │
│  └─ [ ] Get all records                                 │
│                                                          │
│  Phase 3: Error Handling                                │
│  ├─ [ ] Order exceeding stock → 400 error               │
│  ├─ [ ] Delete category with products → 400 error       │
│  ├─ [ ] Invalid UUID → 404 error                        │
│  ├─ [ ] Invalid email config → Logged error             │
│  └─ [ ] Malformed JSON → 400 error                      │
│                                                          │
│  Phase 4: Email Verification                            │
│  ├─ [ ] All 3 notification types received               │
│  ├─ [ ] Email contains correct data                     │
│  ├─ [ ] Product names (not IDs) in order emails         │
│  ├─ [ ] Emails arrive within 10 seconds                 │
│  └─ [ ] Failed emails logged but don't break API        │
│                                                          │
│  Phase 5: Performance                                    │
│  ├─ [ ] API responds in <100ms                          │
│  ├─ [ ] Email doesn't block response                    │
│  ├─ [ ] Multiple concurrent requests work               │
│  └─ [ ] Database transactions are atomic                │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### Step-by-Step Test Script

```bash
#!/bin/bash
# save as test_api.sh

BASE_URL="http://localhost:8080/api/v1"

echo "═══════════════════════════════════════════"
echo "  INVENTORY API TEST SUITE"
echo "═══════════════════════════════════════════"

# Test 1: Create Category
echo ""
echo "Test 1: Creating category..."
CATEGORY_RESPONSE=$(curl -s -X POST $BASE_URL/categories \
  -H "Content-Type: application/json" \
  -d '{"Name":"Test Electronics"}')

CATEGORY_ID=$(echo $CATEGORY_RESPONSE | jq -r '.ID')
echo "✓ Category ID: $CATEGORY_ID"
echo "→ CHECK EMAIL: 'New Category Created'"

sleep 2

# Test 2: Create Product
echo ""
echo "Test 2: Creating product..."
PRODUCT_RESPONSE=$(curl -s -X POST $BASE_URL/products \
  -H "Content-Type: application/json" \
  -d "{
    \"ProductName\":\"Test USB Cable\",
    \"Description\":\"Test product\",
    \"Price\":15.99,
    \"Quantity\":100,
    \"CategoryID\":\"$CATEGORY_ID\",
    \"IsActive\":true
  }")

PRODUCT_ID=$(echo $PRODUCT_RESPONSE | jq -r '.ID')
echo "✓ Product ID: $PRODUCT_ID"
echo "→ CHECK EMAIL: 'New Product Added'"

sleep 2

# Test 3: Create Order
echo ""
echo "Test 3: Creating order..."
ORDER_RESPONSE=$(curl -s -X POST $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -d "{
    \"ProductID\":\"$PRODUCT_ID\",
    \"Quantity\":5
  }")

ORDER_ID=$(echo $ORDER_RESPONSE | jq -r '.ID')
echo "✓ Order ID: $ORDER_ID"
echo "→ CHECK EMAIL: 'New Order Placed' with product name"

sleep 2

# Test 4: Verify Stock Reduction
echo ""
echo "Test 4: Verifying stock..."
PRODUCTS=$(curl -s $BASE_URL/products)
CURRENT_STOCK=$(echo $PRODUCTS | jq -r ".[] | select(.ID==\"$PRODUCT_ID\") | .Quantity")
echo "✓ Current stock: $CURRENT_STOCK (should be 95)"

# Test 5: Test Insufficient Stock
echo ""
echo "Test 5: Testing insufficient stock..."
INSUFFICIENT_ORDER=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -d "{
    \"ProductID\":\"$PRODUCT_ID\",
    \"Quantity\":200
  }")

echo "$INSUFFICIENT_ORDER" | grep "HTTP_CODE:400" && echo "✓ Correctly rejected" || echo "✗ Failed"

# Test 6: Test FK Constraint
echo ""
echo "Test 6: Testing foreign key protection..."
DELETE_RESULT=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X DELETE $BASE_URL/categories/$CATEGORY_ID)
echo "$DELETE_RESULT" | grep "HTTP_CODE:400" && echo "✓ Correctly prevented" || echo "✗ Failed"

echo ""
echo "═══════════════════════════════════════════"
echo "  TEST SUMMARY"
echo "═══════════════════════════════════════════"
echo "Please verify:"
echo "  1. Three emails received in inbox"
echo "  2. Order email contains 'Test USB Cable' (not UUID)"
echo "  3. Stock reduced from 100 to 95"
echo "  4. All tests passed"
echo "═══════════════════════════════════════════"
```

Run with:
```bash
chmod +x test_api.sh
./test_api.sh
```

### Email Testing Checklist

```
Email Test Verification:
├─ [ ] Category Email
│   ├─ Subject: "New Category Created"
│   ├─ Body contains: Category name
│   └─ Received within 10 seconds
│
├─ [ ] Product Email
│   ├─ Subject: "New Product Added"
│   ├─ Body contains: Product name, price, quantity, category
│   └─ Received within 10 seconds
│
└─ [ ] Order Email
    ├─ Subject: "New Order Placed"
    ├─ Body contains: Product NAME (not UUID)
    ├─ Body contains: Quantity ordered
    ├─ Body contains: Remaining stock
    └─ Received within 10 seconds
```

---

## Troubleshooting

### Email Issues

```
┌─────────────────────────────────────────────────────────┐
│             EMAIL TROUBLESHOOTING GUIDE                  │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Problem: "Email not received"                          │
│  ├─ Check spam/junk folder                              │
│  ├─ Verify RECEIVER_EMAIL in .env                       │
│  ├─ Check application logs for errors                   │
│  └─ Test SMTP credentials manually                      │
│                                                          │
│  Problem: "535 Authentication failed"                   │
│  ├─ Gmail: Use App Password (not account password)      │
│  ├─ Enable 2FA first                                    │
│  ├─ Generate new App Password                           │
│  └─ Update SENDER_PASSWORD in .env                      │
│                                                          │
│  Problem: "Email sent but shows UUID in order"          │
│  ├─ Check product lookup in usecase layer               │
│  ├─ Verify product exists in database                   │
│  └─ Review email formatting logic                       │
│                                                          │
│  Problem: "Connection timeout to SMTP"                  │
│  ├─ Verify SMTP_HOST and SMTP_PORT                      │
│  ├─ Check firewall/network restrictions                 │
│  ├─ Try port 465 (SSL) instead of 587 (TLS)            │
│  └─ Test with telnet: telnet smtp.gmail.com 587        │
│                                                          │
│  Problem: "Email delays API response"                   │
│  ├─ Verify goroutine is used (go keyword)               │
│  ├─ Check if email logic is synchronous                 │
│  └─ Review usecase implementation                       │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

### Gmail App Password Setup

```
Step-by-Step:
├─ 1. Go to myaccount.google.com
├─ 2. Security → 2-Step Verification
├─ 3. Enable 2-Step Verification (if not enabled)
├─ 4. App passwords → Select app: Mail
├─ 5. Select device: Other (Custom name)
├─ 6. Name it: "Inventory API"
├─ 7. Copy 16-character password
├─ 8. Paste into .env as SENDER_PASSWORD
└─ 9. Restart application
```

### Database Issues

```
Problem: "connection refused"
└─► Check if PostgreSQL is running
    └─► Linux: sudo systemctl status postgresql
    └─► macOS: brew services list
    └─► Windows: Services → postgresql

Problem: "password authentication failed"
└─► Verify .env credentials match database user
    └─► Test in pgAdmin first
    └─► Check pg_hba.conf for auth method

Problem: "relation does not exist"
└─► Tables not created
    └─► Run SQL script in correct database
    └─► Verify you're connected to inventory_db

Problem: "foreign key constraint"
└─► Trying to delete referenced data
    └─► Delete child records (products) first
    └─► Then delete parent (category)
```

### Application Logs

```
Useful log messages to watch for:

✓ Success indicators:
  "Database connected successfully"
  "Email service initialized"
  "SMTP configured: smtp.gmail.com:587"
  "Server running on port 8080"
  "Email sent successfully to: receiver@example.com"

✗ Error indicators:
  "Email send failed: ..." (logged but not critical)
  "Database connection failed"
  "SMTP authentication failed"
  "Product not found for order email"
```

---

## Performance Benchmarks

```
┌─────────────────────────────────────────────────────────┐
│              EXPECTED PERFORMANCE METRICS                │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  API Response Times (with email):                       │
│  ├─ CREATE Category:    10-20ms                         │
│  ├─ CREATE Product:     15-30ms                         │
│  ├─ CREATE Order:       20-40ms                         │
│  ├─ GET All Products:   5-15ms                          │
│  └─ UPDATE Product:     10-25ms                         │
│                                                          │
│  Email Delivery Times:                                   │
│  ├─ Local SMTP:         100-500ms                       │
│  ├─ Gmail SMTP:         500-2000ms                      │
│  ├─ SendGrid:           200-800ms                       │
│  └─ AWS SES:            300-1000ms                      │
│                                                          │
│  Note: Email times don't affect API response            │
│        because of async goroutines                      │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

---

## Security Considerations

```
┌─────────────────────────────────────────────────────────┐
│                 SECURITY BEST PRACTICES                  │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  Environment Variables:                                  │
│  ├─ NEVER commit .env to git                            │
│  ├─ Use different credentials per environment           │
│  ├─ Rotate SMTP passwords regularly                     │
│  └─ Use read-only database users where possible         │
│                                                          │
│  Email Security:                                         │
│  ├─ Use App Passwords (not account passwords)           │
│  ├─ Enable TLS for SMTP (port 587)                      │
│  ├─ Don't log email passwords                           │
│  └─ Validate recipient email addresses                  │
│                                                          │
│  Database Security:                                      │
│  ├─ Use prepared statements (GORM does this)            │
│  ├─ Limit user permissions (GRANT specific tables)      │
│  ├─ Enable SSL for production database connections      │
│  └─ Regular backups                                     │
│                                                          │
│  API Security (Future Enhancements):                    │
│  ├─ Add authentication (JWT tokens)                     │
│  ├─ Rate limiting                                       │
│  ├─ Input validation and sanitization                   │
│  └─ HTTPS in production                                 │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

---

## Future Enhancements

```
Planned Features:
├─ Email Templates
│  └─ HTML email templates with branding
│
├─ Email Queue System
│  └─ Redis/RabbitMQ for reliable delivery
│
├─ Notification Preferences
│  └─ User-configurable notification settings
│
├─ Multiple Recipients
│  └─ CC/BCC support for team notifications
│
├─ Email Retry Logic
│  └─ Automatic retry on failure
│
├─ Webhook Support
│  └─ Alternative to email for integrations
│
└─ Email Analytics
   └─ Track delivery rates and open rates
```

---

## Quick Reference

```bash
# ============================================
# START SERVICES
# ============================================

# PostgreSQL
sudo systemctl start postgresql    # Linux
brew services start postgresql     # macOS
# Services → postgresql             # Windows

# Application
go run cmd/main.go


# ============================================
# TEST EMAIL CONFIGURATION
# ============================================

# Test SMTP connection (manual)
telnet smtp.gmail.com 587

# Create test category (triggers email)
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{"Name":"Test Category"}'

# Check application logs
tail -f application.log  # or check console output


# ============================================
# COMMON OPERATIONS
# ============================================

# View all data
curl http://localhost:8080/api/v1/get-tasks

# Check specific product stock
curl http://localhost:8080/api/v1/products | jq '.[] | {name: .ProductName, stock: .Quantity}'

# View recent orders
curl http://localhost:8080/api/v1/orders | jq '.[-5:]'


# ============================================
# DEBUGGING
# ============================================

# Check if email goroutine is spawned
# Look for log: "Sending email notification..."

# Verify email service initialization
# Look for log: "Email service initialized"
# Look for log: "SMTP configured: smtp.gmail.com:587"

# Test database connection
psql -h localhost -U inventory_user -d inventory_db

# Check table contents
psql -h localhost -U inventory_user -d inventory_db \
  -c "SELECT * FROM products;"
```

---

## License

MIT License

---

## Support

For issues or questions:

1. Check the [Troubleshooting](#troubleshooting) section
2. Review application logs for errors
3. Verify email credentials and SMTP settings
4. Test with simple curl commands
5. Check Gmail App Password setup

---

**Built with Go, PostgreSQL, and SMTP**

**Features:** RESTful API • Real-time Stock Management • Email Notifications • Concurrent Processing • Database Integrity