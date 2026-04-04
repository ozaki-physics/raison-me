BEGIN;

-- 以降の SQL を app, public スキーマ の 順で使う
SET search_path TO app, public;

INSERT INTO passwords (
    password_id
  , account_id
  , password_hash
  , iat
)
VALUES
  (
      '018f2f4e-8c1d-7b33-a1aa-4c0b0f0e1d01'
    , '018f2f4e-8c1d-7c44-b2bb-5d1c1a1f2e01'
    , 'hash_ozaki_001'
    , '2022-12-17T12:28:00'
  )
;

COMMIT;
