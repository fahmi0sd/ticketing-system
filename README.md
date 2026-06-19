# Ticketing System API

REST API untuk pemesanan tiket transportasi (bus, kereta, pesawat) 

## Tech Stack

- **Language:** Go
- **Framework:** Echo
- **Database:** PostgreSQL (GORM)
- **Auth:** JWT
- **3rd Party:** Midtrans (Payment Gateway)

## Fitur

- Register & Login dengan JWT
- Cari rute transportasi (filter by origin, destination, type, tanggal)
- Booking tiket dengan pengecekan kuota otomatis
- Integrasi pembayaran via Midtrans Snap
- Webhook untuk update status pembayaran otomatis
- Cancel booking

## Cara Menjalankan

1. Clone repo & masuk ke folder project

   ```bash
   git clone <repo-url>
   cd ticketing-api
   ```

2. Copy `.env.example` jadi `.env`, lalu isi konfigurasi (DB, JWT secret, Midtrans key)

   ```bash
   cp .env.example .env
   ```

3. Install dependencies

   ```bash
   go mod tidy
   ```

4. Jalankan migration database

   ```bash
   psql -U postgres -d ticketing_db -f ddl.sql
   ```

5. Jalankan server

   ```bash
   go run ./app/echo-server
   ```

   Server berjalan di `http://localhost:8080`

## Testing

```bash
go test ./service/booking/... -v
```

## Endpoint

| Method | Endpoint                      | Keterangan                  |
|--------|--------------------------------|------------------------------|
| POST   | `/api/users/register`          | Daftar akun baru             |
| POST   | `/api/users/login`             | Login, dapat JWT token       |
| GET    | `/api/users/me`                | Profil user (JWT)            |
| GET    | `/api/routes`                  | Cari rute transportasi       |
| GET    | `/api/routes/:id`              | Detail rute                  |
| POST   | `/api/bookings`                | Buat booking + payment link  |
| GET    | `/api/bookings`                | Riwayat booking (JWT)        |
| GET    | `/api/bookings/:id`            | Detail booking (JWT)         |
| PUT    | `/api/bookings/:id/cancel`     | Cancel booking (JWT)         |
| POST   | `/api/bookings/webhook`        | Callback Midtrans            |
| GET    | `/api/bookings/:id/payment`    | Cek status pembayaran (JWT)  |

## Struktur Folder

```
app/echo-server/     → controller, router, main.go
service/             → business logic (+ unit test)
repository/          → akses database (GORM)
pkg/                 → client 3rd party Midtrans
util/                → response
sql/                 → schema database
```

---

- **Link Deployment:** https://ticketing-system-production-ec0f.up.railway.app 
- **Link Documentation API:** https://documenter.getpostman.com/view/39224648/2sBXwvJ8gj
- **Link PPT:** https://docs.google.com/presentation/d/19WO23XmIxjRtmxiUdIz7dYjXOziNdf5-/edit?usp=sharing&ouid=109819096331892243259&rtpof=true&sd=true