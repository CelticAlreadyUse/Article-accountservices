package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/config"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/helper"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/sirupsen/logrus"
)

type accountUsecase struct {
	accountRepository model.AccountRepository
}

func NewAccountUsecase(accountRepository model.AccountRepository) model.AccountUsecase {
	return &accountUsecase{accountRepository: accountRepository}
}

func (u *accountUsecase) Create(ctx context.Context, data model.Register) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"data": data,
	})
	passwordHashed, err := helper.Hashpassword(data.Password)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	Email := u.accountRepository.FindByEmail(ctx, data.Email)
	if Email != nil {
		return "", errors.New("account was already exist")
	}
	newAccount, err := u.accountRepository.Store(ctx, model.Account{
		Username: data.Username,
		Email:    data.Email,
		Password: passwordHashed,
	})
	if err != nil {
		logger.Error(err)
		return "", err
	}
	accesToken, err := helper.GenerateToken(newAccount.ID)
	if err != nil {
		logger.Error(err)
		return "", err
	}
	return accesToken, nil
}
func (u *accountUsecase) FindByID(ctx context.Context, data model.Account, id int64) (*model.Account, error) {
	logrus.WithFields(logrus.Fields{
		"id":   id,
		"data": data,
	})

	account, err := u.accountRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("id not found")
	}

	return account, err
}

func (*accountUsecase) FindByIDs(ctx context.Context, ids []int64) ([]*model.Account, error) {
	return nil, errors.New("err")
}
func (u *accountUsecase) Login(ctx context.Context, data model.Login) (*model.Login, error) {
	logger := logrus.WithFields(logrus.Fields{
		"data": data,
	})
	user := u.accountRepository.FindByEmail(ctx, data.Email)
	if user == nil {
		err := errors.New("email not found")
		return nil, err
	}
	if !helper.CheckPasswword(data.Password, user.Password) {
		logger.Errorf("miss match password for  %d", user.ID)
		return nil, errors.New("miss match password")
	}
	token, err := helper.GenerateToken(user.ID)
	if err != nil {
		logger.Error(err)
	}
	Login := model.Login{
		ID: user.ID,
		Email: user.Email,
		Username: user.Username,
		Role: user.Role,
		Token: token,
	}
	claims, _ := helper.ValidateToken(token, model.ConfigJWT{SigningKey: config.JWTSigningKey(), ExpTime: config.JWTExp().String()})
	logger.Infof("Token akan expired pada: %v\n", claims.ExpiresAt.Time)
	logger.Infof("Waktu sekarang: %v\n", time.Now())
	logger.Infof("Sisa waktu: %v\n", claims.ExpiresAt.Time.Sub(time.Now()))
	return &Login, nil
}
func (u *accountUsecase) Update(ctx context.Context, data model.Account, id int64) (*model.Account, error) {
	logger := logrus.WithFields(logrus.Fields{
		"email": data.Email,
		"id":    id,
	})
	account, err := u.accountRepository.Update(ctx, data, id)
	if err != nil {
		logger.Error("failed to update account", err)
		return nil, err
	}

	logger.Info("Account update sucessfully")
	return account, nil
}
func (u *accountUsecase) Search(ctx context.Context, search model.SearchParam) []*model.SearchModelResponse {
	logrus.WithFields(logrus.Fields{
		"data": search.Username,
	})
	account, err := u.accountRepository.FindByUserName(ctx, search)
	if err != nil {
		return nil
	}
	return account
}
func (u *accountUsecase) SetVerify(ctx context.Context, email string) error {
	logrus.WithFields(logrus.Fields{
		"Data": email,
	})
	err := u.accountRepository.SetVerify(ctx, email)
	if err != nil {
		return err
	}
	return nil
}
