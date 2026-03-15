BEGIN;

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

-- passwords テーブルの作成
CREATE TABLE IF NOT EXISTS passwords (
  passwords_id TEXT PRIMARY KEY,
  account_id TEXT NOT NULL REFERENCES accounts(account_id),
  password_hash TEXT NOT NULL,
  iat TIMESTAMP NOT NULL
);

COMMIT;
