# Catatan Arsitektur

## Keputusan utama
- **Backend berlapis:** handler → service → repository, dengan aturan transisi status sebagai fungsi murni di `domain/`. Service bergantung pada interface repository, jadi bisa diuji dengan fake tanpa database.
- **Huma di atas Gin:** spesifikasi OpenAPI dibuat dari tipe Go, sehingga dokumentasi tidak menyimpang dari kode, dan validasi input (enum, panjang) tertulis di tipe yang sama.
- **Migrasi SQL (golang-migrate)**, bukan AutoMigrate, agar skema eksplisit dan bisa dijalankan ulang. Seed idempotent (`ON CONFLICT DO NOTHING`).
- **Akses dibatasi di backend, bukan di UI:** daftar tiket requester diberi `WHERE requester_id`, dan catatan internal difilter di SQL. Tiket orang lain dijawab 404, bukan 403.
- **Penugasan:** agent hanya mengambil tiket untuk dirinya sendiri; memindahkan ke orang lain khusus admin. Aturan ini ditegakkan di service, bukan hanya disembunyikan di frontend.
- **Update tiket bersyarat** (`WHERE status = status_lama`) mencegah dua perubahan serentak saling menimpa (409). Perubahan peran admin berjalan dalam transaksi dengan row lock agar admin terakhir tidak bisa hilang.
- **Frontend tidak menduplikasi aturan transisi:** tombol aksi dibaca dari `allowed_transitions` buatan backend.
- **Sesi memakai cookie httpOnly + SameSite=Lax**, ditambah header kustom wajib untuk request berbasis cookie yang mengubah data sebagai proteksi CSRF. Frontend dan API satu origin lewat proxy (Vite saat dev, nginx saat produksi), jadi CORS tidak dibutuhkan di produksi. Bearer tetap diterima untuk Swagger dan curl.

## Sengaja ditinggalkan
- Peran disimpan di JWT, jadi perubahan peran baru berlaku setelah token kedaluwarsa. Tidak ada refresh token atau revokasi.
- Belum ada rate limiting pada login.
- Assignee tidak bisa dilepas (dikosongkan). Tiket `closed` masih bisa menerima komentar atau perubahan prioritas lewat API.
- Pagination memakai offset, cukup untuk skala tugas ini.
- Tipe TypeScript ditulis manual mengikuti respons API.
- Bonus lampiran, WebSocket/SSE, SLA, dan statistik tidak dikerjakan.

## Perbaikan jika ada waktu
1. Cek peran ke DB per request, atau TTL pendek + refresh token.
2. Blokir komentar dan perubahan pada tiket `closed` di service.
3. Integration test terhadap Postgres asli (filter, sort, pagination, migrasi ulang, konflik serentak) dan pipeline CI untuk lint dan test.
4. Hasilkan tipe frontend dari `/openapi.json` (`openapi-typescript`).
5. Rate limit login, audit trail aktivitas tiket, lampiran dengan validasi tipe dan ukuran.
6. SLA untuk tiket urgent dan halaman statistik admin.
