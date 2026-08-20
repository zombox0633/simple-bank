ALTER TABLE transfers
  DROP CONSTRAINT IF EXISTS transfers_different_accounts_check;

ALTER TABLE users
  DROP COLUMN IF EXISTS password_changed_at;
