package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/govalues/decimal"
	"github.com/jackc/pgx/v5"

	db "simplebank/db/sqlc"
	"simplebank/internal/common"
)

// Service contains account business rules and keeps Gin handlers independent
// from sqlc parameter types.
type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (service *Service) CreateAccount(
	ctx context.Context,
	owner string,
	currency string,
) (db.Account, error) {
	if !common.IsSupportedCurrency(currency) {
		return db.Account{}, common.ErrUnsupportedCurrency
	}

	account, err := service.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    owner,
		Balance:  decimal.Zero,
		Currency: currency,
	})
	if err != nil {
		return db.Account{}, mapCreateAccountError(err)
	}

	return account, nil
}

func (service *Service) GetAccount(ctx context.Context, id int64) (db.Account, error) {
	account, err := service.store.GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Account{}, common.ErrNotFound
		}
		return db.Account{}, fmt.Errorf("get account: %w", err)
	}

	return account, nil
}

func (service *Service) ListAccounts(
	ctx context.Context,
	owner string,
	pageID int32,
	pageSize int32,
) ([]db.Account, error) {
	accounts, err := service.store.ListAccounts(ctx, db.ListAccountsParams{
		Owner:  owner,
		Limit:  pageSize,
		Offset: (pageID - 1) * pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}

	return accounts, nil
}

func mapCreateAccountError(err error) error {
	switch db.SQLState(err) {
	case db.SQLStateUniqueViolation: // an owner already has this currency
		return common.ErrConflict
	case db.SQLStateForeignKeyViolation: // owner does not exist
		return common.ErrInvalidReference
	}

	return fmt.Errorf("create account: %w", err)
}
