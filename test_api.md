# 📖 Panduan Pengujian & Referensi API Customs Clearance Service

Dokumen ini berisi daftar lengkap endpoint yang tersedia pada aplikasi **Customs Clearance API**, format payload request/response, serta langkah demi langkah skenario pengujian menggunakan **cURL** (dapat diimpor langsung ke Postman).

---

## ⚙️ Variabel Global & Lingkungan Lokal
* **Base URL**: `http://localhost:8082/api/v1`
* **Health Check URL**: `http://localhost:8082/health`
* **Format Response Sukses**:
  ```json
  {
    "success": true,
    "message": "Pesan deskriptif sukses",
    "data": { ... }
  }
  ```
* **Format Response Gagal**:
  ```json
  {
    "success": false,
    "message": "Pesan kegagalan",
    "errors": "Rincian detail error / validasi"
  }
  ```

---

## 🔒 1. Endpoint Autentikasi & Profil (Public & Protected)

### A. Registrasi Akun Baru
Membuat akun importir umum (`User`) atau petugas bea cukai (`Officer`).
* **Method**: `POST`
* **URL**: `/api/v1/register`
* **Payload JSON**:
  ```json
  {
    "name": "Supardi Officer",
    "email": "supardi@customs.go.id",
    "password": "password123",
    "role": "Officer"
  }
  ```
  *(Role dapat diisi `"Officer"` atau `"User"`. Jika kosong defaultnya `"User"`)*
* **Contoh cURL**:
  ```bash
  curl -X POST http://localhost:8082/api/v1/register \
    -H "Content-Type: application/json" \
    -d '{"name": "Supardi Officer", "email": "supardi@customs.go.id", "password": "password123", "role": "Officer"}'
  ```

### B. Login Akun
Melakukan login untuk mendapatkan token akses JWT.
* **Method**: `POST`
* **URL**: `/api/v1/login`
* **Payload JSON**:
  ```json
  {
    "email": "supardi@customs.go.id",
    "password": "password123"
  }
  ```
* **Contoh cURL**:
  ```bash
  curl -X POST http://localhost:8082/api/v1/login \
    -H "Content-Type: application/json" \
    -d '{"email": "supardi@customs.go.id", "password": "password123"}'
  ```

### C. Get Profil Saya (Me)
Mengambil data profil user yang sedang login menggunakan token JWT.
* **Method**: `GET`
* **URL**: `/api/v1/me`
* **Headers**: `Authorization: Bearer <JWT_TOKEN>`
* **Contoh cURL**:
  ```bash
  curl -X GET http://localhost:8082/api/v1/me \
    -H "Authorization: Bearer <JWT_TOKEN>"
  ```

---

## 📊 2. Endpoint Seeder & Dashboard (Public & Protected)

### A. Data Seeder (Public untuk Development)
Membersihkan seluruh database (`TRUNCATE CASCADE`) dan mengisi data tiruan Bea Cukai dalam jumlah besar untuk keperluan pengetesan visual.
* **Method**: `POST`
* **URL**: `/api/v1/seed`
* **Contoh cURL**:
  ```bash
  curl -X POST http://localhost:8082/api/v1/seed
  ```

### B. Ringkasan Dashboard (Protected - Officer Only)
Menampilkan statistik pengajuan barang, kategori risiko, top komoditas, top pelabuhan, dan 5 transaksi teranyar.
* **Method**: `GET`
* **URL**: `/api/v1/dashboard`
* **Headers**: `Authorization: Bearer <JWT_TOKEN_OFFICER>`
* **Contoh cURL**:
  ```bash
  curl -X GET http://localhost:8082/api/v1/dashboard \
    -H "Authorization: Bearer <JWT_TOKEN_OFFICER>"
  ```

---

## 📦 3. Endpoint Master Data CRUD (Protected - Officer Only)
*Semua endpoint di bawah ini membutuhkan Header `Authorization: Bearer <JWT_TOKEN_OFFICER>`.*

### A. Negara (Country)
* **GET Semua Negara**: `GET /api/v1/master/countries`
  ```bash
  curl -X GET http://localhost:8082/api/v1/master/countries -H "Authorization: Bearer <TOKEN>"
  ```
* **POST Buat Negara Baru**: `POST /api/v1/master/countries`
  * Payload: `{"code": "SGP", "name": "Singapore"}`
  ```bash
  curl -X POST http://localhost:8082/api/v1/master/countries \
    -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" \
    -d '{"code": "SGP", "name": "Singapore"}'
  ```
* **GET Detail Negara**: `GET /api/v1/master/countries/:id`
* **PUT Update Negara**: `PUT /api/v1/master/countries/:id`
  * Payload: `{"code": "SGP", "name": "Singapore Baru"}`
* **DELETE Hapus Negara**: `DELETE /api/v1/master/countries/:id`

### B. Pelabuhan (Port)
* **GET Semua Pelabuhan**: `GET /api/v1/master/ports`
* **POST Buat Pelabuhan Baru**: `POST /api/v1/master/ports`
  * Payload: `{"code": "SGPIN", "name": "Jurong Port", "country_id": 1}`
  ```bash
  curl -X POST http://localhost:8082/api/v1/master/ports \
    -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" \
    -d '{"code": "SGPIN", "name": "Jurong Port", "country_id": 1}'
  ```
* **PUT Update Pelabuhan**: `PUT /api/v1/master/ports/:id`
* **DELETE Hapus Pelabuhan**: `DELETE /api/v1/master/ports/:id`

### C. Komoditas (Commodity)
* **GET Semua Komoditas**: `GET /api/v1/master/commodities`
* **POST Buat Komoditas Baru**: `POST /api/v1/master/commodities`
  * Payload: `{"hs_code": "85171200", "description": "Handphone", "import_duty_rate": 10.0, "vat_rate": 11.0}`
  ```bash
  curl -X POST http://localhost:8082/api/v1/master/commodities \
    -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" \
    -d '{"hs_code": "85171200", "description": "Handphone", "import_duty_rate": 10.0, "vat_rate": 11.0}'
  ```
* **PUT Update Komoditas**: `PUT /api/v1/master/commodities/:id`
* **DELETE Hapus Komoditas**: `DELETE /api/v1/master/commodities/:id`

---

## 🔄 4. Endpoint Workflow & Transisi Status (Protected - Officer Only)
*Semua endpoint di bawah ini membutuhkan Header `Authorization: Bearer <JWT_TOKEN_OFFICER>`.*

### A. Inisialisasi Workflow & Cek Risiko
Memicu evaluasi Risk Engine pada clearance yang baru masuk (`SUBMITTED`).
* **Method**: `POST`
* **URL**: `/api/v1/workflow/init`
* **Payload JSON**:
  ```json
  {
    "clearance_id": 1
  }
  ```
* **Contoh cURL**:
  ```bash
  curl -X POST http://localhost:8082/api/v1/workflow/init \
    -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" \
    -d '{"clearance_id": 1}'
  ```

### B. Input Hasil Pemeriksaan Fisik (Khusus Jalur Merah/HIGH Risk)
Menentukan kelolosan pemeriksaan fisik barang.
* **Method**: `POST`
* **URL**: `/api/v1/workflow/inspection`
* **Payload JSON**:
  ```json
  {
    "clearance_id": 1,
    "result": "PASS" 
  }
  ```
  *(Result wajib bernilai `"PASS"` untuk lolos, atau `"FAIL"` untuk mengunci ke status `"HOLD"`)*
* **Contoh cURL**:
  ```bash
  curl -X POST http://localhost:8082/api/v1/workflow/inspection \
    -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" \
    -d '{"clearance_id": 1, "result": "PASS"}'
  ```

### C. Persetujuan Dokumen (Approve)
Memberikan persetujuan Bea Cukai pada dokumen clearance.
* **Method**: `POST`
* **URL**: `/api/v1/workflow/approve`
* **Payload JSON**:
  ```json
  {
    "clearance_id": 1
  }
  ```
* **Contoh cURL**:
  ```bash
  curl -X POST http://localhost:8082/api/v1/workflow/approve \
    -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" \
    -d '{"clearance_id": 1}'
  ```

### D. Penerbitan SPPB (Release)
Menerbitkan dokumen pengeluaran barang resmi (SPPB).
* **Method**: `POST`
* **URL**: `/api/v1/workflow/release`
* **Payload JSON**:
  ```json
  {
    "clearance_id": 1
  }
  ```
* **Contoh cURL**:
  ```bash
  curl -X POST http://localhost:8082/api/v1/workflow/release \
    -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" \
    -d '{"clearance_id": 1}'
  ```

### E. Pengeluaran Pabean (Gate Out)
Mencatat keluarnya barang secara fisik dari gerbang pelabuhan/pabean.
* **Method**: `POST`
* **URL**: `/api/v1/workflow/gate-out`
* **Payload JSON**:
  ```json
  {
    "clearance_id": 1
  }
  ```
* **Contoh cURL**:
  ```bash
  curl -X POST http://localhost:8082/api/v1/workflow/gate-out \
    -H "Authorization: Bearer <TOKEN>" -H "Content-Type: application/json" \
    -d '{"clearance_id": 1}'
  ```

---

## 🧪 Skenario Pengetesan End-to-End

### Skenario Jalur Merah (Barang Nilai Tinggi > Rp 50 Juta)
1. **Jalankan Seeder**: Panggil `/api/v1/seed` untuk memastikan database bersih & terisi.
2. **Login Officer**: Login dengan email `supardi@customs.go.id` dan password `password123`. Salin token JWT yang dihasilkan ke variable Postman/cURL.
3. **Simulasi Pengajuan**: Masukkan data clearance dengan nilai tinggi (seeder sudah otomatis membuat data ini dengan status `SUBMITTED`, misalnya pada `id: 35`).
4. **Mulai Workflow**: Panggil `/workflow/init` untuk `clearance_id: 35`. Respon akan menunjukkan status berubah menjadi `"INSPECTION"` (Jalur Merah).
5. **Cek Validasi Blokir**: Coba panggil `/workflow/approve` untuk `id: 35`. Server akan menolak dengan error `400 Bad Request` karena pemeriksaan fisik belum dilakukan.
6. **Input Hasil Periksa Fisik**: Panggil `/workflow/inspection` dengan payload `{"clearance_id": 35, "result": "PASS"}`. Status berubah menjadi `"INSPECTION_PASSED"`.
7. **Beri Persetujuan**: Panggil `/workflow/approve` dengan `clearance_id: 35`. Status berubah menjadi `"APPROVED"`.
8. **Rilis SPPB**: Panggil `/workflow/release` dengan `clearance_id: 35`. Status berubah menjadi `"RELEASED"`.
9. **Gerbang Keluar**: Panggil `/workflow/gate-out` dengan `clearance_id: 35`. Status berubah menjadi `"GATE_OUT"`.

### Skenario Jalur Hijau (Bypass Pemeriksaan Fisik)
1. **Pilih Clearance Jalur Hijau**: Cari clearance dengan nilai di bawah Rp 50 juta (misalnya `id: 1` atau `id` lain yang terdeteksi `LOW` risk setelah di `/workflow/init`).
2. **Mulai Workflow**: Panggil `/workflow/init` untuk `clearance_id: 2` (LOW risk). Status akan tetap `"SUBMITTED"`.
3. **Bypass Inspeksi**: Lewati `/workflow/inspection` dan langsung panggil `/workflow/approve` dengan payload `{"clearance_id": 2}`. Server akan langsung menyetujui dokumen dan mengubah status menjadi `"APPROVED"`.
4. **Rilis SPPB & Gate Out**: Lanjutkan memanggil `/workflow/release` dan `/workflow/gate-out` untuk menyelesaikan proses.
