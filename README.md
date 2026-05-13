# Cinema API — Mitra Kasih Perkasa Backend Test

**Nama Lengkap:** [Mohamad Afid Yoga Pratama Putra]  
**Email:** [afidyoga45dr@gmail.com]  
**Phone:** [081233721538]

---

## Isi Repository

| Folder / File | Isi |
|---|---|
| `cmd/main.go` | Entry point aplikasi |
| `internal/` | Config, handler, middleware, model, repository, service |
| `migrations/001_init.sql` | Schema PostgreSQL lengkap |
| `migrations/002_seed.sql` | Data awal untuk testing |
| `docs/SYSTEM_DESIGN_EXPLANATION.md` | Penjelasan lengkap Point A |
| `docs/Cinema_API.postman_collection.json` | Export Postman siap import |
| `.env.example` | Template environment variable |

---

## Stack Teknologi

- **Language:** Go 1.26
- **Framework:** Gin
- **Database:** PostgreSQL 17.2
- **Cache / Seat Lock:** Redis (TTL-based distributed lock)
- **Auth:** JWT HS256 + bcrypt
- **Payment:** Midtrans / Xendit (webhook-based)

---

## Setup & Menjalankan

### 1. Clone dan install dependency

```bash
git clone https://github.com/afidyoga/cinema-api
cd cinema-api
go mod tidy
```

### 2. Buat database dan jalankan migrasi

```bash
createdb cinema_db
psql -U postgres -d cinema_db -f migrations/001_init.sql
psql -U postgres -d cinema_db -f migrations/002_seed.sql
```

### 3. Konfigurasi environment

```bash
cp .env.example .env
```

Edit `.env` sesuai konfigurasi lokal:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=cinema_db
JWT_SECRET=ganti_dengan_secret_panjang_acak
APP_PORT=8080
```

### 4. Jalankan server

```bash
go run cmd/main.go
```

Server berjalan di `http://localhost:8080`. Cek health: `GET /health`

---

## Endpoint API

### Auth

| Method | Endpoint | Akses | Deskripsi |
|--------|----------|-------|-----------|
| `POST` | `/api/v1/auth/register` | Public | Daftar akun baru |
| `POST` | `/api/v1/auth/login` | Public | Login, mendapat JWT token |
| `GET` | `/api/v1/auth/me` | Bearer Token | Info user yang sedang login |

### Jadwal Tayang

| Method | Endpoint | Akses | Deskripsi |
|--------|----------|-------|-----------|
| `GET` | `/api/v1/schedules` | Bearer Token | Daftar semua jadwal (pagination) |
| `GET` | `/api/v1/schedules/:id` | Bearer Token | Detail satu jadwal |
| `POST` | `/api/v1/schedules` | Admin only | Buat jadwal baru |
| `PUT` | `/api/v1/schedules/:id` | Admin only | Update jadwal |
| `DELETE` | `/api/v1/schedules/:id` | Admin only | Hapus jadwal |

Query params GET list: `?page=1&limit=20`

---

## Akun Seed Data untuk Testing

| Email | Password | Role |
|-------|----------|------|
| `admin@mkp.com` | `password` | admin |

Untuk buat akun customer baru gunakan endpoint `/api/v1/auth/register`.

---

## Import Postman

1. Buka Postman → **Import** → pilih `docs/Cinema_API.postman_collection.json`
2. Set variabel `base_url` = `http://localhost:8080/api/v1`
3. Jalankan **Login** dulu — token otomatis tersimpan ke `{{token}}`
4. Semua endpoint lain langsung bisa dicoba

---

## Penjelasan Desain Sistem (Point A)

Penjelasan lengkap seat locking, restok tiket, dan alur refund tersedia di:
`docs/SYSTEM_DESIGN_EXPLANATION.md`

---

## Struktur Project

```
cinema-api/
├── cmd/main.go
├── internal/
│   ├── config/config.go
│   ├── handler/auth_handler.go
│   ├── handler/schedule_handler.go
│   ├── middleware/auth.go
│   ├── model/model.go
│   ├── repository/repository.go
│   ├── service/auth_service.go
│   └── service/schedule_service.go
├── migrations/
│   ├── 001_init.sql
│   └── 002_seed.sql
├── docs/
│   ├── SYSTEM_DESIGN_EXPLANATION.md
│   └── Cinema_API.postman_collection.json
├── .env.example
├── .gitignore
└── go.mod
```
