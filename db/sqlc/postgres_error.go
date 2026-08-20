package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// PostgreSQL SQLSTATE codes that can be produced by the current schema and
// transaction flow. Keep the codes here so callers do not repeat magic strings.
const (
	SQLStateNumericValueOutOfRange = "22003"
	SQLStateNotNullViolation       = "23502"
	SQLStateForeignKeyViolation    = "23503"
	SQLStateUniqueViolation        = "23505"
	SQLStateCheckViolation         = "23514"
	SQLStateSerializationFailure   = "40001"
	SQLStateDeadlockDetected       = "40P01"
)

// SQLState returns the PostgreSQL error code even when the PgError is wrapped.
// An empty string means the error did not come from PostgreSQL.
func SQLState(err error) string {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return ""
	}
	return pgErr.Code
}

// IsRetryableTransactionError identifies transaction failures that are safe to
// retry from the beginning. It does not retry automatically because callers
// must decide the retry count, backoff, and whether the operation is idempotent.
func IsRetryableTransactionError(err error) bool {
	switch SQLState(err) {
	case SQLStateSerializationFailure, SQLStateDeadlockDetected:
		return true
	default:
		return false
	}
}
