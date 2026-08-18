package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store รวม generated Queries สำหรับคำสั่งเดี่ยว และเก็บ connection pool
// ไว้เริ่ม transaction ที่ต้องรันหลายคำสั่งเป็นหน่วยเดียวกัน
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore สร้าง Store ที่ใช้ database connection pool เดียวกันทั้ง query ปกติ
// และ transaction
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx ครอบ callback ด้วย database transaction:
// callback สำเร็จจึง commit และถ้า callback ล้มเหลวจะ rollback ทุกคำสั่งที่ผ่านมา
func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Queries ตัวนี้ผูกกับ tx ทำให้ทุก query ใน callback ใช้ transaction เดียวกัน
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rollback error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams คือข้อมูลที่จำเป็นสำหรับการโอนเงินหนึ่งครั้ง
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult รวมแถวทั้งหมดที่ถูกสร้างหรืออัปเดตภายใน transaction
// เพื่อให้ caller ตรวจสอบทั้ง transfer, ledger entries และยอดบัญชีหลังโอนได้
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx บันทึกการโอน เงินเข้า/ออกใน ledger และปรับยอดสองบัญชีแบบ atomic
// ถ้าขั้นตอนใดผิดพลาด execTx จะ rollback ทุกขั้นตอน จึงไม่มีข้อมูลค้างเพียงบางส่วน
func (s *Store) TransferTx(ctx context.Context, params TransferTxParams) (*TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		// 1) เก็บหัวรายการโอน เพื่อรู้ว่าเงินย้ายจากบัญชีใดไปบัญชีใดเท่าไร
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: params.FromAccountID,
			ToAccountID:   params.ToAccountID,
			Amount:        params.Amount,
		})

		if err != nil {
			return err
		}

		// 2) เก็บ ledger entry ของต้นทางเป็นค่าลบและปลายทางเป็นค่าบวก
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.FromAccountID,
			Amount:    -params.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: params.ToAccountID,
			Amount:    params.Amount,
		})
		if err != nil {
			return err
		}

		// 3) AddAccountBalance เป็น atomic UPDATE และจะล็อกแถวบัญชีที่กำลังแก้
		// ทุก transaction จึงต้องอัปเดตบัญชีตามลำดับ ID เดียวกัน เพื่อไม่ให้
		// การโอนสวนทางกันถือ lock คนละบัญชีแล้วรอกันเป็นวงกลมจนเกิด deadlock
		if params.FromAccountID < params.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx, q,
				params.FromAccountID, -params.Amount,
				params.ToAccountID, params.Amount,
			)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(
				ctx, q,
				params.ToAccountID, params.Amount,
				params.FromAccountID, -params.Amount,
			)
		}
		return err
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// addMoney ปรับยอดบัญชีแรกให้เสร็จก่อนบัญชีที่สองเสมอ โดย caller จะส่งบัญชี
// ตามลำดับ ID จากน้อยไปมากเพื่อให้ทุก transaction ล็อกแถวในลำดับเดียวกัน
// balanceChange เป็นค่าที่เพิ่มจากยอดปัจจุบัน: ค่าบวกคือเงินเข้า ค่าลบคือเงินออก
func addMoney(
	ctx context.Context,
	q *Queries,
	firstAccountID int64,
	firstBalanceChange int64,
	secondAccountID int64,
	secondBalanceChange int64,
) (firstUpdatedAccount Account, secondUpdatedAccount Account, err error) {
	firstUpdatedAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     firstAccountID,
		Amount: firstBalanceChange,
	})
	if err != nil {
		return
	}

	secondUpdatedAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     secondAccountID,
		Amount: secondBalanceChange,
	})
	return
}
