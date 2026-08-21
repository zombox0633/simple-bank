package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"simplebank/internal/common"
)

// createRandomUser สร้าง user สุ่ม ไว้ใช้ซ้ำในเทสต์อื่น
func createRandomUser(t *testing.T) User {
	t.Helper()

	username := "user_" + randomString(8)
	plainPassword := "Secret123!"
	hashedPassword, err := common.HashPassword(plainPassword)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username: username,
		Password: hashedPassword,
		FullName: "Test User",
		Email:    username + "@test.com", // unique เพราะ username สุ่มไม่ซ้ำ
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Password, user.Password)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotEqual(t, plainPassword, user.Password)
	require.NoError(t, common.CheckPassword(plainPassword, user.Password))
	require.True(t, user.PasswordChangedAt.Valid)
	require.True(t, user.PasswordChangedAt.Time.IsZero())
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateUserParams{
		Username: user1.Username,
		FullName: "Updated Name",
	}
	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)    // key เดิม
	require.Equal(t, arg.FullName, user2.FullName)      // ชื่อเปลี่ยนแล้ว
	require.NotEqual(t, user1.FullName, user2.FullName) // ต้องต่างจากของเดิม
}

func TestChangePassword(t *testing.T) {
	user1 := createRandomUser(t)

	newPassword := "New1!" + randomString(6)
	newPasswordHash, err := common.HashPassword(newPassword)
	require.NoError(t, err)

	err = testQueries.ChangePassword(context.Background(), ChangePasswordParams{
		Username: user1.Username,
		Password: newPasswordHash,
	})
	require.NoError(t, err)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.Equal(t, newPasswordHash, user2.Password)
	require.NotEqual(t, user1.Password, user2.Password)
	require.NoError(t, common.CheckPassword(newPassword, user2.Password))
	require.Error(t, common.CheckPassword(newPassword, user1.Password))
	require.True(t, user2.PasswordChangedAt.Valid)
	require.WithinDuration(t, time.Now(), user2.PasswordChangedAt.Time, time.Second)
	require.WithinDuration(t, user2.PasswordChangedAt.Time, user2.UpdatedAt.Time, time.Second)
}

func TestListUsers(t *testing.T) {
	// สร้างเพิ่ม 5 คน ให้แน่ใจว่ามีข้อมูลพอ list
	for i := 0; i < 5; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{Limit: 5, Offset: 0}
	users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5) // limit 5 → ได้ 5 พอดี

	for _, u := range users {
		require.NotEmpty(t, u)
		// u เป็น ListUsersRow — ไม่มี field Password ให้เข้าถึงเลย (ปิดที่ query แล้ว)
	}
}
