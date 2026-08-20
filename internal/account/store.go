package account

import (
	"context"

	db "simplebank/db/sqlc"
)

//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=store.go -destination=mock_store_test.go -package=account

// Store ระบุเฉพาะ database operations ที่ account feature ต้องใช้
// *db.Store ทำตาม interface นี้อยู่แล้ว และภายหลังสามารถใช้ mock ใน API tests ได้
type Store interface {
	CreateAccount(context.Context, db.CreateAccountParams) (db.Account, error)
	GetAccount(context.Context, int64) (db.Account, error)
	ListAccounts(context.Context, db.ListAccountsParams) ([]db.Account, error)
}
