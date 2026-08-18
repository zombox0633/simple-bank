-- name: CreateAccount :one
-- สร้างบัญชีใหม่และคืนแถวที่เพิ่งสร้าง รวมถึง id/timestamp ที่ PostgreSQL สร้างให้
INSERT INTO accounts (owner, balance, currency)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateAccount :one
-- แทนที่ balance ด้วยยอดใหม่โดยตรง เหมาะกับงานที่รู้ยอดปลายทางแน่นอน
-- สำหรับการเพิ่ม/ลดจากยอดเดิม โดยเฉพาะ money transfer ให้ใช้ AddAccountBalance
UPDATE accounts
SET balance = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
-- เพิ่มหรือลดยอดแบบ atomic ภายใน PostgreSQL: amount บวกคือเงินเข้า ลบคือเงินออก
-- การใช้ UPDATE คำสั่งเดียวช่วยไม่ให้ concurrent transaction อ่านยอดเดิมแล้วเขียนทับกัน
-- sqlc.arg กำหนดชื่อ field ใน AddAccountBalanceParams ให้เป็น Amount และ ID อย่างชัดเจน
UPDATE accounts
SET balance = balance + sqlc.arg(amount),
    updated_at = now()
WHERE id = sqlc.arg(id)
    AND balance + sqlc.arg(amount) >= 0
RETURNING *;

-- name: DeleteAccount :exec
-- ลบบัญชีด้วย primary key; foreign key constraints จะป้องกันการลบที่ทำให้ข้อมูลอ้างอิงเสีย
DELETE FROM accounts
WHERE id = $1;
