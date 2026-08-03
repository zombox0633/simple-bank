---
title: Phase 1 — Database Schema
type: phase
tags: [project/simple-bank, database, postgresql, schema]
created: 2026-07-31
updated: 2026-07-31
---

# Phase 1 — ออกแบบ Database Schema

กลับไปหน้ารวม → [[simple-bank]]

![[simple-bank.png]]

## ทำอะไรไป
ออกแบบ schema บน **dbdiagram.io** แล้ว export ออกมาเป็น:
- `simple-bank.sql` — โครงตาราง (PostgreSQL) → ภายหลังย้ายไป `db/migration/000001_init_schema.up.sql` ใน [[phase-3-docker-migrate-server]]
- `simple-bank.png` — ER diagram

## ตารางทั้งหมด (4 ตาราง)
| ตาราง | หน้าที่ | จุดสำคัญ |
|---|---|---|
| `users` | ผู้ใช้ | PK = `username`, `email` UNIQUE |
| `accounts` | บัญชีเงิน | มี `owner`, `balance`, `currency` |
| `entries` | รายการเดินบัญชี | `amount` เป็นบวก/ลบได้ |
| `transfers` | การโอนเงิน | `amount` ต้องเป็นบวก |

## ความสัมพันธ์ (Foreign Keys)
- `accounts.owner` → `users.username`
- `entries.account_id` → `accounts.id`
- `transfers.from_account_id` / `to_account_id` → `accounts.id`
- ทุก FK เป็น `DEFERRABLE INITIALLY IMMEDIATE`

## Constraints / Index ที่ตั้งไว้
- `UNIQUE (owner, currency)` บน `accounts` → 1 คนมีได้บัญชีเดียวต่อสกุลเงิน
- Index บน `entries(account_id)`, `transfers(to_account_id)`, `transfers(from_account_id, to_account_id)`

## หมายเหตุ
- ตอน phase นี้ยังเป็นแค่ไฟล์ `.sql` เฉย ๆ (ต่อมาถูกรันเข้า DB จริงด้วย `migrate` ใน [[phase-3-docker-migrate-server]])
- schema นี้ถูกเอาไปให้ [[phase-2-sqlc-setup|sqlc]] อ่านเพื่อ gen โค้ด Go
