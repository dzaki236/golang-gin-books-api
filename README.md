# Books API

REST API untuk mengelola pengguna, kategori, dan buku. API dibuat dengan Go dan Gin, menggunakan PostgreSQL, serta melindungi endpoint data dengan JWT Bearer token.

## Teknologi

- Go (versi sesuai `go.mod`)
- Gin
- PostgreSQL
- JWT (HS256)
- SQL migration dengan `sql-migrate`

## Prasyarat

- Go versi yang kompatibel dengan versi pada `go.mod`
- PostgreSQL yang dapat diakses dari mesin aplikasi

## Konfigurasi

Buat berkas `.env` di root project:

```env
DATABASE_URL=postgres://USER:PASSWORD@localhost:5432/books_api?sslmode=disable
JWT_SECRET=ganti_dengan_secret_acak_minimal_32_karakter
PORT=8080
```

`DATABASE_URL` wajib berisi URL koneksi PostgreSQL. `JWT_SECRET` wajib memiliki setidaknya 32 karakter. `PORT` opsional dan bernilai `8080` jika tidak diatur. Jangan masukkan secret atau kredensial asli ke repository.

## Menyiapkan Database dan Menjalankan API

Dari root repository, jalankan migrasi untuk membuat tabel:

```sh
go run ./migrations/execute --path=./migrations/database
```

Setelah migrasi berhasil, jalankan API:

```sh
go run ./cmd/api
```

Server akan tersedia di `http://localhost:8080` (atau port dari `PORT`). Migrasi dapat dijalankan terhadap satu file `.sql` atau direktori yang berisi file migrasi.

## Format Respons

Respons API menggunakan envelope berikut:

```json
{
  "status": "success",
  "message": "daftar buku berhasil diambil",
  "data": []
}
```

Nilai `data` dapat berupa objek, array, atau `null`, bergantung endpoint. Respons error menggunakan format yang sama dengan pesan kesalahan dan status HTTP yang sesuai.

## Autentikasi

Endpoint register dan login tidak memerlukan token. Endpoint logout, kategori, dan buku memerlukan header:

```http
Authorization: Bearer <access_token>
```

Token yang diterbitkan saat login berlaku selama 24 jam. Panggil logout dengan token tersebut untuk mencabutnya sebelum kedaluwarsa.

## Daftar Endpoint

Semua path di bawah memiliki prefix `/api`.

| Method | Path | Auth | Kegunaan |
| --- | --- | --- | --- |
| `POST` | `/api/auth/register` | Tidak | Membuat akun pengguna |
| `POST` | `/api/auth/login` | Tidak | Login dan memperoleh access token |
| `POST` | `/api/auth/logout` | Ya | Logout dan mencabut token aktif |
| `GET` | `/api/categories` | Ya | Mengambil semua kategori |
| `POST` | `/api/categories` | Ya | Membuat kategori |
| `GET` | `/api/categories/:id` | Ya | Mengambil detail kategori |
| `PUT` | `/api/categories/:id` | Ya | Mengubah nama kategori |
| `DELETE` | `/api/categories/:id` | Ya | Menghapus kategori |
| `GET` | `/api/categories/:id/books` | Ya | Mengambil kategori beserta buku-bukunya |
| `GET` | `/api/books` | Ya | Mengambil semua buku |
| `POST` | `/api/books` | Ya | Membuat buku |
| `GET` | `/api/books/:id` | Ya | Mengambil detail buku |
| `PUT` | `/api/books/:id` | Ya | Mengubah data buku |
| `DELETE` | `/api/books/:id` | Ya | Menghapus buku |

## Contoh Penggunaan

Contoh berikut menggunakan `curl` dan dijalankan dari terminal. Ganti nilai placeholder sesuai kebutuhan.

### Registrasi

Username harus terdiri dari 3-100 karakter dan password 8-72 karakter.

```sh
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"reader01","password":"password-aman"}'
```

### Login

Simpan nilai `data.access_token` dari respons untuk permintaan yang memerlukan autentikasi.

```sh
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"reader01","password":"password-aman"}'
```

### Membuat Kategori

Nama kategori wajib diisi dan maksimal 150 karakter.

```sh
curl -X POST http://localhost:8080/api/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{"name":"Pemrograman"}'
```

### Membuat Buku

Semua field berikut wajib. `category_id` harus merujuk kategori yang sudah ada. Tahun terbit yang diterima saat ini adalah 1980-2024.

```sh
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "title":"Belajar Go",
    "description":"Pengantar bahasa Go",
    "image_url":"https://example.com/belajar-go.jpg",
    "release_year":2024,
    "price":95000,
    "total_page":240,
    "category_id":1
  }'
```

### Mengambil Buku dan Buku dalam Kategori

```sh
curl http://localhost:8080/api/books \
  -H "Authorization: Bearer <access_token>"

curl http://localhost:8080/api/categories/1/books \
  -H "Authorization: Bearer <access_token>"
```

### Logout

```sh
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer <access_token>"
```

## Validasi dan Status HTTP

- `400 Bad Request`: JSON, ID, atau field permintaan tidak valid.
- `401 Unauthorized`: kredensial salah, token tidak valid, kedaluwarsa, atau sudah dicabut.
- `404 Not Found`: resource yang diminta tidak ditemukan.
- `409 Conflict`: username atau nama kategori sudah digunakan.
- `201 Created`: akun, kategori, atau buku berhasil dibuat.
- `200 OK`: operasi berhasil.
- `500 Internal Server Error`: operasi gagal akibat masalah internal.

Untuk data buku, judul, deskripsi, URL gambar, tahun terbit, harga, jumlah halaman, dan ID kategori harus diisi. Harga minimal 0, jumlah halaman harus lebih dari 0, dan tahun terbit harus di antara 1980 dan 2024. Ketebalan buku tersimpan dari logika aplikasi; field tersebut tidak diterima di request.

## Struktur Project

```text
.
├── app/
│   ├── configs/       # Koneksi database
│   ├── handlers/      # Handler HTTP
│   ├── helpers/       # Helper respons dan error
│   ├── logics/        # Operasi bisnis dan query data
│   ├── middlewares/   # Autentikasi JWT dan pencabutan token
│   ├── routes/        # Registrasi endpoint
│   └── structs/       # Model dan struktur request/response
├── cmd/api/           # Entry point server API
├── migrations/
│   ├── database/      # File migrasi SQL
│   └── execute/       # Program untuk menjalankan migrasi
├── go.mod
└── README.md
```
