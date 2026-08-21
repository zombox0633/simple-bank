package user

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"

	db "simplebank/db/sqlc"
	"simplebank/internal/common"
)

type stubStore struct {
	createUser func(context.Context, db.CreateUserParams) (db.User, error)
}

func (store *stubStore) CreateUser(
	ctx context.Context,
	arg db.CreateUserParams,
) (db.User, error) {
	return store.createUser(ctx, arg)
}

func TestServiceCreateUserHashesPassword(t *testing.T) {
	const password = "Secret123!"
	want := db.User{
		Username: "charlie_01",
		FullName: "Charlie Example",
		Email:    "charlie@example.com",
	}
	store := &stubStore{
		createUser: func(_ context.Context, arg db.CreateUserParams) (db.User, error) {
			require.Equal(t, want.Username, arg.Username)
			require.Equal(t, want.FullName, arg.FullName)
			require.Equal(t, want.Email, arg.Email)
			require.NotEqual(t, password, arg.Password)
			require.NoError(t, common.CheckPassword(password, arg.Password))

			createdUser := want
			createdUser.Password = arg.Password
			return createdUser, nil
		},
	}

	got, err := NewService(store).CreateUser(
		context.Background(),
		want.Username,
		password,
		want.FullName,
		want.Email,
	)

	require.NoError(t, err)
	require.Equal(t, want.Username, got.Username)
	require.NoError(t, common.CheckPassword(password, got.Password))
}

func TestServiceCreateUserMapsDatabaseErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want error
	}{
		{
			name: "duplicate username or email",
			err:  &pgconn.PgError{Code: db.SQLStateUniqueViolation},
			want: common.ErrConflict,
		},
		{
			name: "unexpected error",
			err:  errors.New("connection closed"),
			want: errors.New("connection closed"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := &stubStore{
				createUser: func(_ context.Context, _ db.CreateUserParams) (db.User, error) {
					return db.User{}, test.err
				},
			}

			_, err := NewService(store).CreateUser(
				context.Background(),
				"charlie",
				"Secret123!",
				"Charlie Example",
				"charlie@example.com",
			)

			if test.name == "unexpected error" {
				require.ErrorIs(t, err, test.err)
				return
			}
			require.ErrorIs(t, err, test.want)
		})
	}
}
