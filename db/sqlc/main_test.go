package db

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx แบบ database/sql
)

const (
	dbDriver = "pgx" // ชื่อที่ pgx/stdlib register
	dbSource = "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("failed to connect to db 😿 : ", err)
	}
	if err := testDB.Ping(); err != nil { // Open ไม่ได้ต่อจริง → Ping ให้ fail เร็วถ้า DB ไม่ขึ้น
		log.Fatal("failed to ping db 😹 : ", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
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

func randomMoney() int64 {
	return randomInt(0, 1000)
}

func randomCurrency() string {
	currencies := []string{"USD", "EUR", "THB"}
	return currencies[rand.Intn(len(currencies))]
}
