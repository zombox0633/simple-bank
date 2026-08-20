package transfer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"

	db "simplebank/db/sqlc"
	"simplebank/internal/common"
)

type Service struct {
	store    Store
	accounts AccountReader
}

func NewService(store Store, accounts AccountReader) *Service {
	return &Service{store: store, accounts: accounts}
}

func (service *Service) CreateTransfer(
	ctx context.Context,
	fromAccountID int64,
	toAccountID int64,
	amount decimal.Decimal,
	currency string,
) (*db.TransferTxResult, error) {
	if fromAccountID == toAccountID {
		return nil, common.ErrSameAccount
	}
	if !common.IsSupportedCurrency(currency) {
		return nil, common.ErrUnsupportedCurrency
	}
	if !common.IsValidMoneyAmount(amount) {
		return nil, common.ErrInvalidAmount
	}

	fromAccount, err := service.accounts.GetAccount(ctx, fromAccountID)
	if err != nil {
		return nil, err
	}
	if fromAccount.Currency != currency {
		return nil, fmt.Errorf("%w: account %d uses %s", common.ErrCurrencyMismatch, fromAccount.ID, fromAccount.Currency)
	}

	toAccount, err := service.accounts.GetAccount(ctx, toAccountID)
	if err != nil {
		return nil, err
	}
	if toAccount.Currency != currency {
		return nil, fmt.Errorf("%w: account %d uses %s", common.ErrCurrencyMismatch, toAccount.ID, toAccount.Currency)
	}
	if fromAccount.Balance.LessThan(amount) {
		return nil, common.ErrInsufficientBalance
	}

	result, err := service.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, common.ErrInsufficientBalance
		case errors.Is(err, db.ErrInvalidTransferAmount), errors.Is(err, db.ErrInvalidTransferAmountPrecision):
			return nil, common.ErrInvalidAmount
		default:
			return nil, fmt.Errorf("transfer money: %w", err)
		}
	}

	return result, nil
}
