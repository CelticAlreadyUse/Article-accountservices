package model

import (
	"time"
)

type OTPRepository interface {
	StoreOTP(email string, otp string, ttl time.Duration) error
	ValidateOTP( data OTPRequestValidate) (bool, error)
}
type OTPUsecase interface {
	GenerateAndSendOTP(email OTPRequestGenerateAndSend) error
	ValidateOTP(data *OTPRequestValidate) (bool, error)
}
type OTPRequestGenerateAndSend struct {
	Email string `json:"email" validate:"required,email"`
}
type OTPRequestValidate struct {
	Email string `json:"email" validate:"required,email"`
	OTPCode string `json:"otp" validate:"required"`
}
