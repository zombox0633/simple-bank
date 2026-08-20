package account

import (
	"context"

	db "simplebank/db/sqlc"
)

// Store ระบุเฉพาะ database operations ที่ account feature ต้องใช้
// *db.Store ทำตาม interface นี้อยู่แล้ว และภายหลังสามารถใช้ mock ใน API tests ได้
type Store interface {
	CreateAccount(context.Context, db.CreateAccountParams) (db.Account, error)
	GetAccount(context.Context, int64) (db.Account, error)
	ListAccounts(context.Context, db.ListAccountsParams) ([]db.Account, error)
}
