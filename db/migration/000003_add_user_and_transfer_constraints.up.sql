ALTER TABLE users
  ADD COLUMN password_changed_at timestamptz NOT NULL
  DEFAULT '0001-01-01 00:00:00Z';

ALTER TABLE transfers
  ADD CONSTRAINT transfers_different_accounts_check
  CHECK (from_account_id <> to_account_id);
