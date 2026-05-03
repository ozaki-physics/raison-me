-- bundle generated at 2026-01-18_084303
BEGIN;
CREATE TABLE IF NOT EXISTS public.schema_migrations (
  id         bigserial PRIMARY KEY,
  filename   text NOT NULL UNIQUE,
  sha256     text NOT NULL,
  applied_at timestamptz NOT NULL DEFAULT now()
);

-- >>> START 0000-init.sql >>>

-- アプリ用スキーマ を 作成
CREATE SCHEMA IF NOT EXISTS app;
-- 以降の SQL を app, public スキーマ の 順で使う
-- トランザクション や セッション のスコープ
SET search_path TO app, public;

-- <<< END 0000-init.sql <<<

-- >>> START 0001-capital-coin.sql >>>

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

-- <<< END 0001-capital-coin.sql <<<

-- >>> START 0002-capital-coin_transaction.sql >>>

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

-- <<< END 0002-capital-coin_transaction.sql <<<

-- >>> START 0003-capital-coin_data.sql >>>

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

INSERT INTO capital_coin (
    id
  , symbol
  , coin_market_cap_id
  , created_at
  , updated_at
)
VALUES
  (1, 'BTC', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (2, 'ETH', 1027, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (3, 'DOT', 6636, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (4, 'BAT', 1697, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (5, 'ENJ', 2130, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (6, 'XTZ', 2011, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (7, 'ATOM', 3794, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
;

-- <<< END 0003-capital-coin_data.sql <<<

-- >>> START 0004-capital-coin_transaction_data.sql >>>

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

INSERT INTO capital_coin_transaction (
    id
  , coin_id
  , symbol
  , side
  , price_rate
  , size
  , fee
  , traded_at
  , created_at
  , updated_at
)
VALUES
  (1001, 1, 'BTC', 0, 2300000, 0.003, 0, '2022-11-13 18:06:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1002, 1, 'BTC', 0, 2350000, 0.002, 0, '2022-11-11 23:20:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1003, 1, 'BTC', 0, 2400000, 0.003, 0, '2022-11-10 12:11:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1004, 1, 'BTC', 0, 2550000, 0.005, 1, '2022-11-09 22:09:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1005, 1, 'BTC', 0, 2590000, 0.003, 0, '2022-11-09 20:17:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1006, 1, 'BTC', 0, 2705000, 0.001, 0, '2022-11-09 08:52:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1007, 1, 'BTC', 0, 2715000, 0.01, 2, '2022-11-09 08:46:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1008, 1, 'BTC', 0, 2709823, 0.002, -3, '2022-11-09 08:14:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1009, 1, 'BTC', 0, 2740000, 0.002, 0, '2022-08-29 08:03:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1010, 1, 'BTC', 0, 2900000, 0.002, 0, '2022-08-22 16:51:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1011, 1, 'BTC', 0, 2580000, 0.0018, 0, '2022-07-02 12:42:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1012, 1, 'BTC', 0, 2853096, 0.0012, -2, '2022-06-15 14:45:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1013, 1, 'BTC', 0, 3750000, 0.003, 1, '2022-05-12 05:54:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1014, 1, 'BTC', 0, 4150000, 0.001, 0, '2022-01-22 09:06:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1015, 1, 'BTC', 0, 4838000, 0.001, -3, '2022-01-07 23:16:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1016, 1, 'BTC', 0, 5300000, 0.001, 0, '2022-01-06 03:11:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1017, 1, 'BTC', 0, 5388000, 0.001, 0, '2022-01-05 20:27:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1018, 2, 'ETH', 0, 160000, 0.04, 0, '2022-11-10 08:24:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1019, 2, 'ETH', 0, 170000, 0.07, 0, '2022-11-10 00:44:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1020, 2, 'ETH', 0, 175000, 0.05, 0, '2022-11-09 23:32:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1021, 2, 'ETH', 0, 177000, 0.04, 0, '2022-11-09 20:25:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1022, 2, 'ETH', 0, 177000, 0.01, 0, '2022-11-09 20:24:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1023, 2, 'ETH', 0, 194952, 0.1, -10, '2022-11-09 08:42:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1024, 2, 'ETH', 0, 194100, 0.01, -1, '2022-11-09 08:15:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1025, 2, 'ETH', 0, 185000, 0.025, 0, '2022-09-24 02:27:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1026, 2, 'ETH', 0, 186000, 0.015, 0, '2022-09-19 11:17:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1027, 2, 'ETH', 0, 200000, 0.05, 1, '2022-09-19 01:23:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1028, 2, 'ETH', 0, 200000, 0.03, 0, '2022-08-28 01:23:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1029, 2, 'ETH', 0, 215000, 0.03, 0, '2022-08-21 15:43:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1030, 2, 'ETH', 0, 140000, 0.05, 0, '2022-07-02 12:33:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1031, 2, 'ETH', 0, 138700, 0.01, -1, '2022-06-15 17:56:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1032, 2, 'ETH', 0, 150867, 0.03, -3, '2022-06-15 14:44:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1033, 2, 'ETH', 0, 250000, 0.05, 0, '2022-05-12 13:20:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1034, 2, 'ETH', 0, 275000, 0.05, 1, '2022-02-24 23:02:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1035, 2, 'ETH', 0, 250000, 0.01, 0, '2022-01-24 21:43:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1036, 2, 'ETH', 0, 296000, 0.01, 0, '2022-01-22 09:03:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1037, 2, 'ETH', 0, 326000, 0.02, 0, '2022-01-21 12:24:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1038, 2, 'ETH', 0, 373000, 0.01, 0, '2022-01-07 15:10:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1039, 2, 'ETH', 0, 440000, 0.01, 0, '2022-01-05 20:33:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1040, 3, 'DOT', 0, 0, 0.002020958, 0, '2022-02-24 23:01:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1041, 3, 'DOT', 0, 1763, 1, 0, '2022-02-24 23:01:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1042, 3, 'DOT', 0, 2093, 1, 0, '2022-01-24 00:34:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1043, 4, 'BAT', 0, 84.242, 10, 0, '2022-01-24 00:34:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1044, 5, 'ENJ', 0, 143.525, 5, 0, '2022-02-24 23:00:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1045, 5, 'ENJ', 0, 182.727, 5, 0, '2022-01-24 00:37:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1046, 6, 'XTZ', 0, 0, 0.032891, 0, '2022-12-12 11:18:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1047, 6, 'XTZ', 0, 0, 0.035775, 0, '2022-11-10 11:16:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1048, 6, 'XTZ', 0, 0, 0.032075, 0, '2022-10-11 11:15:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1049, 6, 'XTZ', 0, 0, 0.039129, 0, '2022-09-12 13:53:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1050, 6, 'XTZ', 0, 0, 0.038452, 0, '2022-08-10 13:21:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1051, 6, 'XTZ', 0, 0, 0.048924, 0, '2022-07-11 12:05:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1052, 6, 'XTZ', 0, 0, 0.033157, 0, '2022-06-10 12:02:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1053, 6, 'XTZ', 0, 0, 0.027322, 0, '2022-05-10 11:57:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1054, 6, 'XTZ', 0, 0, 0.013571, 0, '2022-04-11 10:47:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1055, 6, 'XTZ', 0, 317.035, 5, 0, '2022-02-24 22:59:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1056, 6, 'XTZ', 0, 348.615, 5, 0, '2022-01-24 00:41:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1057, 7, 'ATOM', 0, 0, 0.001732, 0, '2022-12-13 12:52:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
, (1058, 7, 'ATOM', 0, 2620, 1, 0, '2022-02-24 23:02:00', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
;


-- <<< END 0004-capital-coin_transaction_data.sql <<<

-- >>> START 0005-info-user.sql >>>

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

-- users テーブルの作成
CREATE TABLE IF NOT EXISTS users (
  account_id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL
);

-- <<< END 0005-info-user.sql <<<

-- >>> START 0006-info-user_data.sql >>>

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

INSERT INTO users (
    account_id
  , user_id
  , name
)
VALUES
  ('a-001', 'ozaki', 'オザキ')
, ('a-002', 'sena', 'セナ')
, ('a-003', 'physics', '物理')
;

-- <<< END 0006-info-user_data.sql <<<

-- >>> START 0007-info-pass.sql >>>

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

-- users テーブルの作成
CREATE TABLE IF NOT EXISTS pass (
  id TEXT PRIMARY KEY,
  account_id TEXT NOT NULL REFERENCES users(account_id),
  password_hash TEXT NOT NULL,
  iat TIMESTAMP NOT NULL
);

-- <<< END 0007-info-pass.sql <<<

-- >>> START 0008-info-pass_data.sql >>>

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

INSERT INTO pass (
    id
  , account_id
  , password_hash
  , iat
)
VALUES
  ('p-001', 'a-001', 'hash_ozaki_001', '2022-12-17T12:28:00')
, ('p-002', 'a-002', 'hash_sena_001', '2023-03-26T18:48:00')
, ('p-003', 'a-003', 'hash_physics_001', '2023-09-10T10:23:00')
;

-- <<< END 0008-info-pass_data.sql <<<

DO $$ BEGIN
  IF EXISTS (SELECT 1 FROM public.schema_migrations WHERE filename = '0000-init.sql' AND sha256 <> 'e1d7c46c0d8b936c8abeadbd8510602f1b7295d54ec0067e973d1d3a241d24c6') THEN
    RAISE EXCEPTION 'Hash mismatch for 0000-init.sql';
  END IF;
  INSERT INTO public.schema_migrations(filename, sha256)
  VALUES ('0000-init.sql', 'e1d7c46c0d8b936c8abeadbd8510602f1b7295d54ec0067e973d1d3a241d24c6')
  ON CONFLICT (filename) DO NOTHING;
  IF EXISTS (SELECT 1 FROM public.schema_migrations WHERE filename = '0001-capital-coin.sql' AND sha256 <> '79c4a2f6f32d90e574b49e3efda29bfb229ef156098acd176c61e7c43ca49d4b') THEN
    RAISE EXCEPTION 'Hash mismatch for 0001-capital-coin.sql';
  END IF;
  INSERT INTO public.schema_migrations(filename, sha256)
  VALUES ('0001-capital-coin.sql', '79c4a2f6f32d90e574b49e3efda29bfb229ef156098acd176c61e7c43ca49d4b')
  ON CONFLICT (filename) DO NOTHING;
  IF EXISTS (SELECT 1 FROM public.schema_migrations WHERE filename = '0002-capital-coin_transaction.sql' AND sha256 <> '1c0441e33d1f5fafbcb8e6c04447e4a6f5f6ea496f69c495f37ed5d5a4474352') THEN
    RAISE EXCEPTION 'Hash mismatch for 0002-capital-coin_transaction.sql';
  END IF;
  INSERT INTO public.schema_migrations(filename, sha256)
  VALUES ('0002-capital-coin_transaction.sql', '1c0441e33d1f5fafbcb8e6c04447e4a6f5f6ea496f69c495f37ed5d5a4474352')
  ON CONFLICT (filename) DO NOTHING;
  IF EXISTS (SELECT 1 FROM public.schema_migrations WHERE filename = '0003-capital-coin_data.sql' AND sha256 <> '3b61e517e18b5158f9a6bc68b082d72a1bc12e73e713b4e2f47515c54b480a99') THEN
    RAISE EXCEPTION 'Hash mismatch for 0003-capital-coin_data.sql';
  END IF;
  INSERT INTO public.schema_migrations(filename, sha256)
  VALUES ('0003-capital-coin_data.sql', '3b61e517e18b5158f9a6bc68b082d72a1bc12e73e713b4e2f47515c54b480a99')
  ON CONFLICT (filename) DO NOTHING;
  IF EXISTS (SELECT 1 FROM public.schema_migrations WHERE filename = '0004-capital-coin_transaction_data.sql' AND sha256 <> 'c3f0bf5f4ab4f9f97b77aa513950eb30f1502d762ab4ff6e7b1c3b55ec9739ed') THEN
    RAISE EXCEPTION 'Hash mismatch for 0004-capital-coin_transaction_data.sql';
  END IF;
  INSERT INTO public.schema_migrations(filename, sha256)
  VALUES ('0004-capital-coin_transaction_data.sql', 'c3f0bf5f4ab4f9f97b77aa513950eb30f1502d762ab4ff6e7b1c3b55ec9739ed')
  ON CONFLICT (filename) DO NOTHING;
  IF EXISTS (SELECT 1 FROM public.schema_migrations WHERE filename = '0005-info-user.sql' AND sha256 <> 'cfbbd45eed997f442b2ccb6eb2b2ce4f726e16701f5a59677246ffbd41c1e619') THEN
    RAISE EXCEPTION 'Hash mismatch for 0005-info-user.sql';
  END IF;
  INSERT INTO public.schema_migrations(filename, sha256)
  VALUES ('0005-info-user.sql', 'cfbbd45eed997f442b2ccb6eb2b2ce4f726e16701f5a59677246ffbd41c1e619')
  ON CONFLICT (filename) DO NOTHING;
  IF EXISTS (SELECT 1 FROM public.schema_migrations WHERE filename = '0006-info-user_data.sql' AND sha256 <> '9752f1defd31a092a84b0337990260c7970a9f6727cfa82a0b2dc74e7b36bcbe') THEN
    RAISE EXCEPTION 'Hash mismatch for 0006-info-user_data.sql';
  END IF;
  INSERT INTO public.schema_migrations(filename, sha256)
  VALUES ('0006-info-user_data.sql', '9752f1defd31a092a84b0337990260c7970a9f6727cfa82a0b2dc74e7b36bcbe')
  ON CONFLICT (filename) DO NOTHING;
  IF EXISTS (SELECT 1 FROM public.schema_migrations WHERE filename = '0007-info-pass.sql' AND sha256 <> '9fc11f05757d9bbd5f8356b414f02b7020c421370169ff1e3112db6a7b69402d') THEN
    RAISE EXCEPTION 'Hash mismatch for 0007-info-pass.sql';
  END IF;
  INSERT INTO public.schema_migrations(filename, sha256)
  VALUES ('0007-info-pass.sql', '9fc11f05757d9bbd5f8356b414f02b7020c421370169ff1e3112db6a7b69402d')
  ON CONFLICT (filename) DO NOTHING;
  IF EXISTS (SELECT 1 FROM public.schema_migrations WHERE filename = '0008-info-pass_data.sql' AND sha256 <> '94d1f2a31f8fccf2c626a463c9e54ace0d46d8b3be065a22d0b9d197f3b67a0c') THEN
    RAISE EXCEPTION 'Hash mismatch for 0008-info-pass_data.sql';
  END IF;
  INSERT INTO public.schema_migrations(filename, sha256)
  VALUES ('0008-info-pass_data.sql', '94d1f2a31f8fccf2c626a463c9e54ace0d46d8b3be065a22d0b9d197f3b67a0c')
  ON CONFLICT (filename) DO NOTHING;
END; $$;
COMMIT;
