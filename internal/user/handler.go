package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	db "simplebank/db/sqlc"
	"simplebank/internal/httpapi"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (handler *Handler) createUser(ctx *gin.Context) {
	var request createUserRequest
	if !httpapi.BindJSON(ctx, createUserSchema, &request) {
		return
	}

	createdUser, err := handler.service.CreateUser(
		ctx.Request.Context(),
		request.Username,
		request.Password,
		request.FullName,
		request.Email,
	)
	if err != nil {
		httpapi.WriteError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, newCreateUserResponse(createdUser))
}

func newCreateUserResponse(createdUser db.User) createUserResponse {
	return createUserResponse{
		Username:          createdUser.Username,
		FullName:          createdUser.FullName,
		Email:             createdUser.Email,
		PasswordChangedAt: createdUser.PasswordChangedAt.Time.UTC(),
		CreatedAt:         createdUser.CreatedAt.Time,
		UpdatedAt:         createdUser.UpdatedAt.Time,
	}
}
