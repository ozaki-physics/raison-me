BEGIN;

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

-- passwords テーブルの作成
CREATE TABLE IF NOT EXISTS passwords (
  password_id UUID PRIMARY KEY,
  account_id UUID NOT NULL REFERENCES accounts(account_id),
  password_hash TEXT NOT NULL,
  iat TIMESTAMPTZ NOT NULL
);

COMMIT;
