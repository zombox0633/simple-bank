package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    createRandomUser(t).Username,
		Balance:  randomMoney(),
		Currency: randomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	requireDecimalEqual(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	requireDecimalEqual(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: randomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	requireDecimalEqual(t, arg.Balance, account2.Balance) // balance ต้องเปลี่ยนเป็นค่าใหม่
	require.Equal(t, account1.Currency, account2.Currency)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	// ดึงมาอีกที ต้องไม่เจอแล้ว → pgx.ErrNoRows
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	owner := createRandomUser(t).Username

	// 1 owner มีได้หลาย account ถ้าคนละ currency (เพราะติด unique(owner, currency))
	for _, cur := range testCurrencies {
		arg := CreateAccountParams{Owner: owner, Balance: randomMoney(), Currency: cur}
		_, err := testQueries.CreateAccount(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListAccountsParams{Owner: owner, Limit: 2, Offset: 1} // ข้าม 1 ตัวแรก
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, len(testCurrencies)-1) // 3 - 1 = 2

	for _, a := range accounts {
		require.NotEmpty(t, a)
		require.Equal(t, owner, a.Owner)
	}
}
