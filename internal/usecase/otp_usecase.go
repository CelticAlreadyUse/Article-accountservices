package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/helper"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/sirupsen/logrus"
)

type usecaseOTP struct {
	otpRepository model.OTPRepository
	accRepo       model.AccountRepository
}

func InitUsecaseOTP(otpRepository model.OTPRepository, accRepo model.AccountRepository) model.OTPUsecase {
	return &usecaseOTP{
		otpRepository: otpRepository,
		accRepo:       accRepo,
	}
}

func (u *usecaseOTP) GenerateAndSendOTP(data model.OTPRequestGenerateAndSend) error {
	logger := logrus.WithFields(logrus.Fields{
		"data": data.Email,
	})
	OTP := helper.GenerateOTP()
	err := u.otpRepository.StoreOTP(data.Email, OTP, 2*time.Minute)
	if err != nil {
		return err
	}
	err = helper.SendEmail(data.Email,
		"verification code",
		fmt.Sprintf("this is your code,Dont share it to anyone else : %s", OTP))
	if err != nil {
		return err
	}
	logger.Infof("OTP %s sent to email: %s\n", OTP, data.Email)
	return nil
}

func (u *usecaseOTP) ValidateOTP(data *model.OTPRequestValidate) (bool, error) {
	logrus.WithFields(logrus.Fields{
		"data": data,
	})

	return u.otpRepository.ValidateOTP(*data)
}
func (u *usecaseOTP) SendOTPPass(data model.OTPRequestGenerateAndSend) error {
	logger := logrus.WithFields(logrus.Fields{
		"data": data.Email,
	})
	account := u.accRepo.FindByEmail(context.Background(), data.Email)
	if account == nil {
		return errors.New("ops something wrong,email not found")
	}
	OTP := helper.GenerateOTP()
	err := u.otpRepository.StoredOTPPass(data.Email, OTP, 10*time.Minute)
	if err != nil {
		return err
	}
	err = helper.SendEmail(data.Email,
		"verification code",
		fmt.Sprintf("this is your password code,Dont share it to anyone else : %s", OTP))
	if err != nil {
		return err
	}
	logger.Infof("OTP %s sent to email: %s\n", OTP, data.Email)
	return nil
}
func (u *usecaseOTP) ValidateOTPGenerateToken(data *model.OTPRequestValidate) (string, error) {
	logrus.WithFields(logrus.Fields{
		"data": data,
	})
	ok, err := u.otpRepository.ValidateOTPPass(*data)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("invalid otp")
	}
	token, err := u.otpRepository.GenerateTokenPass(data.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}
func (u *usecaseOTP) ChangePassword(ctx context.Context, data model.ResetPasswordReq) error {
	return u.accRepo.UpdatePassword(ctx, data)
}
