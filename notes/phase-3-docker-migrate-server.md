---
title: Phase 3 — Docker, Migration & Server
type: phase
tags: [project/simple-bank, docker, postgresql, migration, gin, air]
created: 2026-07-31
updated: 2026-07-31
---

# Phase 3 — Docker Postgres + Migration + Gin/Air

กลับไปหน้ารวม → [[simple-bank]] · คำสั่งทั้งหมด → [[commands]]

## ทำอะไรไป
1. **ย้าย schema เข้า migration** — เดิม `simple-bank.sql` อยู่ที่ root
   ย้ายไป `db/migration/000001_init_schema.up.sql` + เขียนไฟล์คู่ `..down.sql` (DROP tables)
   แล้วแก้ `sqlc.yaml` ให้ `schema: "db/migration"` → [[phase-2-sqlc-setup]]
2. **ติดตั้ง tool เพิ่ม**
   ```bash
   brew install golang-migrate                       # migrate CLI
   go install github.com/air-verse/air@latest        # live-reload
   go get github.com/gin-gonic/gin                    # web framework
   ```
3. **Docker Postgres** — เขียน `docker-compose.yml` (postgres:18-alpine)
4. **Makefile** — รวมคำสั่งที่ใช้บ่อยไว้ที่เดียว
5. **main.go** — Gin server ขั้นต่ำ มี route `/ping`
6. **.air.toml** — config ให้ air build+reload อัตโนมัติ

## ⚠️ เรื่อง port 5432 (สำคัญ)
เครื่องมี container `postgres_db` (postgres:15) จองพอร์ต **5432** อยู่ก่อนแล้ว (คนละโปรเจกต์)
→ เลยตั้งของ simple-bank ไว้ที่ host port **5433** แทน (ใน `docker-compose.yml` + `Makefile`)
- ต่อ DB จริงที่: `localhost:5433`
- ถ้าอยากใช้ 5432 ค่อยลบ/หยุด `postgres_db` เก่าก่อน แล้วแก้พอร์ตกลับ

## ⚠️ PG18 เปลี่ยน volume path (ต้องรู้)
Postgres 18 image เก็บ data ใน subdir ตาม major version → **ห้าม** mount volume ที่ `/var/lib/postgresql/data` (container จะ exit code 1)
- ✅ ถูก: `pgdata:/var/lib/postgresql`
- ❌ ผิด (แบบ PG≤17): `pgdata:/var/lib/postgresql/data`
- อัปเกรด major version (17→18) ใช้ volume เดิมไม่ได้ ต้อง `docker compose down -v` แล้ว `migrate up` ใหม่

## วิธีใช้งาน
คำสั่งรันทั้งหมดดูที่ → [[commands]]
เริ่มเร็ว: `make postgres` → `make migrateup` → `make server` แล้วทดสอบ `curl localhost:8080/ping` → ได้ `{"message":"pong"}`

## ผลลัพธ์ที่ยืนยันแล้ว ✅
- `docker compose up` → postgres รันที่ 5433
- `migrate up` → ได้ตาราง `users, accounts, entries, transfers` (+ `schema_migrations`)
- `go build ./...` ผ่าน, `air` อยู่ใน PATH (`~/go/bin`)

## ค้างไว้ / ต่อไป
- [ ] เพิ่ม driver (`pgx` หรือ `lib/pq`) แล้วต่อ DB จาก main.go จริง
- [ ] เขียน Store layer + `TransferTx` (db transaction)
- [ ] query ที่เหลือ (entry/transfer/user) แล้ว `make sqlc`
- ดู stack ทั้งหมดที่ [[tools-and-stack]]
