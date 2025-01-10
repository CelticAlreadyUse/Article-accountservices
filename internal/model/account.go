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
	FindByID(ctx context.Context, id int64) (*Account, error)
	Update(ctx context.Context, account Account, id int64) (*Account, error)
	FindByIDs(ctx context.Context, ids []int64) ([]*Account, error)
	FindByUserName(ctx context.Context, search SearchParam) ([]*SearchModelResponse, error)
}
type AccountUsecase interface {
	Create(ctx context.Context, data Register) (token string, err error)
	Login(ctx context.Context, data Login) (token string, err error)
	FindByID(ctx context.Context, data Account, id int64) (*Account, error)
	Update(ctx context.Context, data Account, id int64) (*Account, error)
	FindByIDs(ctx context.Context, ids []int64) ([]*Account, error)
	Search(ctx context.Context, search SearchParam) []*SearchModelResponse

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
	ID         int64          `json:"id"`
	Fullname   sql.NullString `json:"fullname"`
	SortBio    sql.NullString `json:"sort_bio"`
	Gender     Gender         `json:"gender"`
	PictureUrl sql.NullString `json:"picture_url"`
	Username   string         `json:"-"`
	Email      string         `json:"-"`
	Verify     sql.NullBool   `json:"-"`
	Password   string         `json:"-"`
	Role       Role           `json:"role,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  time.Time      `json:"deleted_at"`
}

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}
type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}
type ConfigJWT struct {
	SigningKey string
	ExpTime    string
}
type Login struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
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
