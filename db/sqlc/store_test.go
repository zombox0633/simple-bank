package db

import (
	"context"
	"testing"
	"time"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/require"

	"simplebank/internal/common"
)

func TestTransferTx(t *testing.T) {
	// ใช้ Store เพราะ TransferTx ต้องรวมหลาย SQL query ให้ทำงานใน transaction เดียวกัน
	// ส่วนบัญชีสองตัวนี้เก็บยอดเงินก่อนโอนเอาไว้ เพื่อใช้เปรียบเทียบกับยอดหลังโอน
	store := NewStore(testDB)
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	fromAccount, err := testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      fromAccount.ID,
		Balance: mustDecimal("1000.0000"),
	})
	require.NoError(t, err)
	t.Logf("fromAccount: %+v, toAccount: %+v", fromAccount, toAccount)

	// จำลองการโอนเงินพร้อมกัน 5 ครั้ง ครั้งละ 10.1234 หน่วย
	// การรันพร้อมกันช่วยทดสอบว่า database transaction ป้องกันยอดเงินสูญหาย
	// หรือถูกเขียนทับเมื่อหลาย goroutine อัปเดตบัญชีเดียวกันหรือไม่
	const n = 5
	amount := mustDecimal("10.1234")

	// ใช้ buffered channels ขนาด n เพื่อรับ error และผลลัพธ์จากทุก goroutine
	// โดย goroutine สามารถส่งผลกลับมาได้โดยไม่ต้องรอให้ test เริ่มอ่านทันที
	errs := make(chan error, n)
	results := make(chan *TransferTxResult, n)

	// เริ่ม transaction n ชุดพร้อมกัน ทุกชุดโอนจากบัญชีเดียวกันไปบัญชีเดียวกัน
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	// executed เก็บลำดับยอดเงินที่เกิดขึ้นหลังแต่ละ transaction สำเร็จ
	// เช่น โอนครั้งละ 10 จะต้องเห็นส่วนต่าง 10, 20, 30, 40 และ 50 อย่างละหนึ่งครั้ง
	// ถ้าลำดับซ้ำ แปลว่าอาจมี concurrent update ที่เขียนทับผลของกันและกัน
	executed := make(map[int]bool)
	for i := 0; i < n; i++ {
		// ทุก transaction ต้องสำเร็จ และต้องคืนผลลัพธ์กลับมา
		require.NoError(t, <-errs)
		result := <-results
		require.NotNil(t, result)

		// ตรวจว่า transfer record บันทึกต้นทาง ปลายทาง และจำนวนเงินถูกต้อง
		transfer := result.Transfer
		require.NotZero(t, transfer.ID)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		requireDecimalEqual(t, amount, transfer.Amount)

		// อ่าน transfer จากฐานข้อมูลอีกครั้ง เพื่อยืนยันว่า record ถูก commit จริง
		dbTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		require.Equal(t, transfer, dbTransfer)

		// บัญชีต้นทางต้องมี entry ติดลบ เพราะเป็นเงินออก
		fromEntry := result.FromEntry
		require.NotZero(t, fromEntry.ID)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		requireDecimalEqual(t, amount.Neg(), fromEntry.Amount)

		// บัญชีปลายทางต้องมี entry เป็นบวก เพราะเป็นเงินเข้า
		toEntry := result.ToEntry
		require.NotZero(t, toEntry.ID)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		requireDecimalEqual(t, amount, toEntry.Amount)

		// เงินที่หายจากต้นทางต้องเท่ากับเงินที่เพิ่มให้ปลายทางเสมอ
		// และส่วนต่างต้องเพิ่มทีละ amount โดยไม่มีค่าค้างระหว่างทาง
		fromDiff := mustSub(t, fromAccount.Balance, result.FromAccount.Balance)
		toDiff := mustSub(t, result.ToAccount.Balance, toAccount.Balance)
		t.Logf("fromDiff: %s, toDiff: %s", fromDiff, toDiff)

		requireDecimalEqual(t, fromDiff, toDiff)
		require.True(t, fromDiff.IsPos())

		// แปลงส่วนต่างเป็นลำดับ transaction แล้วตรวจว่าครบและไม่ซ้ำกัน
		// ไม่จำเป็นต้องเรียง 1 ถึง n เพราะ goroutine อาจทำเสร็จคนละลำดับกับตอนเริ่ม
		stepValue, err := fromDiff.QuoExact(amount, 0)
		require.NoError(t, err)
		stepNumber, _, ok := stepValue.Int64(0)
		require.True(t, ok)
		step := int(stepNumber)
		require.GreaterOrEqual(t, step, 1)
		require.LessOrEqual(t, step, n)
		require.NotContains(t, executed, step)
		executed[step] = true
	}

	// อ่านยอดล่าสุดจากฐานข้อมูลหลัง transaction ทั้งหมดจบ
	// ต้นทางต้องลดรวม n*amount และปลายทางต้องเพิ่มรวมเท่ากัน
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)
	updatedToAccount, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	t.Logf("updatedFromAccount: %+v, updatedToAccount: %+v", updatedFromAccount, updatedToAccount)
	totalAmount := mustMul(t, amount, decimal.MustNew(n, 0))
	requireDecimalEqual(t, mustSub(t, fromAccount.Balance, totalAmount), updatedFromAccount.Balance)
	requireDecimalEqual(t, mustAdd(t, toAccount.Balance, totalAmount), updatedToAccount.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	var err error
	account1, err = testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID: account1.ID, Balance: mustDecimal("1000.0000"),
	})
	require.NoError(t, err)
	account2, err = testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID: account2.ID, Balance: mustDecimal("1000.0000"),
	})
	require.NoError(t, err)

	// ครึ่งหนึ่งโอน 1 -> 2 และอีกครึ่งโอน 2 -> 1 พร้อมกัน
	// ถ้าแต่ละ transaction lock ต้นทางก่อน จะเกิดวงจรรอ lock ได้
	const n = 10
	amount := mustDecimal("10.1234")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := make(chan struct{})
	errs := make(chan error, n)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID, toAccountID = toAccountID, fromAccountID
		}

		go func() {
			<-start
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	// ปล่อยทุก goroutine ให้เริ่มใกล้เคียงกันเพื่อเพิ่มโอกาสแย่ง row lock
	close(start)
	for i := 0; i < n; i++ {
		require.NoError(t, <-errs)
	}

	// มีการโอนเข้าและออกเท่ากัน ยอดสุทธิจึงต้องไม่เปลี่ยน
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	requireDecimalEqual(t, account1.Balance, updatedAccount1.Balance)
	requireDecimalEqual(t, account2.Balance, updatedAccount2.Balance)
}

func TestTransferTxRejectsInvalidAmount(t *testing.T) {
	store := &Store{}

	_, err := store.TransferTx(context.Background(), TransferTxParams{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        decimal.Zero,
	})
	require.ErrorIs(t, err, common.ErrInvalidAmount)

	_, err = store.TransferTx(context.Background(), TransferTxParams{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        mustDecimal("1.23456"),
	})
	require.ErrorIs(t, err, common.ErrInvalidAmount)

	_, err = store.TransferTx(context.Background(), TransferTxParams{
		FromAccountID: 1,
		ToAccountID:   2,
		Amount:        mustAdd(t, common.MaxMoneyAmount, mustDecimal("0.0001")),
	})
	require.ErrorIs(t, err, common.ErrInvalidAmount)
}

func TestTransferTxRejectsSameAccount(t *testing.T) {
	store := &Store{}

	_, err := store.TransferTx(context.Background(), TransferTxParams{
		FromAccountID: 42,
		ToAccountID:   42,
		Amount:        mustDecimal("1.0000"),
	})
	require.ErrorIs(t, err, common.ErrSameAccount)
}

func TestTransferTxRejectsBalanceOverflow(t *testing.T) {
	store := NewStore(testDB)
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	fromAccount, err := testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID: fromAccount.ID, Balance: mustDecimal("100.0000"),
	})
	require.NoError(t, err)
	toAccount, err = testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID: toAccount.ID, Balance: common.MaxMoneyAmount,
	})
	require.NoError(t, err)

	_, err = store.TransferTx(context.Background(), TransferTxParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        mustDecimal("1.0000"),
	})
	require.ErrorIs(t, err, common.ErrBalanceLimitExceeded)

	updatedFromAccount, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)
	updatedToAccount, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)
	requireDecimalEqual(t, fromAccount.Balance, updatedFromAccount.Balance)
	requireDecimalEqual(t, toAccount.Balance, updatedToAccount.Balance)
}
