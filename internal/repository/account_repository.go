package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type accountRepository struct {
	db *sql.DB
}

func InitAccountRepository(db *sql.DB) model.AccountRepository {
	return &accountRepository{db: db}
}
var ErrDuplicateEntry = errors.New("username or email already exist")

func (r *accountRepository) Store(ctx context.Context, data model.Account) (*model.Account, error) {
	now := time.Now().UTC()
	result, err := sq.Insert("accounts").Columns("username", "email","verify","password", "created_at", "updated_at").
		Values(data.Username, data.Email,false ,data.Password, now, now).RunWith(r.db).ExecContext(ctx)
	if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
		return nil, ErrDuplicateEntry
	}
	if err != nil {
		logrus.WithField("data", data).Error(err)
		return nil, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		logrus.Error("data", err)
		return nil, err
	} else {
		logrus.Infof("last insert ID : %d", lastInsertId)
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		logrus.Error("data", err)
		return nil, err
	} else {
		logrus.Infof("last insert ID : %d", rowAffected)
	}
	newAccount := &data
	newAccount.ID = lastInsertId
	newAccount.CreatedAt = now
	return newAccount, nil
}
func (r *accountRepository) FindByEmail(ctx context.Context, email string) *model.Login {
	row := sq.Select("id", "email", "password").
		From("accounts").
		Where(sq.Eq{"email": email}).
		RunWith(r.db).
		QueryRowContext(ctx)
	var data model.Login
	err := row.Scan(
		&data.ID,
		&data.Email,
		&data.Password,
	)
	if err != nil {
		return nil
	}
	return &data
}
func (r *accountRepository) FindByID(ctx context.Context, id int64) (*model.Account, error) {
	logg := logrus.WithFields(logrus.Fields{
		"data": id,
	})
	row := sq.Select("fullname", "sort_bio", "gender", "picture_url").From("accounts").Where(sq.Eq{"id":id}).RunWith(r.db).QueryRowContext(ctx)
	var data model.Account
	err := row.Scan(
		&data.Fullname,
		&data.SortBio,
		&data.Gender,
		&data.PictureUrl)
	if err != nil {
		logg.Error(err.Error())
		return nil, err
	}
	return &data, nil
}
func (r *accountRepository) Update(ctx context.Context, data model.Account, id int64) (*model.Account, error) {
	timeNow := time.Now().UTC()
	_, err := sq.Update("accounts").
		Set("fullname", data.Fullname).
		Set("sort_bio", data.SortBio).
		Set("gender", data.Gender).
		Set("picture_url", data.PictureUrl).
		Set("updated_at",timeNow).
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		ExecContext(ctx)
	if err != nil {
		return nil, err
	}
	updatedAccount, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return updatedAccount, nil
}
func (r *accountRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.Account, error) {
	panic("IMplement")
}
func (r *accountRepository) FindByUserName(ctx context.Context, search model.SearchParam) ([]*model.SearchModelResponse, error) {
	query := `SELECT id, username, picture_url, sort_bio, created_at 
	          FROM accounts 
	          WHERE username LIKE ? 
	          LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, "%"+search.Username+"%", search.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accounts []*model.SearchModelResponse
	for rows.Next() {
		var (

			 id int64
			 pictureUrl,sortBio sql.NullString
			 username string
			 createdAt time.Time
		)
		if err := rows.Scan(&id, &username, &pictureUrl, &sortBio, &createdAt); err != nil {
			return nil, err
		}
	searchAccounnts := model.SearchModelResponse{
		ID: id,
		Username: username,
		PictureUrl: pictureUrl.String,
		SortBio: sortBio.String,
		CreatedAt: createdAt,
	}
	accounts = append(accounts, &searchAccounnts)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	logrus.Infof("Searched Username: %s", search.Username)
	return accounts, nil
}

func(r *accountRepository)SetVerify(ctx context.Context,email string)error{
	_,err := sq.Update("accounts").Set("verify",true).Where(sq.Eq{"email":email}).RunWith(r.db).ExecContext(ctx)
	if err !=nil{
		return err
	}
	return nil
}