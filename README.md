# 🚀 English Course Registration API

RESTful API backend untuk sistem pendaftaran kursus bahasa Inggris yang dibangun menggunakan **Go (Golang)** dengan arsitektur modular berlapis (*Clean Layered Architecture*), **Gin Framework**, **GORM ORM**, **PostgreSQL Database**, **Zerolog Structured Logging**, dan **Docker Containerization**.

---

## 📑 Daftar Isi

- [Tech Stack & Technologies](#-tech-stack--technologies)
- [Arsitektur & Struktur Folder](#-arsitektur--struktur-folder)
- [Alur Bisnis Utama (Core Workflow)](#-alur-bisnis-utama-core-workflow)
- [Panduan Menjalankan Aplikasi](#-panduan-menjalankan-aplikasi)
- [Automated Unit Testing](#-automated-unit-testing)
- [Format Standar Response API](#-format-standar-response-api)
- [Dokumentasi Lengkap API & Contoh Request](#-dokumentasi-lengkap-api--contoh-request)
  - [1. Health Checks](#1-health-check-endpoints)
  - [2. Modul Student](#2-modul-student-peserta-kursus)
  - [3. Modul Course](#3-modul-course-katalog-kursus)
  - [4. Modul Class](#4-modul-class-kelas-kursus)
  - [5. Modul Course Registration](#5-modul-course-registration-pendaftaran)
  - [6. Modul Payment (Simulasi Pembayaran)](#6-modul-payment-simulasi-pembayaran)
  - [7. Modul Class Placement (Penempatan Kelas)](#7-modul-class-placement-penempatan-kelas)

---

## 🛠️ Tech Stack & Technologies

| Komponen | Teknologi | Keterangan |
| :--- | :--- | :--- |
| **Language** | Go (Golang `1.22+` / `1.26+`) | Strongly-typed, performa tinggi, native concurrency. |
| **Web Framework** | [Gin Web Framework](https://github.com/gin-gonic/gin) | HTTP router cepat, JSON binding, middleware pipeline. |
| **ORM** | [GORM](https://gorm.io) | Relational mapping, query builder, database transaction. |
| **Database** | PostgreSQL 16 | Relational database dengan constraint & foreign key. |
| **Logging** | [Zerolog](https://github.com/rs/zerolog) | High-performance structured JSON logger. |
| **Validation** | Validator v10 | Struct tag validation dengan human-readable error formatter. |
| **Env Loader** | [godotenv](https://github.com/joho/godotenv) | Environment variable reader dari file `.env`. |
| **DevOps** | Docker & Docker Compose | Multi-stage build containerization. |

---

## 📐 Arsitektur & Struktur Folder

Project menerapkan pola **Layered Architecture** (*Separation of Concerns*):

```text
english-course-api/
├── config/             # Inisialisasi database GORM & load environment .env
│   ├── config.go
│   └── database.go
├── models/             # Definisi struct entitas domain & GORM tags
│   ├── student.go
│   ├── course.go
│   ├── class.go
│   ├── registration.go
│   ├── payment.go
│   └── class_placement.go
├── repositories/       # Layer akses database murni via GORM
│   ├── student_repository.go
│   ├── course_repository.go
│   ├── class_repository.go
│   ├── registration_repository.go
│   ├── payment_repository.go
│   └── class_placement_repository.go
├── services/           # Layer logika bisnis, validasi aturan & transaksi DB
│   ├── student_service.go
│   ├── course_service.go
│   ├── class_service.go
│   ├── registration_service.go
│   ├── payment_service.go
│   └── class_placement_service.go
├── handlers/           # Layer HTTP controller (JSON binding & status code)
│   ├── student_handler.go
│   ├── course_handler.go
│   ├── class_handler.go
│   ├── registration_handler.go
│   ├── payment_handler.go
│   └── class_placement_handler.go
├── routes/             # Pendaftaran router Gin & dependency injection
│   └── routes.go
├── middleware/         # HTTP Middleware (Zerolog request logger, recovery)
│   └── logger.go
├── utils/              # Helper JSON response & validation error formatter
│   └── response.go
├── Dockerfile          # Multi-stage Docker build
├── docker-compose.yml  # Orkestrasi container PostgreSQL & Go API
├── .env.example        # Template konfigurasi environment
├── go.mod / go.sum     # Go dependency management
└── main.go             # Entrypoint aplikasi
```

---

## 🔄 Alur Bisnis Utama (Core Workflow)

```text
Student + Course ──► Registration (Status: pending)
                            │
                            ▼
                    Payment (Status: pending ──► paid)
                            │ (pembayaran lunas)
                            ▼
               Class Placement (Assign ke Class)
                            │
                            ▼
           Class (Status: open ──► full jika kapasitas tercapai)
```

### 📋 Aturan Bisnis Kunci:
1. **Validasi Email Unik:** Setiap siswa harus memiliki alamat email unik di database.
2. **Pencegahan Pendaftaran Ganda:** Siswa tidak dapat mendaftar course yang sama jika masih memiliki registrasi berstatus aktif (`pending` / `registered`).
3. **Pembayaran Wajib:** Saat registrasi dibuat, tagihan `Payment` otomatis terbentuk secara atomik. Siswa **hanya bisa ditempatkan ke kelas jika status registrasinya sudah `registered`** (artinya tagihan `Payment` berstatus `paid`).
4. **Course Matching:** Siswa hanya bisa masuk ke kelas yang sesuai dengan course yang didaftarkan.
5. **Proteksi Kapasitas Kelas:** Penempatan siswa ditolak jika kapasitas kelas (`Capacity`) telah terpenuhi.

---

## 🚀 Panduan Menjalankan Aplikasi

### Opsi 1: Menggunakan Docker Compose (Direkomendasikan)
Pastikan Docker Desktop aktif di komputer Anda:

```powershell
# 1. Build dan jalankan seluruh container di background
docker compose up --build -d

# 2. Melihat log aplikasi secara realtime
docker compose logs -f app

# 3. Menghentikan container
docker compose down
```
- API Server: `http://localhost:8080`
- PostgreSQL: `localhost:5432`

### Opsi 2: Menjalankan Secara Lokal (Direct Go)
1. Salin template environment:
   ```powershell
   copy .env.example .env
   ```
2. Pastikan database PostgreSQL lokal aktif sesuai konfigurasi `.env`.
3. Jalankan aplikasi:
   ```powershell
   go run main.go
   ```

---

## 🧪 Automated Unit Testing

Project dilengkapi dengan automated unit testing untuk layer Service:

```powershell
# Menjalankan seluruh unit test
go test -v ./services/...
```

---

## 📦 Format Standar Response API

Seluruh endpoint menghasilkan format JSON response yang konsisten:

### Response Sukses (200 OK / 201 Created):
```json
{
  "success": true,
  "message": "Pesan deskriptif keberhasilan",
  "data": {}
}
```

### Response Error Validasi (422 Unprocessable Entity):
```json
{
  "success": false,
  "message": "Validasi request gagal",
  "errors": {
    "Email": "Field 'Email' harus berupa alamat email yang valid",
    "Name": "Field 'Name' wajib diisi"
  }
}
```

### Response Error Bisnis / Konflik (400 Bad Request / 409 Conflict / 404 Not Found):
```json
{
  "success": false,
  "message": "student masih memiliki pendaftaran aktif untuk course ini",
  "errors": null
}
```

---

## 📚 Dokumentasi Lengkap API & Contoh Request

Base URL: `http://localhost:8080/api/v1`

---

### 1. Health Check Endpoints

#### `GET /health` & `GET /api/v1/health`
Mengecek status ketersediaan server dan API.

**Contoh Request (cURL):**
```bash
curl -X GET http://localhost:8080/api/v1/health
```

**Contoh Response (`200 OK`):**
```json
{
  "success": true,
  "message": "API v1 is healthy 🚀",
  "data": {
    "status": "UP",
    "version": "v1"
  }
}
```

---

### 2. Modul Student (Peserta Kursus)

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/v1/students` | Mendaftarkan student baru |
| `GET` | `/api/v1/students` | Mengambil seluruh data student |
| `GET` | `/api/v1/students/:id` | Mengambil detail student |
| `PUT` | `/api/v1/students/:id` | Mengubah data student |
| `DELETE` | `/api/v1/students/:id` | Menghapus data student |
| `GET` | `/api/v1/students/:id/registrations` | Mengambil riwayat pendaftaran student |

#### 🔹 Create Student (`POST /api/v1/students`)
```bash
curl -X POST http://localhost:8080/api/v1/students \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Budi Santoso",
    "email": "budi.santoso@example.com",
    "phone": "081234567890"
  }'
```

**Response (`201 Created`):**
```json
{
  "success": true,
  "message": "Student berhasil dibuat",
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

#### 🔹 Update Student (`PUT /api/v1/students/:id`)
```bash
curl -X PUT http://localhost:8080/api/v1/students/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Budi Santoso, S.Kom",
    "email": "budi.santoso@example.com",
    "phone": "081299999999"
  }'
```

---

### 3. Modul Course (Katalog Kursus)

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/v1/courses` | Membuat kursus baru |
| `GET` | `/api/v1/courses` | Mengambil seluruh kursus |
| `GET` | `/api/v1/courses/:id` | Mengambil detail kursus beserta daftar kelasnya |
| `PUT` | `/api/v1/courses/:id` | Mengubah data kursus |
| `DELETE` | `/api/v1/courses/:id` | Menghapus kursus |
| `GET` | `/api/v1/courses/:id/registrations` | Mengambil daftar registrasi pada kursus |

#### 🔹 Create Course (`POST /api/v1/courses`)
```bash
curl -X POST http://localhost:8080/api/v1/courses \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Basic English",
    "description": "Dasar grammar dan conversation untuk pemula",
    "price": 750000,
    "duration": "3 Bulan",
    "status": "active"
  }'
```

**Response (`201 Created`):**
```json
{
  "success": true,
  "message": "Course berhasil dibuat",
  "data": {
    "id": 1,
    "name": "Basic English",
    "description": "Dasar grammar dan conversation untuk pemula",
    "price": 750000,
    "duration": "3 Bulan",
    "status": "active",
    "created_at": "2026-09-01T08:00:00Z",
    "updated_at": "2026-09-01T08:00:00Z"
  }
}
```

#### 🔹 Get Course Detail with Preloaded Classes (`GET /api/v1/courses/1`)
```bash
curl -X GET http://localhost:8080/api/v1/courses/1
```

**Response (`200 OK`):**
```json
{
  "success": true,
  "message": "Detail course berhasil diambil",
  "data": {
    "id": 1,
    "name": "Basic English",
    "description": "Dasar grammar dan conversation untuk pemula",
    "price": 750000,
    "duration": "3 Bulan",
    "status": "active",
    "classes": [
      {
        "id": 1,
        "course_id": 1,
        "name": "Basic English Pagi",
        "capacity": 15,
        "schedule": "Senin & Rabu, 09:00 - 11:00",
        "status": "open"
      }
    ]
  }
}
```

---

### 4. Modul Class (Kelas Kursus)

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/v1/classes` | Membuat kelas di bawah kursus tertentu |
| `GET` | `/api/v1/classes` | Mengambil seluruh kelas beserta Course |
| `GET` | `/api/v1/classes/:id` | Mengambil detail kelas |
| `PUT` | `/api/v1/classes/:id` | Mengubah jadwal/kapasitas kelas |
| `DELETE` | `/api/v1/classes/:id` | Menghapus kelas |
| `GET` | `/api/v1/classes/:id/students` | Mengambil daftar siswa yang terdaftar di kelas |

#### 🔹 Create Class (`POST /api/v1/classes`)
```bash
curl -X POST http://localhost:8080/api/v1/classes \
  -H "Content-Type: application/json" \
  -d '{
    "course_id": 1,
    "name": "Basic English Pagi",
    "capacity": 15,
    "schedule": "Senin & Rabu, 09:00 - 11:00",
    "status": "open"
  }'
```

**Response (`201 Created`):**
```json
{
  "success": true,
  "message": "Class berhasil dibuat",
  "data": {
    "id": 1,
    "course_id": 1,
    "name": "Basic English Pagi",
    "capacity": 15,
    "schedule": "Senin & Rabu, 09:00 - 11:00",
    "status": "open",
    "course": {
      "id": 1,
      "name": "Basic English",
      "price": 750000
    }
  }
}
```

---

### 5. Modul Course Registration (Pendaftaran)

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/v1/registrations` | Mendaftarkan siswa ke kursus (otomatis buat payment) |
| `GET` | `/api/v1/registrations` | Mengambil seluruh data registrasi |
| `GET` | `/api/v1/registrations/:id` | Mengambil detail registrasi beserta status payment |
| `PUT` | `/api/v1/registrations/:id/cancel` | Membatalkan pendaftaran |

#### 🔹 Create Registration (`POST /api/v1/registrations`)
```bash
curl -X POST http://localhost:8080/api/v1/registrations \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "course_id": 1
  }'
```

**Response (`201 Created`):**
```json
{
  "success": true,
  "message": "Pendaftaran kursus berhasil dibuat",
  "data": {
    "id": 1,
    "student_id": 1,
    "course_id": 1,
    "registration_date": "2026-09-01T08:30:00Z",
    "status": "pending",
    "student": {
      "id": 1,
      "name": "Budi Santoso"
    },
    "course": {
      "id": 1,
      "name": "Basic English",
      "price": 750000
    },
    "payment": {
      "id": 1,
      "registration_id": 1,
      "amount": 750000,
      "payment_method": "pending",
      "status": "pending"
    }
  }
}
```

---

### 6. Modul Payment (Simulasi Pembayaran)

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `GET` | `/api/v1/payments` | Mengambil daftar seluruh tagihan pembayaran |
| `GET` | `/api/v1/payments/:id` | Mengambil detail tagihan pembayaran |
| `POST` | `/api/v1/payments/:id/pay` | Memproses pembayaran (mengubah status ke `paid`) |

#### 🔹 Process Payment (`POST /api/v1/payments/:id/pay`)
```bash
curl -X POST http://localhost:8080/api/v1/payments/1/pay \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "bank_transfer",
    "amount": 750000
  }'
```

**Response (`200 OK`):**
```json
{
  "success": true,
  "message": "Pembayaran berhasil diproses! Status registrasi aktif 🎉",
  "data": {
    "id": 1,
    "registration_id": 1,
    "amount": 750000,
    "payment_method": "bank_transfer",
    "payment_date": "2026-09-01T08:35:00Z",
    "status": "paid",
    "registration": {
      "id": 1,
      "status": "registered"
    }
  }
}
```

---

### 7. Modul Class Placement (Penempatan Kelas)

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `POST` | `/api/v1/class-placements` | Menempatkan siswa lunas ke kelas tertentu |
| `GET` | `/api/v1/class-placements` | Mengambil seluruh data penempatan kelas |
| `GET` | `/api/v1/class-placements/:id` | Mengambil detail penempatan kelas |

#### 🔹 Place Student into Class (`POST /api/v1/class-placements`)
> **Aturan:** Siswa harus sudah lunas (`status: registered`), course harus cocok, dan kapasitas kelas belum penuh.

```bash
curl -X POST http://localhost:8080/api/v1/class-placements \
  -H "Content-Type: application/json" \
  -d '{
    "registration_id": 1,
    "class_id": 1
  }'
```

**Response (`201 Created`):**
```json
{
  "success": true,
  "message": "Student berhasil ditempatkan ke dalam kelas 🎉",
  "data": {
    "id": 1,
    "registration_id": 1,
    "class_id": 1,
    "placement_date": "2026-09-01T08:40:00Z",
    "registration": {
      "id": 1,
      "student": {
        "id": 1,
        "name": "Budi Santoso"
      }
    },
    "class": {
      "id": 1,
      "name": "Basic English Pagi",
      "capacity": 15,
      "schedule": "Senin & Rabu, 09:00 - 11:00",
      "status": "open"
    }
  }
}
```

#### 🔹 Check Students in Class (`GET /api/v1/classes/1/students`)
```bash
curl -X GET http://localhost:8080/api/v1/classes/1/students
```

**Response (`200 OK`):**
```json
{
  "success": true,
  "message": "Daftar siswa dalam kelas berhasil diambil",
  "data": [
    {
      "id": 1,
      "name": "Budi Santoso",
      "email": "budi.santoso@example.com",
      "phone": "081234567890"
    }
  ]
}
```

---

## 📄 Lisensi
Project ini dibuat sebagai materi pembelajaran backend Go berstandar clean layered architecture.
