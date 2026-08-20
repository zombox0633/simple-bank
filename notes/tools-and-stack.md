---
title: Tools & Stack
type: reference
tags: [project/simple-bank, stack, tools]
created: 2026-07-31
updated: 2026-08-03
---

# 🧰 Tools & Stack — "ลงอะไร ใช้อะไร"

กลับไปหน้ารวม → [[simple-bank]] · คำสั่งรันทั้งหมด → [[commands]]

## ✅ ติดตั้งแล้ว / ใช้อยู่
| ตัว | เวอร์ชัน | ลงยังไง | ใช้ทำอะไร |
|---|---|---|---|
| **Go** | 1.26.5 | มีอยู่แล้ว | ภาษาหลักของ backend |
| **sqlc** | v1.31.1 | `brew install sqlc` | gen Go code จาก SQL ([[phase-2-sqlc-setup]]) |
| **golang-migrate** | (brew) | `brew install golang-migrate` | DB migration up/down |
| **Docker** | 28.1.1 | มีอยู่แล้ว | รัน PostgreSQL แบบ container |
| **PostgreSQL** | 18-alpine | `docker-compose.yml` | database จริง (พอร์ต **5433**) |
| **Gin** | v1.12.0 | `go get github.com/gin-gonic/gin` | HTTP web framework |
| **Air** | latest | `go install github.com/air-verse/air@latest` | live-reload ตอน dev |
| **pgx/v5** | v5.10.0 | `go get github.com/jackc/pgx/v5` | Native PostgreSQL driver และ connection pool |
| **govalues/decimal** | v0.1.36 | `go get github.com/govalues/decimal` | คำนวณเงินแบบ decimal exact ภายในขอบเขต 19 digits |
| **pgx-govalues-decimal** | v0.1.0 | `go get github.com/ColeBurch/pgx-govalues-decimal` | เชื่อม `decimal.Decimal` กับ PostgreSQL `NUMERIC` |
| **godotenv** | v1.5.1 | `go get github.com/joho/godotenv` | โหลด `.env` เข้า env ([[phase-4-db-connect-config-testing]]) |
| **testify** | (go get) | `go get github.com/stretchr/testify` | assertion ใน test (`require`) |

> [!warning] พอร์ต Postgres
> เครื่องมี `postgres_db` (postgres:15) จองพอร์ต 5432 อยู่แล้ว → simple-bank ใช้ **5433** แทน

## ⬜ ยังไม่ได้ลง (แผนอนาคต)
| ตัว | ลงยังไง (ตอนถึงเวลา) | ใช้ทำอะไร |
|---|---|---|
| **x/crypto/bcrypt** | `go get golang.org/x/crypto/bcrypt` | hash password (auth phase) — มีเป็น indirect อยู่แล้ว |
| **paseto / JWT** | `go get ...` | auth token |

## 📁 ไฟล์สำคัญในโปรเจกต์
- `db/migration/*.sql` — schema (source of truth) ← sqlc + migrate อ่านตรงนี้
- `db/query/*.sql` — SQL ที่เราเขียน (account, users)
- `db/sqlc/*.go` — โค้ด Go ที่ sqlc gen (อย่าแก้มือ) + `*_test.go` (db test)
- `sqlc.yaml` — generate แบบ `pgx/v5` และ override money เป็น `govalues/decimal.Decimal` รวมถึง `users.password` → `json:"-"`
- `main.go` — ต่อ DB (pgx) + Gin server
- `.env` / `.env.example` — config (DSN, server addr) · `.env` โดน gitignore
- `docker-compose.yml` · `Makefile` (ดู [[commands]]) · `.air.toml`

## 🔑 คำสั่ง
คำสั่งรันทั้งหมด (docker · migrate · reset · sqlc · **test** · server) → [[commands]]
