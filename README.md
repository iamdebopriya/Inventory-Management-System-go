# Inventory Management System API

A RESTful API service for managing inventory operations with automatic stock management and real-time notifications.

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-316192?style=flat&logo=postgresql)](https://www.postgresql.org/)

---

## Table of Contents

- [What This API Does](#what-this-api-does)
- [System Architecture](#system-architecture)
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
- [API Usage](#api-usage)
- [Project Structure](#project-structure)
- [Key Features Explained](#key-features-explained)

---

## What This API Does

This inventory management system helps you:

- **Organize Products** into categories (Electronics, Clothing, etc.)
- **Track Stock Levels** in real-time
- **Process Orders** with automatic inventory updates
- **Maintain Data Integrity** through foreign key relationships
- **Send Notifications** asynchronously without blocking operations

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
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│         │                  │                  │                  │
│         │    Validation    │    Stock Check   │                 │
│         │    Rules         │    Goroutines    │                 │
└─────────┼──────────────────┼──────────────────┼─────────────────┘
          │                  │                  │
          ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                    REPOSITORY LAYER (GORM)                       │
│                     Database Abstraction                         │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      POSTGRESQL DATABASE                         │
│   ┌────────────┐  ┌────────────┐  ┌────────────┐               │
│   │ Categories │  │  Products  │  │   Orders   │               │
│   │   Table    │  │   Table    │  │   Table    │               │
│   └────────────┘  └────────────┘  └────────────┘               │
└─────────────────────────────────────────────────────────────────┘
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
└─────────────────────┘                    │
                                           │
┌─────────────────────┐                    │
│      ORDERS         │                    │
├─────────────────────┤                    │
│ id (UUID) PK        │                    │
│ product_id FK ──────┼────────────────────┘
│ quantity (INT)      │
│ order_date          │
└─────────────────────┘
```

### Request Flow Diagram

```
POST /api/v1/orders
│
├─ Step 1: Validate Request
│  └─ Check JSON format
│     Check required fields
│
├─ Step 2: Find Product
│  └─ Query database for ProductID
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
├─ Step 5: Send Notification (Goroutine)
│  └─ Spawn async process
│     Log notification
│     (Non-blocking)
│
└─ Step 6: Return Response
   └─ 201 Created with order details
```

---

## Technology Stack

```
┌─────────────────────────────────────────────────┐
│              Application Layer                   │
│  ┌───────────────────────────────────────┐     │
│  │  Go 1.19+                              │     │
│  │  - Gin Web Framework (HTTP)            │     │
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
                      ▼
┌─────────────────────────────────────────────────┐
│                Database Layer                    │
│  ┌───────────────────────────────────────┐     │
│  │  PostgreSQL 13+                        │     │
│  │  - UUID Extension                      │     │
│  │  - Foreign Key Constraints             │     │
│  │  - Indexes                             │     │
│  └───────────────────────────────────────┘     │
└─────────────────────────────────────────────────┘
```

---

## Getting Started

### Prerequisites Checklist

- [ ] Go 1.19 or higher installed
- [ ] PostgreSQL 13 or higher installed
- [ ] pgAdmin 4 (for database management)
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

#### Option A: Quick Setup (Recommended for Beginners)

1. **Open pgAdmin 4**
2. **Create Database:** Right-click "Databases" → Create → Database
   - Name: `inventory_db`
   - Click Save
3. **Run Setup Script:** Select `inventory_db` → Query Tool → Paste and run:
   ```sql
   -- See documentation for full SQL script
   ```

#### Option B: Auto-Migration (For Development)

Use GORM's auto-migration feature:

```go
// In database/database.go
db.AutoMigrate(&domain.Category{}, &domain.Product{}, &domain.Order{})
```

**Setup Flowchart:**

```
START
  ├─► Install PostgreSQL
  ├─► Create Database (inventory_db)
  ├─► Enable UUID Extension
  ├─► Create Application User
  ├─► Choose Setup Method:
  │   ├─► Manual: Run SQL scripts
  │   └─► Auto: Enable AutoMigrate in code
  ├─► Grant Permissions
  └─► Update .env file
END
```

### Step 3: Configuration

Create `.env` file in project root:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=inventory_user
DB_PASS=your_secure_password
DB_NAME=inventory_db
SERVER_PORT=8080
```

### Step 4: Run the Application

```bash
# Run directly
go run cmd/main.go

# Expected output:
# Database connected successfully
# Server running on port 8080...
```

---

## API Usage

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoint Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      API ENDPOINTS                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  CATEGORIES                                                  │
│  ├─ GET    /categories          List all                    │
│  ├─ POST   /categories          Create new                  │
│  ├─ PUT    /categories/{id}     Update                      │
│  └─ DELETE /categories/{id}     Delete                      │
│                                                              │
│  PRODUCTS                                                    │
│  ├─ GET    /products            List all                    │
│  ├─ POST   /products            Create new                  │
│  ├─ PUT    /products/{id}       Update                      │
│  └─ DELETE /products/{id}       Delete                      │
│                                                              │
│  ORDERS                                                      │
│  ├─ GET    /orders              List all                    │
│  └─ POST   /orders              Create (auto-deduct stock)  │
│                                                              │
│  AGGREGATED                                                  │
│  ├─ GET    /get-tasks           All data at once            │
│  └─ POST   /get-by-tablename    Get by table & ID           │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Usage Workflow Diagram

```
Step 1: Create Category
    POST /categories
    {"Name": "Electronics"}
         │
         ├─► Returns: {"ID": "uuid-123...", "Name": "Electronics"}
         │
         ▼
Step 2: Create Product
    POST /products
    {
      "ProductName": "USB Cable",
      "CategoryID": "uuid-123...",  ◄── Use ID from Step 1
      "Price": 15.99,
      "Quantity": 100
    }
         │
         ├─► Returns: Product with ID "uuid-456..."
         │
         ▼
Step 3: Create Order
    POST /orders
    {
      "ProductID": "uuid-456...",  ◄── Use ID from Step 2
      "Quantity": 5
    }
         │
         ├─► Stock automatically reduced: 100 → 95
         │
         ▼
Step 4: Verify
    GET /products
    ├─► Check quantity is now 95
    │
    GET /orders
    └─► See order history
```

### Quick Test Commands

```bash
# 1. Create category
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Content-Type: application/json" \
  -d '{"Name":"Electronics"}'

# 2. Create product (replace CATEGORY_ID with actual UUID)
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "ProductName":"USB Cable",
    "Price":15.99,
    "Quantity":100,
    "CategoryID":"CATEGORY_ID"
  }'

# 3. Create order (replace PRODUCT_ID with actual UUID)
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"ProductID":"PRODUCT_ID","Quantity":5}'

# 4. Get all data
curl http://localhost:8080/api/v1/get-tasks
```

---

## Project Structure

```
inventory-api/
│
├── cmd/
│   └── main.go                    # Application entry point
│
├── database/
│   └── database.go                # Database connection
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
│   │
│   ├── usecase/
│   │   └── usecase.go             # Business logic
│   │                              #    - Validation
│   │                              #    - Stock management
│   │                              #    - Goroutines
│   │
│   └── delivery/
│       └── http/
│           ├── handler.go         # HTTP handlers
│           └── routes.go          # Route definitions
│
├── .env                           # Environment variables
├── go.mod                         # Dependencies
├── go.sum                         # Checksums
└── README.md                      # This file
```

### Architecture Layers

```
┌──────────────────────────────────────────────────┐
│              DELIVERY LAYER                       │
│  (HTTP Handlers, Routes, JSON Serialization)     │
│                                                   │
│  Responsibilities:                                │
│  • Parse HTTP requests                           │
│  • Call use case functions                       │
│  • Return JSON responses                         │
│  • Handle HTTP status codes                      │
└─────────────────┬────────────────────────────────┘
                  │
                  ▼
┌──────────────────────────────────────────────────┐
│              USE CASE LAYER                       │
│  (Business Logic, Validation, Orchestration)     │
│                                                   │
│  Responsibilities:                                │
│  • Validate business rules                       │
│  • Check stock availability                      │
│  • Spawn goroutines                              │
│  • Coordinate operations                         │
└─────────────────┬────────────────────────────────┘
                  │
                  ▼
┌──────────────────────────────────────────────────┐
│            REPOSITORY LAYER                       │
│  (Database Access, GORM Operations)              │
│                                                   │
│  Responsibilities:                                │
│  • Execute SQL queries                           │
│  • Handle transactions                           │
│  • Preload relationships                         │
│  • Map database ↔ Go structs                     │
└─────────────────┬────────────────────────────────┘
                  │
                  ▼
┌──────────────────────────────────────────────────┐
│              DATABASE LAYER                       │
│  (PostgreSQL, Tables, Constraints)               │
│                                                   │
│  Responsibilities:                                │
│  • Store data persistently                       │
│  • Enforce foreign keys                          │
│  • Maintain indexes                              │
│  • Ensure ACID properties                        │
└──────────────────────────────────────────────────┘
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
└─ Response: Order created, stock now 95
```

**Business Rules:**
- Orders cannot exceed available stock
- Stock updates are atomic (all or nothing)
- Real-time inventory tracking

### 2. Foreign Key Protection

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
- Maintains data integrity
- Enforces logical business constraints

### 3. Asynchronous Notifications

**Goroutine Workflow:**

```
POST /products (Create USB Cable)
    │
    ├─ Main Thread:
    │  ├─ Validate request
    │  ├─ Insert into database
    │  ├─ Spawn goroutine ────────┐
    │  └─ Return HTTP response    │
    │     (Total: ~10ms)           │
    │                              │
    │                              ├─ Background Thread:
    │                              │  ├─ Wait 1 second
    │                              │  ├─ Log notification
    │                              │  └─ Complete
    │                              │     (Total: ~1000ms)
    └─ User receives response      │
       immediately                 └─ Runs independently
```

**Benefits:**
- **Fast responses** - API returns instantly
- **Scalable** - Thousands of notifications don't slow requests
- **Non-blocking** - Server handles more concurrent users

### 4. Data Integrity

**Field Naming Convention:**

```
┌──────────────────┬──────────────────┬──────────────────┐
│   Go Struct      │   JSON Request   │   Database       │
├──────────────────┼──────────────────┼──────────────────┤
│ ProductName      │ "ProductName"    │ product_name     │
│ Description      │ "Description"    │ description      │
│ CategoryID       │ "CategoryID"     │ category_id      │
│ IsActive         │ "IsActive"       │ is_active        │
└──────────────────┴──────────────────┴──────────────────┘
```

**Important:**
- JSON fields are **case-sensitive** (use PascalCase)
- Database columns use **snake_case**
- GORM handles conversion automatically

---

## Error Handling Reference

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP STATUS CODES                         │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  200 OK                    Successful GET/PUT/DELETE        │
│  201 Created               Successful POST                  │
│  400 Bad Request           Validation error                 │
│                            Insufficient stock               │
│                            FK constraint violation          │
│  404 Not Found             Resource doesn't exist           │
│  500 Internal Server Error Database/system failure          │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### Common Errors

| Error Message | Cause | Solution |
|---------------|-------|----------|
| `product not found` | Invalid ProductID | Use existing product UUID |
| `product quantity insufficient` | Order > stock | Reduce order quantity |
| `cannot delete category: products exist` | FK constraint | Delete products first |
| `record not found` | Invalid UUID | Check ID exists |

---

## Best Practices

### Testing Your API

**Recommended Test Flow:**

```
Test 1: Happy Path
├─ Create category → Success
├─ Create product → Success
├─ Create order → Success
└─ Verify stock reduced → Success

Test 2: Error Scenarios
├─ Order with insufficient stock → 400 Error
├─ Delete category with products → 400 Error
├─ Invalid UUID → 404 Error
└─ Duplicate category name → 500 Error

Test 3: Edge Cases
├─ Order quantity = 0 → Should fail
├─ Negative price → Should fail
├─ Empty product name → Should fail
└─ Order after stock depleted → Should fail
```

### Development Tips

- **Always save UUIDs** from responses for subsequent requests
- **Check logs** for goroutine notifications
- **Test error cases** to understand constraints
- **Use pgAdmin** to verify database state
- **Start fresh** by deleting test data regularly

---

## Quick Reference

```bash
# Start PostgreSQL
sudo systemctl start postgresql    # Linux
brew services start postgresql     # macOS
# Services → postgresql             # Windows

# Run application
go run cmd/main.go

# Test endpoints
curl http://localhost:8080/api/v1/categories
curl http://localhost:8080/api/v1/products
curl http://localhost:8080/api/v1/orders

# View all data
curl http://localhost:8080/api/v1/get-tasks
```

---

## Troubleshooting

### Connection Issues

```
Problem: "connection refused"
└─► Check if PostgreSQL is running
    └─► Start PostgreSQL service

Problem: "password authentication failed"
└─► Verify .env credentials match pgAdmin
    └─► Test connection in pgAdmin first

Problem: "permission denied for table"
└─► Grant permissions to inventory_user
    └─► Run GRANT statements in pgAdmin
```

### Database Issues

```
Problem: "relation does not exist"
└─► Tables not created in correct database
    └─► Run CREATE TABLE scripts in inventory_db

Problem: "column name mismatch"
└─► Using wrong case in JSON
    └─► Use PascalCase: ProductName, CategoryID

Problem: "foreign key constraint"
└─► Trying to delete referenced data
    └─► Delete child records first
```

---

## License

MIT License

---

**Built with Go and PostgreSQL**