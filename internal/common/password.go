package common

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordCharacters = 8
	MaxPasswordCharacters = 64

	argon2Memory      = 19 * 1024
	argon2Iterations  = 2
	argon2Parallelism = 1
	argon2SaltLength  = 16
	argon2KeyLength   = 32
)

var (
	ErrPasswordMismatch    = errors.New("password does not match")
	ErrInvalidPasswordHash = errors.New("invalid password hash")
)

type argon2Parameters struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
}

// IsValidPassword applies the product's password policy. Spaces are allowed,
// but do not count as special characters. Control characters are rejected.
func IsValidPassword(password string) bool {
	if !utf8.ValidString(password) {
		return false
	}
	length := utf8.RuneCountInString(password)
	if length < MinPasswordCharacters || length > MaxPasswordCharacters {
		return false
	}

	var hasLetter, hasNumber, hasSpecial bool
	for _, character := range password {
		if unicode.IsControl(character) {
			return false
		}
		switch {
		case unicode.IsLetter(character):
			hasLetter = true
		case unicode.IsNumber(character):
			hasNumber = true
		case unicode.IsPunct(character) || unicode.IsSymbol(character):
			hasSpecial = true
		}
	}
	return hasLetter && hasNumber && hasSpecial
}

// HashPassword creates a salted Argon2id hash encoded in PHC string format.
// Store only the returned hash; callers must never persist plaintext.
func HashPassword(password string) (string, error) {
	salt := make([]byte, argon2SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argon2Iterations,
		argon2Memory,
		argon2Parallelism,
		argon2KeyLength,
	)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argon2Memory,
		argon2Iterations,
		argon2Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

// CheckPassword verifies Argon2id hashes and keeps read compatibility with
// bcrypt hashes created before the application migrated to Argon2id.
func CheckPassword(password string, encodedHash string) error {
	if strings.HasPrefix(encodedHash, "$2a$") ||
		strings.HasPrefix(encodedHash, "$2b$") ||
		strings.HasPrefix(encodedHash, "$2y$") {
		err := bcrypt.CompareHashAndPassword([]byte(encodedHash), []byte(password))
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordMismatch
		}
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidPasswordHash, err)
		}
		return nil
	}

	parameters, salt, expectedHash, err := decodeArgon2Hash(encodedHash)
	if err != nil {
		return err
	}
	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		parameters.iterations,
		parameters.memory,
		parameters.parallelism,
		uint32(len(expectedHash)),
	)
	if subtle.ConstantTimeCompare(actualHash, expectedHash) != 1 {
		return ErrPasswordMismatch
	}
	return nil
}

func decodeArgon2Hash(encodedHash string) (argon2Parameters, []byte, []byte, error) {
	if len(encodedHash) > 512 {
		return argon2Parameters{}, nil, nil, ErrInvalidPasswordHash
	}
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" {
		return argon2Parameters{}, nil, nil, ErrInvalidPasswordHash
	}

	version, err := parseHashParameter(parts[2], "v=", 8)
	if err != nil || version != argon2.Version {
		return argon2Parameters{}, nil, nil, ErrInvalidPasswordHash
	}

	parameterParts := strings.Split(parts[3], ",")
	if len(parameterParts) != 3 {
		return argon2Parameters{}, nil, nil, ErrInvalidPasswordHash
	}
	memory, memoryErr := parseHashParameter(parameterParts[0], "m=", 32)
	iterations, iterationsErr := parseHashParameter(parameterParts[1], "t=", 32)
	parallelism, parallelismErr := parseHashParameter(parameterParts[2], "p=", 8)
	if memoryErr != nil || iterationsErr != nil || parallelismErr != nil ||
		memory < 8*1024 || memory > 256*1024 ||
		iterations < 1 || iterations > 10 ||
		parallelism < 1 || parallelism > 16 {
		return argon2Parameters{}, nil, nil, ErrInvalidPasswordHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) < 16 || len(salt) > 64 {
		return argon2Parameters{}, nil, nil, ErrInvalidPasswordHash
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(expectedHash) < 16 || len(expectedHash) > 64 {
		return argon2Parameters{}, nil, nil, ErrInvalidPasswordHash
	}

	return argon2Parameters{
		memory:      uint32(memory),
		iterations:  uint32(iterations),
		parallelism: uint8(parallelism),
	}, salt, expectedHash, nil
}

func parseHashParameter(value string, prefix string, bitSize int) (uint64, error) {
	if !strings.HasPrefix(value, prefix) {
		return 0, ErrInvalidPasswordHash
	}
	parsed, err := strconv.ParseUint(strings.TrimPrefix(value, prefix), 10, bitSize)
	if err != nil {
		return 0, ErrInvalidPasswordHash
	}
	return parsed, nil
}
