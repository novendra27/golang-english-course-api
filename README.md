# 🚀 English Course Registration API

[![Go Version](https://img.shields.io/badge/Go-1.22%20%7C%201.26-00ADD8?style=flat&logo=go)](https://golang.org)
[![Gin Framework](https://img.shields.io/badge/Framework-Gin%20v1.12-008ECF?style=flat&logo=gin)](https://gin-gonic.com)
[![GORM](https://img.shields.io/badge/ORM-GORM%20v1.31-7952B3?style=flat)](https://gorm.io)
[![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL%2016-336791?style=flat&logo=postgresql)](https://www.postgresql.org)
[![Zerolog](https://img.shields.io/badge/Logging-Zerolog-brightgreen?style=flat)](https://github.com/rs/zerolog)
[![i18n](https://img.shields.io/badge/Localization-go--i18n%20(EN%20%2F%20ID)-blue?style=flat)](https://github.com/nicksnyder/go-i18n)
[![Swagger](https://img.shields.io/badge/API%20Docs-Swagger%202.0-85EA2D?style=flat&logo=swagger)](http://localhost:8080/swagger/index.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> [!NOTE]
> 📖 **Bahasa Indonesia:** Tersedia dokumentasi dalam Bahasa Indonesia di [README.id.md](file:///d:/Code%20Learning/Learn%20Project/Golang%20Project/english-course-api/README.id.md).

A production-grade RESTful API backend for an English Course Registration and Class Placement management system. Built with **Go (Golang)** following clean **Layered Architecture** (*Separation of Concerns*), featuring **Dynamic i18n Localization**, strict transactional state management, automated unit tests, and full Docker containerization.

---

## 📑 Table of Contents

- [Key Features](#-key-features)
- [Tech Stack](#%EF%B8%8F-tech-stack)
- [System Architecture & Directory Structure](#-system-architecture--directory-structure)
- [Business Workflow & State Machine](#-business-workflow--state-machine)
- [Dynamic Localization (i18n)](#-dynamic-localization-i18n)
- [Getting Started](#-getting-started)
  - [Option 1: Docker Compose (Recommended)](#option-1-docker-compose-recommended)
  - [Option 2: Direct Local Go Execution](#option-2-direct-local-go-execution)
- [5-Minute Quick Tour (End-to-End API Walkthrough)](#-5-minute-quick-tour-end-to-end-api-walkthrough)
- [Environment Configuration](#-environment-configuration)
- [Automated Testing](#-automated-testing)
- [API Documentation & Endpoints](#-api-documentation--endpoints)
  - [Interactive Swagger UI](#interactive-swagger-ui)
  - [Standard JSON API Response Format](#standard-json-api-response-format)
  - [Endpoints Overview](#endpoints-overview)
- [License](#-license)

---

## 🌟 Key Features

- **Clean Layered Architecture:** Strict decoupled boundaries between Handler $\rightarrow$ Service $\rightarrow$ Repository $\rightarrow$ Database.
- **Transactional State Machine:** Course Registration, Atomic Payment Settlement, and Capacity-Guarded Class Placement.
- **Dynamic Localization Engine (i18n):** Embedded bilingual dictionaries (`EN` & `ID`) supporting runtime switching via query parameter (`?lang=`) or HTTP header (`Accept-Language`).
- **Data Integrity & Validation:** Field-level struct validation via `validator/v10` with human-readable localized error messages.
- **High-Performance Observability:** Structured JSON request logging powered by `Zerolog` with latency, IP, and status tracing.
- **Interactive OpenAPI / Swagger Documentation:** Live Swagger 2.0 UI embedded for frictionless API testing.
- **Zero-Dependency Production Builds:** Multi-stage Docker containerization ready for cloud deployments.

---

## 🛠️ Tech Stack

| Component | Technology | Description |
| :--- | :--- | :--- |
| **Language** | Go `1.22+` / `1.26+` | High performance, strict type safety, native concurrency. |
| **Web Framework** | [Gin Web Framework](https://github.com/gin-gonic/gin) | Fast HTTP routing, middleware chaining, JSON binding. |
| **ORM** | [GORM](https://gorm.io) | Schema migration, relation preloading, ACID transactions. |
| **Database** | PostgreSQL 16 | Relational data persistence with strict foreign key constraints. |
| **Localization** | [go-i18n v2](https://github.com/nicksnyder/go-i18n) | Dynamic runtime translation with embedded JSON files. |
| **Logging** | [Zerolog](https://github.com/rs/zerolog) | Blazing-fast zero-allocation JSON structured logger. |
| **Validation** | Go Playground Validator v10 | Struct validation tags with multi-language error parsing. |
| **API Docs** | [Swaggo / Swagger 2.0](https://github.com/swaggo/swag) | Declarative Swagger documentation generator & Swagger UI. |
| **DevOps** | Docker & Docker Compose | Multi-stage scratch/alpine build container orchestration. |

---

## 📐 System Architecture & Directory Structure

```text
english-course-api/
├── config/             # Config loader (.env) & PostgreSQL connection pool setup
│   ├── config.go
│   └── database.go
├── docs/               # Auto-generated Swagger 2.0 specs (docs.go, swagger.json, swagger.yaml)
├── handlers/           # HTTP controllers: request parsing, DTO binding, status responses
│   ├── student_handler.go
│   ├── course_handler.go
│   ├── class_handler.go
│   ├── registration_handler.go
│   ├── payment_handler.go
│   └── class_placement_handler.go
├── locales/            # Embedded translation dictionaries (//go:embed *.json)
│   ├── locales.go
│   ├── en.json
│   └── id.json
├── middleware/         # Gin middleware pipeline (Zerolog HTTP logger, i18n detection)
│   ├── logger.go
│   └── i18n.go
├── models/             # Domain entities & GORM relational models
│   ├── student.go
│   ├── course.go
│   ├── class.go
│   ├── registration.go
│   ├── payment.go
│   └── class_placement.go
├── repositories/       # Data Access Layer: direct GORM queries and transactions
│   ├── student_repository.go
│   ├── course_repository.go
│   ├── class_repository.go
│   ├── registration_repository.go
│   ├── payment_repository.go
│   └── class_placement_repository.go
├── routes/             # Dependency injection wiring & Gin router endpoint group setup
│   └── routes.go
├── services/           # Business Logic Layer: state validations, transaction coordination
│   ├── student_service.go
│   ├── course_service.go
│   ├── class_service.go
│   ├── registration_service.go
│   ├── registration_service_test.go
│   ├── payment_service.go
│   ├── class_placement_service.go
│   └── class_placement_service_test.go
├── utils/              # Standard JSON API response helper & i18n translation engine
│   ├── i18n.go
│   ├── i18n_test.go
│   └── response.go
├── Dockerfile          # Multi-stage container build definition
├── docker-compose.yml  # PostgreSQL & Go API orchestration
├── .env.example        # Environment variable configuration template
├── go.mod / go.sum     # Go dependency management
└── main.go             # Application bootstrap & entrypoint
```

---

## 🔄 Business Workflow & State Machine

```text
Student + Course ──► Registration (Status: pending)
                            │
                            ▼
                    Payment (Status: pending ──► paid)
                            │ (Invoice settled)
                            ▼
                Class Placement (Assigned to Class)
                            │
                            ▼
            Class (Status: open ──► full when capacity reached)
```

### 📋 Core Business Rules:
1. **Unique Email Enforcement:** Each student record must contain a unique email address.
2. **Duplicate Registration Prevention:** A student cannot register for the same course if an active registration (`pending` or `registered`) exists.
3. **Atomic Billing:** When a registration is submitted, a linked `Payment` invoice is created atomically.
4. **Payment Gate:** A student **can only be assigned to a class if their registration status is `registered`** (meaning the invoice is `paid`).
5. **Course Matching:** The class selected for placement must match the exact course of the student's registration.
6. **Class Capacity Guard:** Placements are rejected if the target class is `closed` or has reached its maximum seat `capacity`. When filled, the class status transitions to `full`.

---

## 🌐 Dynamic Localization (i18n)

The API dynamically resolves response messages and validation errors based on client preference:
- **Default Language:** English (`en`)
- **Supported Locales:** `en` (English), `id` (Indonesian)

### How to specify the language:
1. **Query Parameter (Highest priority):**
   ```bash
   GET /api/v1/courses?lang=id
   ```
2. **HTTP Request Header:**
   ```bash
   Accept-Language: id-ID,id;q=0.9,en;q=0.8
   ```

---

## 🚀 Getting Started

### Option 1: Docker Compose (Recommended)

Run both the PostgreSQL database and the API container seamlessly:

```bash
# 1. Build and start containers in the background
docker compose up --build -d

# 2. Follow live structured application logs
docker compose logs -f app

# 3. Stop containers
docker compose down
```

- **API Base URL:** `http://localhost:8080/api/v1`
- **Swagger UI:** `http://localhost:8080/swagger/index.html`
- **Health Check:** `http://localhost:8080/health`

---

### Option 2: Direct Local Go Execution

#### Prerequisites:
- Go `1.22+` installed
- PostgreSQL `14+` running locally

#### Steps:
1. **Clone the repository:**
   ```bash
   git clone https://github.com/novendra27/golang-english-course-api.git
   cd golang-english-course-api
   ```
2. **Setup environment variables:**
   ```bash
   # Linux/macOS
   cp .env.example .env

   # Windows PowerShell
   copy .env.example .env
   ```
3. **Run database migrations & start server:**
   ```bash
   go run main.go
   ```

---

## ⚡ 5-Minute Quick Tour (End-to-End API Walkthrough)

Experience the complete lifecycle of student registration, invoice settlement, and class placement in under 5 minutes using these sequential `curl` commands:

### Step 1: Register a Student
```bash
curl -X POST http://localhost:8080/api/v1/students \
  -H "Content-Type: application/json" \
  -d '{"name": "Alex Johnson", "email": "alex@example.com", "phone": "081234567890"}'
```
> **Result:** Creates Student record (ID: `1`).

### Step 2: Create Course & Open Class Section
```bash
# 2a. Create Course
curl -X POST http://localhost:8080/api/v1/courses \
  -H "Content-Type: application/json" \
  -d '{"name": "IELTS Masterclass", "description": "Intensive IELTS band 7.5+ prep", "price": 1250000, "duration": "2 Months", "status": "active"}'

# 2b. Open Class Section (Capacity: 15)
curl -X POST http://localhost:8080/api/v1/classes \
  -H "Content-Type: application/json" \
  -d '{"course_id": 1, "name": "IELTS Weekend Intensive", "capacity": 15, "schedule": "Sat & Sun, 10:00 - 13:00", "status": "open"}'
```
> **Result:** Creates Course (ID: `1`) and Class (ID: `1`).

### Step 3: Enroll Student into Course (Generates Invoice)
```bash
curl -X POST http://localhost:8080/api/v1/registrations \
  -H "Content-Type: application/json" \
  -d '{"student_id": 1, "course_id": 1}'
```
> **State Machine:** Registration created with status `pending` and linked Payment ID `1` (Amount: `1250000`).

### Step 4: Settle Payment Invoice
```bash
curl -X POST http://localhost:8080/api/v1/payments/1/pay \
  -H "Content-Type: application/json" \
  -d '{"payment_method": "bank_transfer", "amount": 1250000}'
```
> **State Machine:** Payment status transitions to `paid` $\rightarrow$ Registration status automatically unlocks to `registered` 🎉.

### Step 5: Assign Student to Class Section
```bash
curl -X POST http://localhost:8080/api/v1/class-placements \
  -H "Content-Type: application/json" \
  -d '{"registration_id": 1, "class_id": 1}'
```
> **Business Rule Check:** System verifies course alignment and capacity before assigning.

### Step 6: Test Dynamic Localization (i18n)
```bash
# 6a. Default English response
curl -X GET http://localhost:8080/api/v1/courses/1

# 6b. Indonesian localized response via query parameter
curl -X GET "http://localhost:8080/api/v1/courses/1?lang=id"
```


---

## ⚙️ Environment Configuration

Configuration is loaded from `.env` using `godotenv`:

| Key | Default | Description |
| :--- | :--- | :--- |
| `APP_NAME` | `english-course-api` | Name identifier for application logs |
| `APP_ENV` | `development` | Environment mode (`development` / `production`) |
| `APP_PORT` | `8080` | HTTP port listener |
| `DB_HOST` | `localhost` | PostgreSQL host address |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database username |
| `DB_PASSWORD`| `postgres` | Database password |
| `DB_NAME` | `english_course_db`| Database name |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL connection mode |
| `DB_TIMEZONE`| `Asia/Jakarta` | Database session timezone |
| `LOG_LEVEL` | `debug` | Zerolog log level (`debug`, `info`, `warn`, `error`) |
| `LOG_PRETTY` | `true` | Enable colorized pretty console log output |

---

## 🧪 Automated Testing

Execute all unit test suites across service business logic and the i18n translation engine:

```bash
# Run all tests with verbose output
go test -v ./...
```

---

## 📚 API Documentation & Endpoints

### Interactive Swagger UI
Access the interactive API documentation and test endpoints directly from your browser:
🔗 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

---

### Standard JSON API Response Format

All endpoints return a uniform response envelope:

#### Success Response (`200 OK` / `201 Created`):
```json
{
  "success": true,
  "message": "Student created successfully",
  "data": {
    "id": 1,
    "name": "Budi Santoso",
    "email": "budi.santoso@example.com",
    "phone": "081234567890",
    "created_at": "2026-09-01T08:00:00Z",
    "updated_at": "2026-09-01T08:00:00Z"
  }
}
```

#### Localized Validation Error (`422 Unprocessable Entity`):
```json
{
  "success": false,
  "message": "Request validation failed",
  "errors": {
    "Email": "The Email field must be a valid email address",
    "Name": "The Name field is required"
  }
}
```

#### Business Conflict / Error Response (`409 Conflict` / `400 Bad Request`):
```json
{
  "success": false,
  "message": "Student already has an active registration for this course",
  "errors": null
}
```

---

### Endpoints Overview

#### 1. System Health Checks
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/health` | Root application health check |
| `GET` | `/api/v1/health` | API v1 subsystem health check |

#### 2. Student Module
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/students` | Register a new student profile |
| `GET` | `/api/v1/students` | Fetch all registered students |
| `GET` | `/api/v1/students/:id` | Fetch student details by ID |
| `PUT` | `/api/v1/students/:id` | Update student profile |
| `DELETE` | `/api/v1/students/:id` | Delete student profile |
| `GET` | `/api/v1/students/:id/registrations` | Fetch student registration history |

#### 3. Course Module
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/courses` | Create a new course offering |
| `GET` | `/api/v1/courses` | Fetch all available courses |
| `GET` | `/api/v1/courses/:id` | Fetch course detail with preloaded classes |
| `PUT` | `/api/v1/courses/:id` | Update course details |
| `DELETE` | `/api/v1/courses/:id` | Delete course |
| `GET` | `/api/v1/courses/:id/registrations` | Fetch all course enrollments |

#### 4. Class Module
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/classes` | Create a class section under a course |
| `GET` | `/api/v1/classes` | Fetch all classes with course relations |
| `GET` | `/api/v1/classes/:id` | Fetch class details by ID |
| `PUT` | `/api/v1/classes/:id` | Update class schedule or capacity |
| `DELETE` | `/api/v1/classes/:id` | Delete class |
| `GET` | `/api/v1/classes/:id/students` | List all enrolled students in a class |

#### 5. Course Registration Module
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/registrations` | Submit enrollment & auto-generate invoice |
| `GET` | `/api/v1/registrations` | Fetch all course registrations |
| `GET` | `/api/v1/registrations/:id` | Fetch registration details and payment status |
| `PUT` | `/api/v1/registrations/:id/cancel` | Cancel a pending course registration |

#### 6. Payment Module
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/payments` | Fetch all invoices |
| `GET` | `/api/v1/payments/:id` | Fetch invoice details by ID |
| `POST` | `/api/v1/payments/:id/pay` | Settle payment invoice and activate registration |

#### 7. Class Placement Module
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/class-placements` | Assign a paid student to a verified class |
| `GET` | `/api/v1/class-placements` | Fetch all student class assignments |
| `GET` | `/api/v1/class-placements/:id` | Fetch class placement details by ID |

---

## 📄 License

Distributed under the **MIT License**. See `LICENSE` for more information.
