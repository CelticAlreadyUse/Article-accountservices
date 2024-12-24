package repository
import (
	"context"
	"database/sql"
	"time"
	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)
type accountRepository struct{
	db *sql.DB
}
func InitAccountRepository (db *sql.DB)model.AccountRepository{
	return &accountRepository{db: db}
}
func (r *accountRepository) Store(ctx context.Context, data model.Account) (*model.Account, error){
	now := time.Now().UTC()
	result,err := sq.Insert("accounts").Columns("username","email","password","created_at","updated_at").
	Values(data.Username,data.Email,data.Password,now,now).RunWith(r.db).ExecContext(ctx)
	if err != nil{
		logrus.WithField("data",data).Error(err)
		return nil,err
	}
	lastInsertId,err:=result.LastInsertId()
	if err !=nil{
		logrus.Error("data",err)
		return nil, err
	}else{
		logrus.Infof("last insert ID : %d",lastInsertId)
	}
	rowAffected,err := result.RowsAffected()
	if err !=nil{
		logrus.Error("data",err)
		return nil, err
	}else{
		logrus.Infof("last insert ID : %d",rowAffected)
	}
	newAccount  := &data
	newAccount.ID = lastInsertId
	newAccount.CreatedAt = now
	return newAccount,nil
}
func (r *accountRepository) FindByEmail(ctx context.Context, email string) *model.Login{
panic("i")
}
func (r *accountRepository)FindByID(ctx context.Context, id int64) (*model.Account, error){
panic("")
}
func (r *accountRepository)Update(ctx context.Context, account model.Account, id int64) (*model.Account, error){
panic("")
}
func (r *accountRepository)FindByIDs(ctx context.Context, ids []int64) ([]*model.Account, error){
panic("")
}