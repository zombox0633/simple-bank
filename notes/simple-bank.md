---
title: Simple Bank
type: MOC
tags: [project/simple-bank, golang, backend]
created: 2026-07-31
updated: 2026-07-31
---

# 🏦 Simple Bank — Map of Content

โปรเจกต์ backend ระบบธนาคารอย่างง่าย (Go + PostgreSQL) — โน้ตรวมทุกอย่างที่ทำ

> ⚡ อยากได้คำสั่งรันด่วน (docker · migrate · reset · sqlc · server) → [[commands]]

## 📌 สถานะปัจจุบัน
- ✅ ออกแบบ database schema เสร็จ → ดู [[phase-1-database-schema]]
- ✅ ติดตั้ง + ตั้งค่า `sqlc` ใช้งานกับ Go ได้แล้ว → ดู [[phase-2-sqlc-setup]]
- ✅ Docker Postgres + migration + Gin/Air ใช้งานได้ → ดู [[phase-3-docker-migrate-server]]
- ✅ ต่อ DB (pgx) + .env config + DB testing (10 tests, 85%) → ดู [[phase-4-db-connect-config-testing]]
- ⬜ ยังไม่ได้: Store/TransferTx, query entry/transfer, API routes, bcrypt

## 🗂️ โน้ตในโปรเจกต์
- [[phase-1-database-schema]] — ออกแบบตาราง 4 ตัว + ความสัมพันธ์
- [[phase-2-sqlc-setup]] — ติดตั้ง sqlc, gen โค้ด Go จาก SQL
- [[phase-3-docker-migrate-server]] — Docker Postgres, golang-migrate, Gin, Air
- [[phase-4-db-connect-config-testing]] — ต่อ DB (pgx), .env, DB testing, ปิด password
- [[commands]] — ⚡ cheatsheet รวมคำสั่งทั้งหมด (docker · migrate · reset · sqlc · test · server)
- [[tools-and-stack]] — สรุปว่า "ลงอะไร ใช้อะไร" ทั้งหมด

## 🧭 แผนถัดไป (roadmap)
- [x] ต่อ DB จาก `main.go` จริง (pgx) ✅
- [x] unit test ของ db layer — account + user (10 tests, coverage 85%) ✅
- [ ] query ที่เหลือ: `entry.sql`, `transfer.sql` แล้ว `make sqlc`
- [ ] `util/password.go` — bcrypt (`HashPassword`/`CheckPassword`)
- [ ] Store layer + `TransferTx` (db transaction — โอนเงิน atomic)
- [ ] HTTP API routes (Gin) + validation (`binding:"required,email"`)
