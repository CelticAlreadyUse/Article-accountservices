package usecase

import (
	"context"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/helper"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/sirupsen/logrus"
)


type accountUsecase struct{
	accountRepository model.AccountRepository
}

func NewAccountUsecase(accountRepository model.AccountRepository) model.AccountUsecase {
	return &accountUsecase{accountRepository: accountRepository}
}

func (u *accountUsecase) Create(ctx context.Context, data model.Register) (token string, err error){
	logger := logrus.WithFields(logrus.Fields{
		"data":data,
	})
	passwordHashed,err := helper.Hashpassword(data.Password)
	if err !=nil{
		logger.Error(err)
		return "",err
	}
	newAccount,err := u.accountRepository.Store(ctx,model.Account{
		Username: data.Username,
		Email: data.Email,
		Password: passwordHashed,
	})
	if err !=nil{
		logger.Error(err)
		return "",nil
	}
	accesToken,err := helper.GenerateToken(newAccount.ID)
	if err !=nil{
		logger.Error(err)
		return "",err
	}

	return accesToken,nil
}
func (*accountUsecase) FindByID(ctx context.Context, data model.Account,id int64) (*model.Account,error){
	panic("implement me")
}
func (*accountUsecase) FindByIDs(ctx context.Context,ids []int64) ([]*model.Account,error){
	panic("implement me")
}
func (*accountUsecase) Login(ctx context.Context,login model.Login) (string,error){
	panic("implement me")
}
func (*accountUsecase) Update(ctx context.Context,data model.Account,id int64) (*model.Account,error){
	panic("implement me")
}