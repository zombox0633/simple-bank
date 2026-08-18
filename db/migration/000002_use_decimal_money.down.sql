ALTER TABLE transfers
  DROP CONSTRAINT IF EXISTS transfers_amount_positive;

ALTER TABLE accounts
  DROP CONSTRAINT IF EXISTS accounts_balance_nonnegative;

-- Rolling back to BIGINT cannot preserve fractional digits. Round explicitly
-- so the lossy behavior is visible instead of relying on an implicit cast.
ALTER TABLE transfers
  ALTER COLUMN amount TYPE bigint
  USING round(amount)::bigint;

ALTER TABLE entries
  ALTER COLUMN amount TYPE bigint
  USING round(amount)::bigint;

ALTER TABLE accounts
  ALTER COLUMN balance TYPE bigint
  USING round(balance)::bigint;
