package model

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID              uuid.UUID         `json:"id"`
	Username        string            `json:"username"`
	Email           string            `json:"email"`
	Role            string            `json:"role"`
	PhoneNumber     string            `json:"phone_number,omitempty"`
	Institution     string            `json:"institution,omitempty"`
	Address         string            `json:"address,omitempty"`
	City            string            `json:"city,omitempty"`
	Province        string            `json:"province,omitempty"`
	Points          int64             `json:"points"`
	Balance         int64             `json:"balance"`
	IsEmailVerified bool              `json:"is_email_verified"`
	Location        *LocationResponse `json:"location,omitempty"`
	AccessToken     string            `json:"access_token,omitempty"`
	RefreshToken    string            `json:"refresh_token,omitempty"`
	CreatedAt       *time.Time        `json:"created_at,omitempty"`
	UpdatedAt       *time.Time        `json:"updated_at,omitempty"`
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Role     string `json:"role" validate:"required,max=100"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type GetUserRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type DeleteUserRequest struct {
	ID string `json:"id" validate:"required"`
}
