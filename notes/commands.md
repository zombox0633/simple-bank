---
title: Commands Cheatsheet
type: reference
tags: [project/simple-bank, commands, cheatsheet]
created: 2026-07-31
updated: 2026-07-31
---

# ⚡ Commands Cheatsheet — simple-bank

กลับไปหน้ารวม → [[simple-bank]]

> คำสั่งที่ใช้บ่อยรวมไว้ที่เดียวตรงนี้ · ส่วนใหญ่มี `make` shortcut (ตั้งไว้ใน `Makefile`)
> คำสั่ง `migrate` แบบเต็มต้องใส่ DB_URL — export ครั้งเดียวต่อ terminal จะพิมพ์สั้นลง:
> ```bash
> export DB_URL="postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable"
> ```

---

## 🐘 Database (Docker + Postgres)
| คำสั่ง | ทำอะไร |
|---|---|
| `make postgres` &nbsp;·&nbsp; `docker compose up -d` | เปิด Postgres (รัน background) |
| `make postgres-down` &nbsp;·&nbsp; `docker compose down` | ปิด container — **ข้อมูลยังอยู่** |
| `docker compose down -v` | ปิด + **ลบข้อมูลทั้งหมด** (ลบ volume `pgdata`) |
| `docker compose ps` | ดูว่า container รันอยู่ไหม |
| `docker compose logs -f postgres` | ดู log แบบ realtime (ออก = `Ctrl+C`) |
| `docker compose restart postgres` | restart container |

## 🔀 Migration (golang-migrate)
| คำสั่ง | ทำอะไร |
|---|---|
| `make migrateup` | รัน up ที่ค้างทั้งหมด (สร้าง/อัปเดตตาราง) |
| `migrate -path db/migration -database "$DB_URL" up 1` | ทำทีละ 1 step |
| `migrate -path db/migration -database "$DB_URL" down 1` | ย้อน 1 step |
| `migrate -path db/migration -database "$DB_URL" down -all` | ย้อนทั้งหมด |
| `migrate -path db/migration -database "$DB_URL" version` | ดูว่า DB อยู่ version ไหน |
| `migrate create -ext sql -dir db/migration -seq <ชื่อ>` | สร้างไฟล์ migration ใหม่ (ได้ทั้ง .up + .down) |

> 📌 แก้ schema ครั้งหน้า **อย่าแก้ไฟล์เก่า** — สร้างอันใหม่ด้วย `migrate create ...` เสมอ

## 🔄 Reset (เลือกตามสถานการณ์)
**1) Soft reset** — ย้อนแล้วทำใหม่ (ทดสอบว่า up/down เขียนถูก, ไม่แตะ container)
```bash
migrate -path db/migration -database "$DB_URL" down -all
migrate -path db/migration -database "$DB_URL" up
```

**2) Hard reset** — ล้าง DB ทั้งลูก เริ่มใหม่หมด (⚠️ data หายเกลี้ยง)
```bash
docker compose down -v && docker compose up -d
make migrateup      # ถ้าเจอ "connection refused" รอ 2-3 วิ แล้วรันซ้ำ (postgres ยังไม่พร้อม)
```

**3) แก้ dirty state** — migration พังกลางคัน DB ติดธง dirty → migrate ไม่ยอมรันต่อ
```bash
migrate -path db/migration -database "$DB_URL" version         # ดู version ที่ค้าง (จะเห็นคำว่า dirty)
migrate -path db/migration -database "$DB_URL" force <version>  # ปลดล็อก: set version โดยไม่รัน SQL
```

## 🏗️ sqlc
| คำสั่ง | ทำอะไร |
|---|---|
| `make sqlc` &nbsp;·&nbsp; `sqlc generate` | gen โค้ด Go จาก SQL — รันทุกครั้งหลังแก้ `db/query/*.sql` |

## 🚀 Server (Gin + Air)
| คำสั่ง | ทำอะไร |
|---|---|
| `make server` &nbsp;·&nbsp; `air` | รัน server + reload อัตโนมัติ → `localhost:8080` |
| `go run .` | รันครั้งเดียว (ไม่ auto-reload) |
| `curl localhost:8080/ping` | ทดสอบ → ได้ `{"message":"pong"}` |

## 🧪 Test (testify)
| คำสั่ง | ทำอะไร |
|---|---|
| `make test` &nbsp;·&nbsp; `go test -v -cover ./...` | รัน test ทั้งหมด + coverage |
| `go test ./db/sqlc/ -v` | รันเฉพาะ package db แบบ verbose |
| `go test ./db/sqlc/ -run TestGetUser -v` | รันเฉพาะ test ตัวเดียว |
| `go test -count=1 ...` | บังคับรันใหม่ (ไม่เอา cache) |

ดู coverage เป็นสี (เขียว=ถูกเทสต์ · แดง=ไม่ถูกเทสต์):
```bash
go test -coverprofile=cover.out ./db/sqlc/ && go tool cover -html=cover.out   # เปิด browser
go tool cover -func=cover.out                                                 # ดู %ต่อ function ใน terminal
```
> ⚠️ Go cache ผล test → โค้ดไม่เปลี่ยนจะขึ้น `(cached)` ไม่รันจริง · ใส่ `-count=1` ถ้าอยากรันใหม่

## 🔍 เข้าไปดู DB ตรง ๆ (psql)
```bash
docker compose exec postgres psql -U root -d simple_bank
```
ในหน้า psql: `\dt` = ดูตารางทั้งหมด · `\d accounts` = ดูโครงตาราง accounts · `\q` = ออก

---

## 🧰 Make shortcuts ที่มีตอนนี้
`postgres` · `postgres-down` · `migrateup` · `migratedown` · `sqlc` · `test` · `server`
> `make migratedown` = ย้อน migration (จะถาม `y/N` ก่อนลบ)

## 📋 Workflows (ทำเรียงตามนี้)
**ครั้งแรก / เพิ่งเปิดโปรเจกต์**
```bash
make postgres      # 1. เปิด DB
make migrateup     # 2. สร้างตาราง
make server        # 3. รัน server
```

**dev ประจำวัน** — `make postgres` → `make server` · (แก้ query เมื่อไหร่ → `make sqlc`)

**เลิกงาน** — `make postgres-down`
