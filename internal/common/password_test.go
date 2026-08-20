package common

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/require"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "correct horse battery staple"

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEqual(t, password, hashedPassword)
	require.NoError(t, CheckPassword(password, hashedPassword))
	require.Error(t, CheckPassword("wrong password", hashedPassword))
}

func TestHashPasswordRejectsMoreThan72Bytes(t *testing.T) {
	_, err := HashPassword(strings.Repeat("a", 73))
	require.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
}

func TestCheckPasswordRejectsInvalidHash(t *testing.T) {
	err := CheckPassword("password", "not-a-bcrypt-hash")
	require.Error(t, err)
}
