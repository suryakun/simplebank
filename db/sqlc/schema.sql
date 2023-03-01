CREATE TABLE accounts (
  id BIGSERIAL PRIMARY KEY,
  owner varchar NOT NULL,
  balance bigint NOT NULL,
  currency varchar NOT NUll,
  created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE entries (
  id BIGSERIAL PRIMARY KEY,
  account_id BIGSERIAL NOT NULL,
  amount BIGINT NOT NULL DEFAULT 0,
  created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE transfers (
  id BIGSERIAL PRIMARY KEY,
  from_account_id BIGINT NOT NULL,
  to_account_id BIGINT NOT NULL,
  amount BIGINT NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now())
);