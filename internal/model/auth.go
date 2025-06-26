package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	ID              string `json:"id"`
	Role            string `json:"role"`
	IsEmailVerified bool   `json:"is_email_verified"`
}

type JWTClaims struct {
	UserID          string `json:"user_id"`
	Role            string `json:"role"`
	IsEmailVerified bool   `json:"is_email_verified"`
	jwt.RegisteredClaims
}
