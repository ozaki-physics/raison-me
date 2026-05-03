BEGIN;

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

-- refresh_tokens テーブルの作成
CREATE TABLE IF NOT EXISTS refresh_tokens (
  refresh_token_id UUID PRIMARY KEY,
  account_id UUID NOT NULL REFERENCES accounts(account_id),
  token_hash TEXT NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

COMMIT;
