package transfer

import (
	"context"

	db "simplebank/db/sqlc"
)

type Store interface {
	TransferTx(context.Context, db.TransferTxParams) (*db.TransferTxResult, error)
}

type AccountReader interface {
	GetAccount(context.Context, int64) (db.Account, error)
}
