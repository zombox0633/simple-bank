package user

import (
	"regexp"
	"strings"
	"unicode/utf8"

	z "github.com/Oudwins/zog"

	"simplebank/internal/common"
	"simplebank/internal/httpapi"
)

var (
	usernamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]{1,30}[a-z0-9]$`)
	fullNamePattern = regexp.MustCompile(`^[\p{L}\p{M} .'’\-]+$`)
)

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

var createUserSchema = z.Struct(z.Shape{
	"username": z.String().Trim().Transform(httpapi.Lowercase).
		Match(usernamePattern, z.IssueCode(httpapi.IssueInvalidUsername)).
		Required(z.IssueCode(httpapi.IssueInvalidUsername)),
	"password": z.String().TestFunc(
		func(value *string, _ z.Ctx) bool { return common.IsValidPassword(*value) },
		z.IssueCode(httpapi.IssueInvalidPassword),
	).Required(z.IssueCode(httpapi.IssueInvalidPassword)),
	"fullName": z.String().Trim().Transform(httpapi.CollapseWhitespace).Match(
		fullNamePattern,
		z.IssueCode(httpapi.IssueInvalidFullName),
	).TestFunc(
		func(value *string, _ z.Ctx) bool {
			length := utf8.RuneCountInString(*value)
			return length > 0 && length <= 100
		},
		z.IssueCode(httpapi.IssueInvalidFullName),
	).Required(z.IssueCode(httpapi.IssueInvalidFullName)),
	"email": z.String().Trim().Transform(httpapi.Lowercase).
		Email(z.IssueCode(httpapi.IssueInvalidEmail)).
		TestFunc(func(value *string, _ z.Ctx) bool {
			if len(*value) > 254 {
				return false
			}
			at := strings.LastIndexByte(*value, '@')
			return at > 0 && at <= 64
		}, z.IssueCode(httpapi.IssueInvalidEmail)).
		Required(z.IssueCode(httpapi.IssueInvalidEmail)),
})
