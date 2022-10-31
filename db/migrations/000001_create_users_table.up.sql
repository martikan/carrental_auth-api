CREATE TABLE IF NOT EXISTS "users" (
    "id" bigserial PRIMARY KEY,
    "email" varchar(255) UNIQUE NOT NULL,
    "password" varchar(100) NOT NULL,
    "first_name" varchar(50) NOT NULL,
    "last_name" varchar(50) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
