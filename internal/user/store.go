package user

import (
	"context"

	db "simplebank/db/sqlc"
)

// Store lists only the database operation required by this feature.
type Store interface {
	CreateUser(context.Context, db.CreateUserParams) (db.User, error)
}
