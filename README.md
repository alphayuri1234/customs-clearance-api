# 🚢 Customs Clearance API Service

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Docker Compose](https://img.shields.io/badge/Docker_Compose-Supported-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![Swagger API](https://img.shields.io/badge/Swagger-Interactive_Docs-85EA2D?style=flat&logo=swagger)](http://localhost:8082/swagger/index.html)

Sistem Simulasi Customs Clearance Bea Cukai (Jalur Merah & Jalur Hijau) berbasis RESTful API menggunakan **Go**, **Gin Gonic**, **GORM (PostgreSQL)**, dan **Swagger Docs**. 

Aplikasi ini mendemonstrasikan otomatisasi penentuan jalur risiko (*Risk Engine*) kepabeanan berdasarkan asal negara, jenis komoditas, tarif bea masuk, serta nilai valuasi barang.

---

## 🎯 Fitur Utama

1. **Risk Engine & Alur Kerja Otomatis**:
   - **Jalur Hijau (Low Risk)**: Melewati pemeriksaan fisik, langsung menuju tahap persetujuan dokumen (*Approved*).
   - **Jalur Merah (High Risk)**: Wajib melalui pemeriksaan fisik oleh petugas (*Inspection*), jika lolos (*Inspection Passed*) baru bisa disetujui (*Approved*).
2. **Autentikasi & Otorisasi**:
   - Berbasis **JWT (JSON Web Token)**.
   - Perbedaan peran pengguna (*Role-based Access Control*): **Importer (User)** & **Officer (Petugas Bea Cukai)**.
3. **Master Data Management**: CRUD lengkap untuk Negara (*Countries*), Pelabuhan (*Ports*), dan Komoditas (*Commodities*).
4. **Data Seeder Otomatis**: Mengisi database dengan puluhan data simulasi secara terstruktur untuk keperluan dashboard analitis.
5. **Dashboard Statistik (Analitik Officer)**:
   - Total Transaksi, Total Valuasi, Persentase Jalur Merah vs Hijau.
   - Top Pelabuhan Teraktif & Top Komoditas Terimpor.
   - Riwayat transaksi teranyar.
6. **Dokumentasi Swagger Terintegrasi**: Akses interaktif ke seluruh endpoint API.

---

## 🛠️ Arsitektur & Teknologi

- **Backend Language**: Go (Golang) v1.26.4
- **Web Framework**: Gin Gonic v1.12.0
- **Database ORM**: GORM v1.25.12 (PostgreSQL Driver)
- **Database**: PostgreSQL 16 (Alpine-based Container)
- **API Documentation**: Swaggo / Gin-Swagger
- **Containerization**: Docker & Docker Compose

---

## 📁 Struktur Folder Proyek

```text
customs-clearance-api/
├── cmd/
│   └── main.go                  # Entry point aplikasi
├── config/
│   └── jwt.go                   # Konfigurasi secret JWT
├── controllers/                 # Handler REST API (Gin)
│   ├── auth_controller.go
│   ├── dashboard_controller.go
│   ├── master_controller.go
│   └── workflow_controller.go
├── database/
│   └── postgres.go              # Koneksi PostgreSQL & AutoMigrate GORM
├── docs/                        # File Swagger Docs auto-generated
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── middleware/                  # JWT & Role Authentication Middleware
│   └── auth_middleware.go
├── models/                      # Defisini Struct GORM / Skema Database
│   ├── clearance.go
│   ├── user.go
│   ├── officer.go
│   ├── master.go
│   ├── risk_profile.go
│   ├── inspection_result.go
│   └── release_order.go
├── routes/                      # Pengelompokan & Registrasi Endpoint
│   ├── router.go
│   ├── auth_routes.go
│   ├── dashboard_routes.go
│   ├── master_routes.go
│   └── workflow_routes.go
├── services/                    # Logika Bisnis & Validasi (Core Logic)
│   ├── auth_service.go
│   ├── dashboard_service.go
│   ├── master_service.go
│   ├── seeder_service.go
│   └── workflow_service.go
├── .env                         # Konfigurasi environment variables lokal
├── Dockerfile                   # Docker build multi-stage untuk Golang
├── docker-compose.yml           # Orkestrasi Docker untuk PostgreSQL & Go API
└── test_api.md                  # Panduan lengkap pengetesan API (cURL)
```

---

## 🚀 Panduan Menjalankan Aplikasi

Aplikasi ini dapat dijalankan dengan sangat mudah menggunakan **Docker Compose** (Semua komponen otomatis siap digunakan) atau secara manual.

### Metode A: Menggunakan Docker Compose (Sangat Direkomendasikan)

Metode ini otomatis mengunduh PostgreSQL, membuat database, melakukan build binary Go secara optimal (multi-stage), serta menjalankan server API dalam satu jaringan virtual Docker yang aman.

1. **Pastikan Docker & Docker Compose sudah terpasang** di mesin Anda.
2. **Jalankan perintah berikut** di direktori utama proyek:
   ```bash
   docker compose up --build -d
   ```
3. **Cek status container** untuk memastikan semuanya berjalan lancar:
   ```bash
   docker compose ps
   ```
4. **Hentikan Container** jika sudah selesai:
   ```bash
   docker compose down
   ```

Setelah sukses dijalankan, API akan aktif pada port **`8082`** di host Anda.

---

### Metode B: Menjalankan secara Manual (Local Go & Docker PostgreSQL)

Jika Anda ingin menjalankan atau mengembangkan kode secara lokal di host Anda:

1. **Jalankan Database PostgreSQL saja** menggunakan Docker Compose:
   ```bash
   docker compose up postgres -d
   ```
   *(DB akan terekspos di host Anda pada port `5435`)*.

2. **Periksa konfigurasi `.env`** di direktori utama, pastikan isinya:
   ```env
   PORT=8082
   DB_HOST=localhost
   DB_PORT=5435
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=customs_clearance
   DB_SSLMODE=disable
   JWT_SECRET=supersecretkey
   ```

3. **Unduh Dependensi Go**:
   ```bash
   go mod download
   ```

4. **Jalankan Aplikasi**:
   ```bash
   go run cmd/main.go
   ```

Server API akan berjalan pada alamat: **`http://localhost:8082`**.

---

## 📌 Langkah Awal Penggunaan API

Setelah server API aktif, ikuti langkah-langkah di bawah ini untuk menginisialisasi data dan memulai pengujian:

### 1. Jalankan Seeder Data (Wajib)
Panggil endpoint Seeder untuk membersihkan dan mengisi database dengan data dummy terstruktur (Pengguna, Petugas, Negara, Pelabuhan, Komoditas, dan 35 transaksi clearance):
* **Method**: `POST`
* **URL**: `http://localhost:8082/api/v1/seed`
* **cURL**:
  ```bash
  curl -X POST http://localhost:8082/api/v1/seed
  ```

### 2. Login ke Sistem (Sebagai Officer/Petugas)
Gunakan salah satu akun Petugas Bea Cukai hasil seeder untuk login:
* **Method**: `POST`
* **URL**: `http://localhost:8082/api/v1/login`
* **Body (JSON)**:
  ```json
  {
    "email": "supardi@customs.go.id",
    "password": "password123"
  }
  ```
* **cURL**:
  ```bash
  curl -X POST http://localhost:8082/api/v1/login \
    -H "Content-Type: application/json" \
    -d '{"email": "supardi@customs.go.id", "password": "password123"}'
  ```
> **Catatan**: Salin token JWT dari response `data.token` untuk digunakan sebagai header Authorization pada request lainnya.

### 3. Akses Dokumentasi Swagger
Buka browser dan akses halaman berikut untuk melihat spesifikasi detail seluruh endpoint dan mencobanya secara interaktif:
👉 **[http://localhost:8082/swagger/index.html](http://localhost:8082/swagger/index.html)**

---

## 🧪 Skenario Pengujian Alur Kerja (Workflow)

Untuk panduan lengkap skenario pengujian fungsional dari hulu ke hilir (End-to-End), silakan merujuk ke berkas dokumentasi pengujian terpisah:
👉 **[Panduan Pengujian & Referensi API (test_api.md)](test_api.md)**

---

## 👨‍💻 Kontributor & Lisensi
- **Sistem**: Customs Clearance Bea Cukai Simulator (Training Go Lang)
- **Lisensi**: Apache License 2.0
