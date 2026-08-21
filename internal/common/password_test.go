package common

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/require"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "Correct horse battery staple 42!"

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEqual(t, password, hashedPassword)
	require.True(t, strings.HasPrefix(hashedPassword, "$argon2id$"))
	require.NoError(t, CheckPassword(password, hashedPassword))
	require.ErrorIs(t, CheckPassword("wrong password", hashedPassword), ErrPasswordMismatch)

	secondHash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEqual(t, hashedPassword, secondHash)
}

func TestCheckPasswordSupportsExistingBcryptHash(t *testing.T) {
	password := "Legacy123!"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	require.NoError(t, CheckPassword(password, string(hash)))
	require.ErrorIs(t, CheckPassword("wrong password", string(hash)), ErrPasswordMismatch)
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{name: "minimum policy", password: "Secret1!", want: true},
		{name: "UTF-8", password: "รหัสผ่าน1!", want: true},
		{name: "maximum characters", password: "A1!" + strings.Repeat("a", MaxPasswordCharacters-3), want: true},
		{name: "too short", password: "Abc1!", want: false},
		{name: "missing letter", password: "1234567!", want: false},
		{name: "missing number", password: "Password!", want: false},
		{name: "missing special", password: "Password1", want: false},
		{name: "spaces are not special", password: "Password1 ", want: false},
		{name: "control character", password: "Pass1!\nword", want: false},
		{name: "too many characters", password: "A1!" + strings.Repeat("a", MaxPasswordCharacters-2), want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, IsValidPassword(test.password))
		})
	}
}

func TestCheckPasswordRejectsInvalidHash(t *testing.T) {
	tests := []string{
		"not-a-password-hash",
		"$argon2id$v=19$m=999999999,t=2,p=1$c2FsdHNhbHRzYWx0c2FsdA$aGFzaGhhc2hoYXNoaGFzaGhhc2g",
		"$argon2id$v=99$m=19456,t=2,p=1$c2FsdHNhbHRzYWx0c2FsdA$aGFzaGhhc2hoYXNoaGFzaGhhc2g",
	}

	for _, hash := range tests {
		require.ErrorIs(t, CheckPassword("Password1!", hash), ErrInvalidPasswordHash)
	}
}
