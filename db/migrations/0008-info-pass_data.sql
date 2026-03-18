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
  ('p-001', 'a-001', 'hash_ozaki_001', '2022-12-17T12:28:00')
, ('p-002', 'a-002', 'hash_sena_001', '2023-03-26T18:48:00')
, ('p-003', 'a-003', 'hash_physics_001', '2023-09-10T10:23:00')
;

COMMIT;
