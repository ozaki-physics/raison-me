BEGIN;

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

-- capital_coin テーブルの作成
CREATE TABLE IF NOT EXISTS capital_coin (
  -- 本当は UUID など で管理したいが 開発初期段階では 一旦 integer 型で運用する
  -- id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  id integer PRIMARY KEY,
  symbol text NOT NULL UNIQUE,
  coin_market_cap_id integer NOT NULL,
  -- 以下は 一般的な管理用カラム
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);

COMMIT;
