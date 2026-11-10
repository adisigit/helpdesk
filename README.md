# Helpdesk Internal

Aplikasi helpdesk untuk menggantikan email dan spreadsheet. Karyawan membuat tiket, agent menanganinya, admin mengelola pengguna dan peran.

**Stack:** Go (Gin + Huma/OpenAPI 3.1, GORM, golang-migrate), PostgreSQL 16, Vue 3 + TypeScript (Vite, Pinia, Vue Router, Tailwind), nginx.

## Menjalankan

Butuh Docker dengan Compose v2.

```bash
cp .env.example .env
docker compose up --build
```

| Layanan | URL |
|---|---|
| Aplikasi (frontend) | http://localhost:5173 |
| API | http://localhost:8080/api |
| Dokumentasi API (OpenAPI/Swagger UI) | http://localhost:8080/docs |
| Spesifikasi OpenAPI | http://localhost:8080/openapi.json dan `/openapi.yaml` |

Migrasi dan seed berjalan otomatis saat backend start dan aman diulang.
Reset database: `docker compose down -v`.

### Akun seed

Password semua akun diambil dari `SEED_PASSWORD` di `.env` (nilai default di `.env.example`: `Password123!`).

| Peran | Email |
|---|---|
| Admin | admin@helpdesk.com |
| Agent | agent@helpdesk.com |
| Requester | requester@helpdesk.com |

Akun baru lewat halaman Register selalu berperan requester. Hanya admin yang bisa mengubah peran.

### Mode pengembangan (hot reload)

```bash
docker compose -f docker-compose.yaml -f docker-compose.dev.yaml up --build
```

Frontend di http://localhost:5173 (Vite), backend di-reload oleh air. Butuh Compose v2.24 ke atas.

## Konfigurasi (environment variable)

Lihat `.env.example`. Tidak ada kredensial di source code.

| Variabel | Fungsi |
|---|---|
| `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` | Database (compose merakit `DATABASE_URL`) |
| `JWT_SECRET` | Kunci tanda tangan JWT, minimal 16 karakter |
| `JWT_TTL` | Masa berlaku token, default `3h` |
| `COOKIE_SECURE` | `true` bila disajikan lewat HTTPS |
| `CORS_ORIGINS` | Daftar origin dipisah koma |
| `SEED_ON_START`, `SEED_PASSWORD` | Seed satu akun per peran |
| `PORT` | Port backend, default `8080` |

## Autentikasi

- Browser memakai cookie `token` (`HttpOnly`, `SameSite=Lax`). Request berbasis cookie yang mengubah data wajib menyertakan header `X-Requested-With` (proteksi CSRF); frontend menambahkannya otomatis.
- Klien API (curl, Swagger) bisa memakai `Authorization: Bearer <access_token>` dari respons login.
- Di `/docs`, klik **Authorize** lalu tempel `access_token`.

## Endpoint

Semua selain auth dan health butuh login.

| Method | Path | Fungsi | Akses |
|---|---|---|---|
| GET | `/api/health` | Health check | Publik |
| POST | `/api/auth/register` | Daftar akun requester | Publik |
| POST | `/api/auth/login` | Login, set cookie, kembalikan JWT | Publik |
| POST | `/api/auth/logout` | Hapus cookie sesi | Publik |
| GET | `/api/me` | Profil pengguna saat ini | Login |
| GET | `/api/tickets` | Daftar tiket: `status`, `priority`, `assignee` (UUID/`me`/`none`), `q`, `sort`, `page`, `page_size` | Login (requester hanya tiket miliknya) |
| POST | `/api/tickets` | Buat tiket | Login |
| GET | `/api/tickets/{id}` | Detail tiket + komentar | Pemilik, agent, admin |
| PATCH | `/api/tickets/{id}` | Ubah status, prioritas, assignee | Sesuai tabel transisi |
| POST | `/api/tickets/{id}/comments` | Komentar atau catatan internal | Pemilik, agent, admin |
| GET | `/api/assignees` | Daftar agent dan admin untuk dropdown assignee | Admin |
| GET | `/api/users` | Daftar pengguna | Admin |
| PATCH | `/api/users/{id}/role` | Ubah peran | Admin |

### Alur status

| Dari | Ke | Oleh |
|---|---|---|
| open | in_progress | Agent, Admin |
| in_progress | resolved | Agent, Admin |
| resolved | closed | Requester (pemilik), Agent, Admin |
| resolved | open | Requester (pemilik) |
| open | closed | Requester (pemilik), Admin |

- Di luar tabel: `422`. Kombinasi valid tetapi peran tidak berhak: `403`.
- Requester yang mengakses tiket orang lain mendapat `404`, jadi keberadaan tiket tidak bocor.
- Respons tiket memuat `allowed_transitions`, dan frontend menampilkan tombol dari field itu.

### Penugasan tiket

- Agent hanya bisa mengambil tiket untuk dirinya sendiri (tombol "Ambil tiket", atau otomatis saat memulai tiket tanpa assignee), dan tidak bisa merebut tiket yang dipegang orang lain.
- Admin bisa menugaskan atau memindahkan tiket ke agent atau admin mana pun.

## Test

```bash
cd backend
go test ./...

cd ../frontend
npm ci && npm run build
```

Cakupan utama:
- tabel transisi status (table-driven, termasuk penjaga jumlah kombinasi yang diizinkan)
- pembatasan akses requester pada baca, ubah, dan komentar tiket orang lain
- catatan internal tidak terlihat oleh requester
- agent hanya boleh mengambil tiket untuk dirinya sendiri
- register dan login (password salah, email duplikat)
- JWT dan middleware peran, termasuk wajib header CSRF pada request berbasis cookie

## Struktur

```
backend/
  cmd/api/              entrypoint
  db/migrations/        SQL migrasi (golang-migrate)
  internal/
    config/ auth/ database/ (models, seed)
    domain/             aturan transisi status (fungsi murni)
    repository/         akses DB (interface + GORM)
    service/            logika bisnis
    server/             handler, router, middleware (Huma)
frontend/src/
  api/ stores/ router/ views/ components/
```
