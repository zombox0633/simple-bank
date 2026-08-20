---
title: Phase 2 — sqlc Setup
type: phase
tags: [project/simple-bank, golang, sqlc, codegen]
created: 2026-07-31
updated: 2026-07-31
---

# Phase 2 — ติดตั้ง + ใช้งาน sqlc กับ Go

กลับไปหน้ารวม → [[simple-bank]]

## sqlc คืออะไร
เครื่องมือ gen **type-safe Go code จาก SQL ที่เราเขียนเอง** — เราเขียน query เป็น SQL ปกติ
แล้ว sqlc สร้าง struct + function Go ให้ (ไม่ใช่ ORM, ไม่มี magic runtime)

## ทำอะไรไป (ตามลำดับ)
1. สร้าง Go module
   ```bash
   go mod init simplebank
   ```
2. ติดตั้ง sqlc ผ่าน Homebrew
   ```bash
   brew install sqlc      # ได้ v1.31.1
   ```
3. สร้างไฟล์ config `sqlc.yaml` — ชี้ไปที่ schema + โฟลเดอร์ query
4. เขียน query ตัวอย่างที่ `db/query/account.sql` (CRUD ของ accounts)
5. gen โค้ด
   ```bash
   sqlc generate
   ```
6. เช็คว่า build ผ่าน
   ```bash
   go build ./...         # BUILD OK
   ```

## โครงไฟล์ที่ได้
```
simple-bank/
├── go.mod                 # module simplebank
├── sqlc.yaml              # config ของ sqlc
└── db/
    ├── migration/         # schema (input ให้ sqlc) — ดู [[phase-3-docker-migrate-server]]
    │   └── 000001_init_schema.up.sql
    ├── query/
    │   └── account.sql    # ← เราเขียน query ตรงนี้
    └── sqlc/              # ← sqlc gen ให้ (ห้ามแก้มือ)
        ├── db.go          # DBTX interface + New()
        ├── models.go      # struct: Account, Entry, Transfer, User
        └── account.sql.go # func: CreateAccount, GetAccount, ...
```
> หมายเหตุ: ตอน phase 2 schema ยังชื่อ `simple-bank.sql` ที่ root — ต่อมาย้ายเข้า `db/migration/`

## sqlc.yaml (ที่ใช้จริง)
```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "db/migration"
    queries: "db/query"
    gen:
      go:
        package: "db"
        out: "db/sqlc"
        emit_json_tags: true
```

## จุดที่น่าจดจำ
- ตั้ง `sql_package: pgx/v5` เพื่อ generate native pgx API และใช้ `pgxpool` เชื่อม PostgreSQL
- money model ใช้ `govalues/decimal` ตาม type override เพื่อเก็บ `NUMERIC(18,4)` แบบ exact
- ใน query ใช้ comment สั่งงาน sqlc:
  - `-- name: CreateAccount :one` → คืน 1 row
  - `:many` → คืนหลาย row, `:exec` → ไม่คืนค่า
- `sqlc.arg(amount)` ใช้ตั้งชื่อ parameter ให้อ่านง่าย (แทน `$1`)
- ทุกครั้งที่แก้ `.sql` ต้องรัน `sqlc generate` ใหม่
