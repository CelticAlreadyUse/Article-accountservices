package usecase

import (
	"context"
	"errors"

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
func (*accountUsecase) FindByID(ctx context.Context, data model.Account, id int64) (*model.Account, error) {
	panic("implement me")
}
func (*accountUsecase) FindByIDs(ctx context.Context, ids []int64) ([]*model.Account, error) {
	panic("implement me")
}
func (u *accountUsecase) Login(ctx context.Context, data model.Login) (string, error) {
	logger := logrus.WithFields(logrus.Fields{
		"data": data,
	})

	user := u.accountRepository.FindByEmail(ctx, data.Email)
	if user == nil {
		err := errors.New("email not found")
		return " ", err
	}
	if !helper.CheckPasswword(data.Password, user.Password) {
		logger.Errorf("miss match password for  %d", user.ID)
		return "", errors.New("miss match password")
	}
	token, err := helper.GenerateToken(user.ID)
	if err != nil {
		logger.Error(err)
	}
	return token, nil
}
func (*accountUsecase) Update(ctx context.Context, data model.Account, id int64) (*model.Account, error) {
	panic("implement me")
}
