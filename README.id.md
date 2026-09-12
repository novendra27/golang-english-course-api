# 🚀 English Course Registration API (Bahasa Indonesia)

RESTful API backend untuk sistem pendaftaran kursus bahasa Inggris yang dibangun menggunakan **Go (Golang)** dengan arsitektur modular berlapis (*Clean Layered Architecture*), **Gin Framework**, **GORM ORM**, **PostgreSQL Database**, **Zerolog Structured Logging**, **go-i18n Dynamic Localization**, dan **Docker Containerization**.

---

## 📑 Daftar Isi

- [Tech Stack & Technologies](#-tech-stack--technologies)
- [Arsitektur & Struktur Folder](#-arsitektur--struktur-folder)
- [Alur Bisnis Utama (Core Workflow)](#-alur-bisnis-utama-core-workflow)
- [Panduan Menjalankan Aplikasi](#-panduan-menjalankan-aplikasi)
- [Automated Unit Testing](#-automated-unit-testing)
- [Dukungan Multi-Bahasa (i18n)](#-dukungan-multi-bahasa-i18n)
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
| **Localization (i18n)**| [go-i18n](https://github.com/nicksnyder/go-i18n) | Dynamic multi-language localization (EN / ID). |
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
├── middleware/         # HTTP Middleware (Zerolog request logger, i18n detection)
│   ├── logger.go
│   └── i18n.go
├── locales/            # Embedded kamus bahasa JSON (en.json, id.json)
│   ├── locales.go
│   ├── en.json
│   └── id.json
├── utils/              # Helper JSON response & i18n translation engine
│   ├── i18n.go
│   ├── i18n_test.go
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

## 🌐 Dukungan Multi-Bahasa (i18n)

Sistem mendukung respon multi-bahasa secara dinamis:
- **Default:** English (`en`)
- **Bahasa Indonesia:** `id`

Cara memilih bahasa pada request:
1. **Query Parameter:** `GET /api/v1/courses?lang=id`
2. **Header HTTP:** `Accept-Language: id-ID,id;q=0.9,en;q=0.8`

---

## 🚀 Panduan Menjalankan Aplikasi

### Opsi 1: Menggunakan Docker Compose (Direkomendasikan)
```powershell
docker compose up --build -d
docker compose logs -f app
docker compose down
```
- API Server: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/swagger/index.html`

### Opsi 2: Menjalankan Secara Lokal (Direct Go)
1. Salin template environment:
   ```powershell
   copy .env.example .env
   ```
2. Jalankan aplikasi:
   ```powershell
   go run main.go
   ```

---

## ⚡ 5-Minute Quick Tour (Uji Coba Alur API)

Coba seluruh siklus pendaftaran siswa, pembayaran tagihan, dan penempatan kelas dalam 5 menit menggunakan perintah `curl` berurutan berikut:

### Langkah 1: Daftarkan Siswa Baru
```bash
curl -X POST http://localhost:8080/api/v1/students \
  -H "Content-Type: application/json" \
  -d '{"name": "Alex Johnson", "email": "alex@example.com", "phone": "081234567890"}'
```
> **Hasil:** Menghasilkan data Student baru (ID: `1`).

### Langkah 2: Buat Kursus & Buka Kelas
```bash
# 2a. Buat Kursus
curl -X POST http://localhost:8080/api/v1/courses \
  -H "Content-Type: application/json" \
  -d '{"name": "IELTS Masterclass", "description": "Intensive IELTS band 7.5+ prep", "price": 1250000, "duration": "2 Bulan", "status": "active"}'

# 2b. Buka Kelas (Kapasitas: 15)
curl -X POST http://localhost:8080/api/v1/classes \
  -H "Content-Type: application/json" \
  -d '{"course_id": 1, "name": "IELTS Weekend Intensive", "capacity": 15, "schedule": "Sabtu & Minggu, 10:00 - 13:00", "status": "open"}'
```
> **Hasil:** Menghasilkan Course (ID: `1`) dan Class (ID: `1`).

### Langkah 3: Daftarkan Siswa ke Kursus (Otomatis Buat Tagihan)
```bash
curl -X POST http://localhost:8080/api/v1/registrations \
  -H "Content-Type: application/json" \
  -d '{"student_id": 1, "course_id": 1}'
```
> **State Machine:** Status pendaftaran `pending` dengan tagihan Payment ID `1` (Nominal: `1250000`).

### Langkah 4: Lunaskan Tagihan Pembayaran
```bash
curl -X POST http://localhost:8080/api/v1/payments/1/pay \
  -H "Content-Type: application/json" \
  -d '{"payment_method": "bank_transfer", "amount": 1250000}'
```
> **State Machine:** Status Payment berubah menjadi `paid` $\rightarrow$ Status pendaftaran otomatis aktif (`registered`) 🎉.

### Langkah 5: Tempatkan Siswa ke Dalam Kelas
```bash
curl -X POST http://localhost:8080/api/v1/class-placements \
  -H "Content-Type: application/json" \
  -d '{"registration_id": 1, "class_id": 1}'
```
> **Validasi Aturan Bisnis:** Sistem memverifikasi kesesuaian kursus dan kapasitas kelas sebelum penempatan.

### Langkah 6: Uji Coba Multi-Bahasa Dinamis (i18n)
```bash
# 6a. Respon default (Bahasa Inggris)
curl -X GET http://localhost:8080/api/v1/courses/1

# 6b. Respon Bahasa Indonesia via query parameter
curl -X GET "http://localhost:8080/api/v1/courses/1?lang=id"
```

---

## 🧪 Automated Unit Testing

```powershell
go test -v ./...
```

