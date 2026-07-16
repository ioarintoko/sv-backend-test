# SV Backend Test - Article API

Backend service untuk use case **Post Article**, dibuat sebagai bagian dari Test Backend Sharing Vision 2023.

## Tech Stack

- **Go** 1.22+ (Gin framework)
- **MySQL** 8.0
- **golang-migrate** untuk database migration
- **go-playground/validator** (built-in di Gin binding) untuk validasi request

## Struktur Project

```
.
├── cmd/api/main.go              # entry point, routing
├── internal/
│   ├── config/                  # load environment variables
│   ├── database/                # koneksi MySQL
│   ├── models/                  # struct Post & PostRequest (dengan validasi)
│   └── handlers/                # logic 5 endpoint
├── migrations/                  # file migration (up/down)
├── postman/                     # Postman collection siap import
└── .env.example
```

## Setup & Menjalankan Project

### 1. Prasyarat
- Go 1.22+
- MySQL Server sudah berjalan
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) sudah terinstall

### 2. Clone & install dependencies
```bash
git clone <repo-url>
cd sv-backend-test
go mod download
```

### 3. Setup database
Buat database secara manual (tabel dibuat lewat migration, bukan manual):
```sql
CREATE DATABASE article;
```

### 4. Setup environment variables
Copy `.env.example` menjadi `.env`, lalu sesuaikan kredensial MySQL:
```bash
cp .env.example .env
```

### 5. Jalankan migration
```bash
migrate -database "mysql://<user>:<password>@tcp(127.0.0.1:3306)/article" -path migrations up
```

### 6. Jalankan server
```bash
go run cmd/api/main.go
```
Server berjalan di `http://localhost:8080` (port bisa diubah lewat `.env`).

## API Endpoints

| No | Method | URL | Deskripsi |
|----|--------|-----|-----------|
| 1 | POST | `/article/` | Membuat article baru |
| 2 | GET | `/article/<limit>/<offset>` | List article dengan paging |
| 3 | GET | `/article/<id>` | Detail article by id |
| 4 | PUT / PATCH | `/article/<id>` | Update article by id |
| 5 | DELETE | `/article/<id>` | Hapus article by id |

### Validasi Request (Create & Update)
| Field | Aturan |
|---|---|
| title | required, minimal 20 karakter |
| content | required, minimal 200 karakter |
| category | required, minimal 3 karakter |
| status | required, harus salah satu dari: `publish`, `draft`, `thrash` |

### Contoh Request - Create Article
```json
POST /article/
{
  "title": "Judul artikel percobaan pertama yang cukup panjang",
  "content": "Konten artikel minimal 200 karakter...",
  "category": "Teknologi",
  "status": "publish"
}
```

## Postman Collection

File Postman collection tersedia di [`postman/SV-Backend-Test-Article.postman_collection.json`](./postman/SV-Backend-Test-Article.postman_collection.json). Berisi request untuk seluruh endpoint di atas, termasuk 1 contoh negative test case untuk validasi.

Cara import: Postman → Import → pilih file tersebut. Variable `base_url` sudah diset ke `http://localhost:8080`, tinggal diganti kalau deploy ke environment lain.

## Catatan Desain Teknis

### Method untuk Update & Delete
Soal menyebutkan endpoint update bisa menggunakan **POST atau PUT atau PATCH**, dan endpoint delete bisa menggunakan **POST atau DELETE**. Implementasi ini secara sadar memilih **PUT/PATCH** untuk update dan **DELETE** untuk hapus (bukan POST), dengan pertimbangan:

1. **Konflik routing** — Gin (dan router berbasis radix tree pada umumnya) tidak mengizinkan path `POST /article/` dan `POST /article/:id` didaftarkan bersamaan karena keduanya bersinggungan secara pattern di level yang sama.
2. **Kesesuaian dengan konvensi REST** — penggunaan method HTTP yang semantically correct (PUT/PATCH untuk update, DELETE untuk hapus) lebih idiomatik dan lebih mudah dipahami/di-maintain dibanding menumpuk banyak aksi berbeda di satu method POST.

Soal secara eksplisit memberi pilihan ("atau"), sehingga keputusan ini valid sesuai spesifikasi sambil tetap menghindari trade-off teknis yang tidak perlu.

### Path Parameter untuk Get List vs Get by ID
Endpoint `GET /article/<limit>/<offset>` dan `GET /article/<id>` juga sempat mengalami konflik serupa di Gin (wildcard dengan nama parameter berbeda di posisi segment yang sama tidak diizinkan). Solusinya, nama parameter internal disamakan (`:id`) untuk kedua route, meski secara nilai tetap dipakai sesuai konteksnya (`limit` untuk route dua-segment). Ini murni penamaan internal dan tidak mengubah kontrak API yang diminta soal.
