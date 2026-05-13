CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('customer', 'admin');
CREATE TYPE cinema_status AS ENUM ('active', 'inactive');
CREATE TYPE seat_type AS ENUM ('regular', 'premium', 'vip');
CREATE TYPE schedule_status AS ENUM ('active', 'cancelled', 'finished');
CREATE TYPE lock_status AS ENUM ('locked', 'released', 'expired');
CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed', 'refunded');
CREATE TYPE ticket_status AS ENUM ('active', 'used', 'cancelled');
CREATE TYPE refund_status AS ENUM ('pending', 'processed', 'failed');

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          VARCHAR(120) NOT NULL,
    email         VARCHAR(120) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone         VARCHAR(20),
    role          user_role NOT NULL DEFAULT 'customer',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE cinemas (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(150) NOT NULL,
    city       VARCHAR(80) NOT NULL,
    address    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE studios (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cinema_id  UUID NOT NULL REFERENCES cinemas(id) ON DELETE CASCADE,
    name       VARCHAR(80) NOT NULL,
    capacity   INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE seats (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    studio_id   UUID NOT NULL REFERENCES studios(id) ON DELETE CASCADE,
    row_label   CHAR(1) NOT NULL,
    seat_number INT NOT NULL,
    seat_type   seat_type NOT NULL DEFAULT 'regular',
    UNIQUE (studio_id, row_label, seat_number)
);

CREATE TABLE movies (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title            VARCHAR(200) NOT NULL,
    duration_minutes INT NOT NULL,
    genre            VARCHAR(80),
    rating           VARCHAR(10),
    synopsis         TEXT,
    poster_url       TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE schedules (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    movie_id   UUID NOT NULL REFERENCES movies(id),
    studio_id  UUID NOT NULL REFERENCES studios(id),
    show_time  TIMESTAMPTZ NOT NULL,
    price      NUMERIC(12,2) NOT NULL,
    status     schedule_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE seat_locks (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    schedule_id UUID NOT NULL REFERENCES schedules(id),
    seat_id     UUID NOT NULL REFERENCES seats(id),
    user_id     UUID NOT NULL REFERENCES users(id),
    locked_until TIMESTAMPTZ NOT NULL,
    status      lock_status NOT NULL DEFAULT 'locked',
    UNIQUE (schedule_id, seat_id, status)
);

CREATE TABLE transactions (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id        UUID NOT NULL REFERENCES users(id),
    schedule_id    UUID NOT NULL REFERENCES schedules(id),
    payment_method VARCHAR(50),
    payment_status payment_status NOT NULL DEFAULT 'pending',
    total_amount   NUMERIC(12,2) NOT NULL,
    pg_reference   VARCHAR(200),
    paid_at        TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tickets (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    seat_id        UUID NOT NULL REFERENCES seats(id),
    ticket_code    VARCHAR(40) UNIQUE NOT NULL,
    status         ticket_status NOT NULL DEFAULT 'active',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE refunds (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    amount         NUMERIC(12,2) NOT NULL,
    reason         TEXT,
    status         refund_status NOT NULL DEFAULT 'pending',
    pg_reference   VARCHAR(200),
    processed_at   TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_schedules_show_time ON schedules(show_time);
CREATE INDEX idx_schedules_status ON schedules(status);
CREATE INDEX idx_seat_locks_schedule ON seat_locks(schedule_id);
CREATE INDEX idx_tickets_transaction ON tickets(transaction_id);
CREATE INDEX idx_transactions_user ON transactions(user_id);
