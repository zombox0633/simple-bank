package db

import (
	"context"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	pgxdecimal "github.com/ColeBurch/pgx-govalues-decimal"
	"github.com/govalues/decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"simplebank/internal/common"
)

const (
	defaultDBSource = "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable"
)

var (
	testQueries    *Queries
	testDB         *pgxpool.Pool
	testCurrencies = [...]string{
		common.CurrencyUSD,
		common.CurrencyEUR,
		common.CurrencyTHB,
	}
)

func TestMain(m *testing.M) {
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = defaultDBSource
	}

	poolConfig, err := pgxpool.ParseConfig(dbSource)
	if err != nil {
		log.Fatal("failed to parse db config 😿 : ", err)
	}
	poolConfig.AfterConnect = func(_ context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}

	startupCtx, cancelStartup := context.WithTimeout(context.Background(), 10*time.Second)
	testDB, err = pgxpool.NewWithConfig(startupCtx, poolConfig)
	if err != nil {
		log.Fatal("failed to create db pool 😿 : ", err)
	}
	if err := testDB.Ping(startupCtx); err != nil {
		log.Fatal("failed to ping db 😹 : ", err)
	}
	cancelStartup()
	testQueries = New(testDB)

	exitCode := m.Run()
	testDB.Close()
	os.Exit(exitCode)
}

// ---------- random helpers (ใช้ร่วมกันทุก test ใน package db) ----------

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomMoney() decimal.Decimal {
	// สุ่มทั้งส่วนจำนวนเต็มและเศษ 4 ตำแหน่ง เช่น 123.4567
	return decimal.MustNew(randomInt(0, 1000)*10_000+randomInt(0, 9_999), common.MoneyScale)
}

func randomPositiveMoney() decimal.Decimal {
	return decimal.MustNew(randomInt(1, 1000)*10_000+randomInt(0, 9_999), common.MoneyScale)
}

func mustDecimal(value string) decimal.Decimal {
	return decimal.MustParse(value)
}

func requireDecimalEqual(t *testing.T, expected, actual decimal.Decimal) {
	t.Helper()
	if !expected.Equal(actual) {
		t.Fatalf("decimal values differ: expected %s, actual %s", expected, actual)
	}
}

func mustAdd(t *testing.T, left, right decimal.Decimal) decimal.Decimal {
	t.Helper()
	result, err := left.AddExact(right, common.MoneyScale)
	require.NoError(t, err)
	return result
}

func mustSub(t *testing.T, left, right decimal.Decimal) decimal.Decimal {
	t.Helper()
	result, err := left.SubExact(right, common.MoneyScale)
	require.NoError(t, err)
	return result
}

func mustMul(t *testing.T, left, right decimal.Decimal) decimal.Decimal {
	t.Helper()
	result, err := left.MulExact(right, common.MoneyScale)
	require.NoError(t, err)
	return result
}

func randomCurrency() string {
	return testCurrencies[rand.Intn(len(testCurrencies))]
}
