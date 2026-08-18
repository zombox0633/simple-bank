-- Store money as an exact decimal with four fractional digits.
-- NUMERIC is exact; unlike REAL/DOUBLE PRECISION it does not introduce
-- binary floating-point rounding errors.
ALTER TABLE accounts
  ALTER COLUMN balance TYPE numeric(18, 4)
  USING balance::numeric(18, 4);

ALTER TABLE entries
  ALTER COLUMN amount TYPE numeric(18, 4)
  USING amount::numeric(18, 4);

ALTER TABLE transfers
  ALTER COLUMN amount TYPE numeric(18, 4)
  USING amount::numeric(18, 4);

ALTER TABLE accounts
  ADD CONSTRAINT accounts_balance_nonnegative CHECK (balance >= 0);

ALTER TABLE transfers
  ADD CONSTRAINT transfers_amount_positive CHECK (amount > 0);
