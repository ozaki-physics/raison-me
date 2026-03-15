BEGIN;

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

-- Create capital_coin_transaction table
CREATE TABLE IF NOT EXISTS capital_coin_transaction (
  -- 本当は UUID などで管理したいが 開発初期段階では 一旦 integer 型で運用する
  -- id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  id integer PRIMARY KEY,
  -- 外部キー: capital_coin テーブルの id を参照
  coin_id integer NOT NULL REFERENCES capital_coin(id) ON DELETE RESTRICT,
  -- 銘柄 symbol は 変わるかもしれないが 取引した時点の文字列で保存する
  symbol text NOT NULL,
  -- 0: buy(買), 1: sell(売)
  side smallint NOT NULL CHECK (side IN (0, 1)),
  -- 約定レート
  price_rate decimal NOT NULL,
  -- 約定数量
  size DOUBLE PRECISION NOT NULL,
  -- 取引手数料
  fee decimal NOT NULL,
  -- 約定日時
  traded_at TIMESTAMP NOT NULL,
  -- 以下は 一般的な管理用カラム
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL
);

COMMIT;
