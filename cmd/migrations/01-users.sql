BEGIN;

CREATE TABLE users (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "password" varchar(100) NOT NULL,
  "name" varchar(100) NOT NULL UNIQUE,
  "admin" BOOLEAN NOT NULL DEFAULT FALSE
);

COMMIT;