# Hackathon Management App (Go)

Port 1:1 dari aplikasi NestJS `hackaton` ke **Go** dengan struktur **MVC** tradisional:
PostgreSQL lokal (raw SQL, tanpa ORM), autentikasi **JWT** + bcrypt, dan role
`ADMIN` / `PARTICIPANT`.

## Pemetaan konsep NestJS → Go

| NestJS | Go (proyek ini) |
|---|---|
| `main.ts` (bootstrap) | [main.go](main.go) |
| `ConfigModule` | [config/config.go](config/config.go) |
| `PrismaService` + `prisma migrate` | [database/database.go](database/database.go) + [database/migrations/001_init.sql](database/migrations/001_init.sql) |
| Model Prisma | [models/](models/) (struct + tag JSON) |
| DTO + class-validator | [dto/](dto/) (tag `binding` + validator custom `future`) |
| `*.controller.ts` | [controllers/](controllers/) |
| `*.service.ts` (Prisma call) | [services/](services/) (raw SQL) |
| Better Auth + `AuthGuard` + `@Roles` | [middleware/auth_middleware.go](middleware/auth_middleware.go) (JWT) |
| `ResponseInterceptor` | [utils/response.go](utils/response.go) |
| Exception layer (`HttpException`) | [utils/apperror.go](utils/apperror.go) |
| Wiring module (`app.module.ts`) | [routes/routes.go](routes/routes.go) |

Alur satu request: `routes` → `middleware` (auth/role) → `controller` (bind +
validasi DTO) → `service` (business logic + SQL) → `utils` (envelope response).

## Menjalankan

1. Pastikan PostgreSQL lokal berjalan, lalu buat database:
   ```sql
   CREATE DATABASE hackathon_db;
   ```
2. Salin `.env.example` menjadi `.env`, sesuaikan `DATABASE_URL` dan `JWT_SECRET`.
3. Jalankan (migrasi schema otomatis dijalankan saat startup):
   ```bash
   go run .
   ```
   Server hidup di `http://localhost:3000`.

## Endpoint

| Method | Path | Akses | Keterangan |
|---|---|---|---|
| GET | `/` | Publik | Hello World |
| POST | `/auth/register` | Publik | Daftar — role selalu `PARTICIPANT` |
| POST | `/auth/login` | Publik | Login → JWT |
| GET | `/me` | Login | User yang sedang login |
| GET | `/user/all` | ADMIN | Semua user |
| GET | `/user/:id` | Login | Satu user |
| GET | `/hackaton` | Publik | Semua hackathon |
| GET | `/hackaton/:id` | Publik | Satu hackathon |
| POST | `/hackaton` | ADMIN | Buat hackathon (author = admin yang login) |
| PATCH | `/hackaton/:id` | ADMIN | Partial update |
| DELETE | `/hackaton/:id` | ADMIN | Hapus (peserta ikut terhapus, cascade) |
| POST | `/hackaton/:id/join` | PARTICIPANT | Join hackathon aktif yang belum berakhir |

Endpoint terproteksi memakai header `Authorization: Bearer <token>`.

### Membuat admin

Sama seperti versi NestJS, role tidak pernah diterima dari input pendaftaran —
promosikan lewat SQL:

```sql
UPDATE users SET role = 'ADMIN' WHERE email = 'admin@mail.com';
```

Perubahan role langsung berlaku tanpa login ulang, karena middleware memuat
user dari database pada tiap request.

## Format response

Sukses (padanan `ResponseInterceptor`):

```json
{ "statusCode": 200, "message": "Success", "data": { ... } }
```

Error (padanan exception layer NestJS):

```json
{ "statusCode": 404, "message": "Hackathon with id \"x\" not found", "error": "Not Found" }
```

Error validasi — `message` berupa array per field (padanan `exceptionFactory`
di `main.ts`):

```json
{
  "statusCode": 400,
  "message": [
    { "property": "name", "message": "name must be longer than or equal to 3 characters" },
    { "property": "startsAt", "message": "startsAt must be a future date" }
  ],
  "error": "Bad Request"
}
```

## Aturan validasi

- `name`: wajib, min 3 karakter
- `description`: opsional, 10–1000 karakter
- `startsAt` / `endsAt`: wajib (create), harus tanggal di masa depan (RFC3339, mis. `2026-08-01T00:00:00Z`)
- `isActive`: opsional, boolean
- Register: `name` wajib, `email` valid, `password` min 8
- PATCH memakai aturan yang sama tetapi semua field opsional (partial update)

## Penyederhanaan dari versi NestJS

- **Better Auth → JWT stateless**: tabel `session`, `account`, `verification`
  tidak dibawa; identitas dibawa di dalam token.
- **Arcjet** (rate limit/bot detection) tidak dibawa — di Go bisa ditambahkan
  belakangan sebagai middleware.
- **Prisma → raw SQL**: schema dikelola lewat file migrasi SQL biasa yang
  idempotent dan dijalankan otomatis saat startup.
- Model `Post` di schema Prisma tidak dibawa karena tidak ada endpoint yang
  memakainya.
- Field `emailVerified` dan `image` pada user (kebutuhan Better Auth) tidak dibawa.
