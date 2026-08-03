---
title: Phase 4 — Connect DB, Config, Testing
type: phase
tags: [project/simple-bank, golang, pgx, testing, config, security]
created: 2026-08-03
updated: 2026-08-03
---

# Phase 4 — ต่อ DB จริง + Config (.env) + DB Testing

กลับไปหน้ารวม → [[simple-bank]] · คำสั่ง → [[commands]]

## ทำอะไรไป (สรุป)
1. **ต่อ DB จริงด้วย pgx** — `main.go` เชื่อม Postgres แล้ว (เห็น `connected to db`)
2. **Config ด้วย .env** — เลิก hardcode DSN ใช้ `godotenv`
3. **user queries** — `db/query/users.sql` (Create/Get/List/Update/ChangePassword)
4. **ปิด password ไม่ให้หลุด** — 2 ชั้น
5. **DB testing (testify)** — 10 tests, coverage **85%**

## 1. ต่อ DB (pgx)
```go
// main.go
_ = godotenv.Load()
conn, _ := sql.Open("pgx", os.Getenv("DB_SOURCE"))  // driver "pgx" (ผ่าน stdlib adapter)
conn.Ping()                                          // Open ไม่ได้ต่อจริง → Ping บังคับต่อ
```
- ใช้ **pgx/v5** ผ่าน `_ "github.com/jackc/pgx/v5/stdlib"` (database/sql) — ไม่ใช้ lib/pq (เก่า)

## 2. Config (.env)
| ไฟล์ | หน้าที่ |
|---|---|
| `.env` | ค่าจริง (DB_SOURCE, SERVER_ADDRESS) — **gitignore ไว้** |
| `.env.example` | template commit ได้ |
| `Makefile` | `-include .env` → migrate ใช้ `$(DB_SOURCE)` แหล่งเดียวกับ app |
- `godotenv.Load()` ทำงานแค่ใน `main.go` (Go process) — Makefile/test ต้องอ่าน .env เอง (test hardcode DSN ไว้)

## 3. Security: ปิด password 2 ชั้น
- **query**: `ListUsers` เลือก column แบบ explicit (ตัด `password`) → ไม่โหลดมาเลย
- **struct**: `sqlc.yaml` override `users.password` → `json:"-"` → ไม่ serialize ออก JSON (แต่ Go ยังใช้ตอน login ได้)
- `GetUser` ยังต้อง return password (ตอน login เอา hash มาเทียบ) → เลยต้องพึ่ง `json:"-"`

## 4. DB Testing (testify) — coverage 85%
```
db/sqlc/
  main_test.go     # TestMain (ต่อ DB) + random helpers + createRandomUser
  account_test.go  # 5 tests + createRandomAccount
  users_test.go    # 5 tests + createRandomUser
```
แพทเทิร์น: `createRandomX()` สร้าง+assert+return → Get/Update/Delete/List เรียกใช้ซ้ำ

## ค้างไว้ / ต่อไป
- [ ] `util/password.go` — `HashPassword`/`CheckPassword` (bcrypt) แล้วใช้จริงตอน register
- [ ] **Store + TransferTx** — db transaction (โอนเงิน atomic) หัวใจของ bank
- [ ] HTTP API (Gin handlers) + validation (`binding:"required,email"`)
- ดู stack → [[tools-and-stack]]
