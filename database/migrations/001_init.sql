-- Schema awal — padanan dari prisma/schema.prisma di proyek NestJS.
-- Tabel session/account/verification milik Better Auth tidak ikut dibawa
-- karena auth di versi Go memakai JWT stateless.

-- enum Role { PARTICIPANT ADMIN }
DO $$ BEGIN
    CREATE TYPE role AS ENUM ('PARTICIPANT', 'ADMIN');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- model user
-- id memakai TEXT (bukan UUID) agar perilakunya sama dengan Prisma:
-- lookup dengan id yang bukan UUID valid tetap menghasilkan "not found",
-- bukan error parsing.
CREATE TABLE IF NOT EXISTS users (
    id         TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL UNIQUE,
    password   TEXT NOT NULL, -- bcrypt hash, tidak pernah dikirim ke client
    role       role NOT NULL DEFAULT 'PARTICIPANT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- model Hackathon
CREATE TABLE IF NOT EXISTS hackathons (
    id          TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    name        TEXT NOT NULL,
    description TEXT,
    start_date  TIMESTAMPTZ NOT NULL,
    end_date    TIMESTAMPTZ NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT false,
    author_id   TEXT NOT NULL REFERENCES users (id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- model HackathonParticipant, dengan @@unique([hackathonId, userId])
CREATE TABLE IF NOT EXISTS hackathon_participants (
    id           TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    hackathon_id TEXT NOT NULL REFERENCES hackathons (id) ON DELETE CASCADE,
    user_id      TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    joined_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (hackathon_id, user_id)
);
