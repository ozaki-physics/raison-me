BEGIN;

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

INSERT INTO accounts (
    account_id
  , user_id
  , user_name
)
VALUES
  ('a-001', 'ozaki', 'オザキ')
, ('a-002', 'sena', 'セナ')
, ('a-003', 'physics', '物理')
;

COMMIT;
