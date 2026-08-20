package db

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestSQLState(t *testing.T) {
	pgErr := &pgconn.PgError{Code: SQLStateUniqueViolation}
	require.Equal(t, SQLStateUniqueViolation, SQLState(fmt.Errorf("create account: %w", pgErr)))
	require.Empty(t, SQLState(errors.New("connection closed")))
}

func TestIsRetryableTransactionError(t *testing.T) {
	require.True(t, IsRetryableTransactionError(&pgconn.PgError{Code: SQLStateSerializationFailure}))
	require.True(t, IsRetryableTransactionError(&pgconn.PgError{Code: SQLStateDeadlockDetected}))
	require.False(t, IsRetryableTransactionError(&pgconn.PgError{Code: SQLStateUniqueViolation}))
	require.False(t, IsRetryableTransactionError(errors.New("connection closed")))
}
