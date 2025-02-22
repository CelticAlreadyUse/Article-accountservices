package model

import (
	"context"
	"time"
)
type OTPRepository interface {
	//verified Email
	StoreOTP(email string, otp string, ttl time.Duration) error
	ValidateOTP( data OTPRequestValidate) (bool, error)

	//Forgot Password
	StoredOTPPass(email string,otp string,ttl time.Duration)error
	ValidateOTPPass(validate OTPRequestValidate)(bool,error)
	GenerateTokenPass(email string)(string,error)
	UpdatePassword(ctx context.Context,newPass string)(error)
}
type OTPUsecase interface {
	//verified Email
	GenerateAndSendOTP(email OTPRequestGenerateAndSend) error
	ValidateOTP(data *OTPRequestValidate) (bool, error)

	//Forgot Password
	SendOTPPass(email OTPRequestGenerateAndSend) error
	ValidateOTPGenerateToken(data *OTPRequestValidate) (string, error)
	ChangePassword(ctx context.Context,newPass string)error
}

//otp for verified Email
type OTPRequestGenerateAndSend struct {
	Email string `json:"email" validate:"required,email"`
}
type OTPRequestValidate struct {
	Email string `json:"email" validate:"required,email"`
	OTPCode string `json:"otp" validate:"required"`
}