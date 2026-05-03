BEGIN;

-- アプリ用スキーマ を 作成
CREATE SCHEMA IF NOT EXISTS app;
-- 以降の SQL を app, public スキーマ の 順で使う
-- トランザクション や セッション のスコープ
SET search_path TO app, public;

COMMIT;
