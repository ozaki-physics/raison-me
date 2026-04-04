BEGIN;

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

INSERT INTO accounts (
    account_id
  , user_id
  , user_name
)
VALUES
  (
      '018f2f4e-8c1d-7c44-b2bb-5d1c1a1f2e01'
    , 'ozaki'
    , 'オザキ'
  )
;

COMMIT;
