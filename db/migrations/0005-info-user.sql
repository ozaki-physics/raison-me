BEGIN;

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

-- accounts テーブルの作成
CREATE TABLE IF NOT EXISTS accounts (
  account_id UUID PRIMARY KEY,
  user_id TEXT NOT NULL UNIQUE,
  user_name TEXT NOT NULL
);

COMMIT;
