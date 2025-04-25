package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextAuthKey string

const BearerAuthKey ContextAuthKey = "BearerAuth"

type AccountRepository interface {
	Store(ctx context.Context, data Account) (*Account, error)
	FindByEmail(ctx context.Context, email string) *Login
	FindByID(ctx context.Context, id string) (*Account, error)
	Update(ctx context.Context, account Account, id string) (*Account, error)
	FindByIDs(ctx context.Context, ids []string) ([]*Account, error)
	FindByUserName(ctx context.Context, search SearchParam) ([]*SearchModelResponse, error)
	SetVerify(ctx context.Context, email string) error
	UpdatePassword(ctx context.Context, req ResetPasswordReq) error
}
type AccountUsecase interface {
	Create(ctx context.Context, data Register) (token string, err error)
	Login(ctx context.Context, data Login) (login *Login, err error)
	FindByID(ctx context.Context, id string) (*Account, error)
	Update(ctx context.Context, data Account, id string) (*Account, error)
	FindByIDs(ctx context.Context, ids []string) ([]*Account, error)
	Search(ctx context.Context, search SearchParam) []*SearchModelResponse
	SetVerify(ctx context.Context, email string) error
}

type Gender string
type Role string

const (
	MALE   Gender = "male"
	Female Gender = "female"
	OTHERS Gender = "others"
)

const (
	USER  Role = "member"
	ADMIN Role = "admin"
)

type SearchParam struct {
	Limit    int64
	Username string
}

type SearchModelResponse struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	PictureUrl string    `json:"picture_url"`
	SortBio    string    `json:"sort_bio"`
	CreatedAt  time.Time `json:"created_at"`
}
type Account struct {
	ID           string          `json:"id"`
	DisplayName  sql.NullString `json:"display_name"`
	ShortBio     sql.NullString `json:"short_bio"`
	Gender       Gender         `json:"gender"`
	PictureUrl   sql.NullString `json:"picture_url"`
	Username     string         `json:"username"`
	Email        string         `json:"email"`
	EmailVerify  sql.NullBool   `json:"EmailVerified"`
	HashPassword string         `json:"-"`
	Role         Role           `json:"role,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}
type CustomClaims struct {
	UserID        string  `json:"user_id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"verified_eamil"`
	Role          Role   `json:"role"`
	jwt.RegisteredClaims
}
type ConfigJWT struct {
	SigningKey string
	ExpTime    string
}
type Login struct {
	ID            string  `json:"id"`
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password,omitempty"`
	Username      string `json:"username"`
	Role          Role   `json:"role"`
	EmailVerified bool   `json:"email_verify"`
	Token         string `json:"access_token"`
}
type VerifyEmail struct {
	ID        int64     `json:"-"`
	UserID    int64     `json:"-"`
	Token     string    `json:"email_token"`
	ExpiresAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type Cookie struct {
	Name    string
	Value   string
	Expired string
}
