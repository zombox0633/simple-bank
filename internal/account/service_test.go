package account

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"

	db "simplebank/db/sqlc"
	"simplebank/internal/common"
)

type stubStore struct {
	createAccount func(context.Context, db.CreateAccountParams) (db.Account, error)
	getAccount    func(context.Context, int64) (db.Account, error)
	listAccounts  func(context.Context, db.ListAccountsParams) ([]db.Account, error)
}

func (store *stubStore) CreateAccount(
	ctx context.Context,
	arg db.CreateAccountParams,
) (db.Account, error) {
	return store.createAccount(ctx, arg)
}

func (store *stubStore) GetAccount(ctx context.Context, id int64) (db.Account, error) {
	return store.getAccount(ctx, id)
}

func (store *stubStore) ListAccounts(
	ctx context.Context,
	arg db.ListAccountsParams,
) ([]db.Account, error) {
	return store.listAccounts(ctx, arg)
}

func TestServiceCreateAccount(t *testing.T) {
	want := db.Account{ID: 42, Owner: "alice", Currency: common.CurrencyTHB}
	store := &stubStore{
		createAccount: func(_ context.Context, arg db.CreateAccountParams) (db.Account, error) {
			require.Equal(t, "alice", arg.Owner)
			require.Equal(t, common.CurrencyTHB, arg.Currency)
			require.True(t, arg.Balance.IsZero())
			return want, nil
		},
	}

	got, err := NewService(store).CreateAccount(context.Background(), "alice", common.CurrencyTHB)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestServiceCreateAccountRejectsUnsupportedCurrency(t *testing.T) {
	storeCalled := false
	store := &stubStore{
		createAccount: func(_ context.Context, _ db.CreateAccountParams) (db.Account, error) {
			storeCalled = true
			return db.Account{}, nil
		},
	}

	_, err := NewService(store).CreateAccount(context.Background(), "alice", "BTC")

	require.ErrorIs(t, err, common.ErrUnsupportedCurrency)
	require.False(t, storeCalled)
}

func TestServiceGetAccountMapsNotFound(t *testing.T) {
	store := &stubStore{
		getAccount: func(_ context.Context, _ int64) (db.Account, error) {
			return db.Account{}, sql.ErrNoRows
		},
	}

	_, err := NewService(store).GetAccount(context.Background(), 42)

	require.ErrorIs(t, err, common.ErrNotFound)
}

func TestServiceListAccountsBuildsPagination(t *testing.T) {
	want := []db.Account{{ID: 42, Owner: "alice", Currency: common.CurrencyTHB}}
	store := &stubStore{
		listAccounts: func(_ context.Context, arg db.ListAccountsParams) ([]db.Account, error) {
			require.Equal(t, "alice", arg.Owner)
			require.Equal(t, int32(25), arg.Limit)
			require.Equal(t, int32(50), arg.Offset)
			return want, nil
		},
	}

	got, err := NewService(store).ListAccounts(context.Background(), "alice", 3, 25)

	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestMapCreateAccountError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want error
	}{
		{
			name: "duplicate account",
			err:  &pgconn.PgError{Code: "23505"},
			want: common.ErrConflict,
		},
		{
			name: "owner not found",
			err:  &pgconn.PgError{Code: "23503"},
			want: common.ErrInvalidReference,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.ErrorIs(t, mapCreateAccountError(test.err), test.want)
		})
	}

	unexpected := errors.New("connection closed")
	require.ErrorIs(t, mapCreateAccountError(unexpected), unexpected)
}
