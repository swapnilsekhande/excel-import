# Excel Import (Go + MySQL + Redis)

This project is a simple and efficient backend service written in Go that allows importing employee data from an Excel file into a MySQL database. It also leverages Redis for caching employee details to improve performance and reduce database load.

## 🔧 Features

- ✅ Import employee data from Excel (`.xlsx`) files
- ✅ Store employee data in MySQL using GORM
- ✅ Cache employee data in Redis for fast retrieval
- ✅ RESTful APIs using Gin framework
- ✅ Pagination support for employee listings
- ✅ Sync Redis cache to MySQL and clear cache
- ✅ Structured and modular codebase

---

## 📦 Tech Stack

- **Go** (Golang)
- **Gin** (Web framework)
- **GORM** (ORM for MySQL)
- **Redis** (In-memory cache)
- **go-redis** (Redis client)
- **github.com/xuri/excelize/v2** (Excel import)
- **MySQL** (Relational database)

---


---

## ⚙️ Tech Stack

- **Golang** — backend language
- **Gin** — web framework
- **GORM** — ORM for MySQL
- **Redis** — in-memory cache
- **go-redis** — Redis client
- **excelize** — Excel import library

---

## 📦 Getting Started

### 1. Clone the repository

bash

```git clone https://github.com/swapnilsekhande/excel-import.git```
```cd excel-import```

Update MySQL and Redis configuration inside database/mysql.go and database/redis.go.

```go run main.go```
