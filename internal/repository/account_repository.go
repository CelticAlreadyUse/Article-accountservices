package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/helper"
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
	data.Username = "@" + data.Username
	result, err := sq.Insert("accounts").Columns("username", "email", "email_verify", "password_hash", "created_at", "updated_at").
		Values(data.Username, data.Email, false, data.HashPassword, now, now).RunWith(r.db).ExecContext(ctx)
	if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
		return nil, ErrDuplicateEntry
	}
	if err != nil {
		logrus.WithField("data", data).Error(err)
		return nil, err
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		logrus.Error("data", err)
		return nil, err
	} else {
		logrus.Infof("last insert ID : %d", rowAffected)
	}
	newAccount := &data
	newAccount.CreatedAt = now
	return newAccount, nil
}
func (r *accountRepository) FindByEmail(ctx context.Context, email string) *model.Login {
	row := sq.Select("id", "email", "password_hash", "username", "role","email_verify").
		From("accounts").
		Where(sq.Eq{"email": email}).
		RunWith(r.db).
		QueryRowContext(ctx)
	var data model.Login
	err := row.Scan(
		&data.ID,
		&data.Email,
		&data.Password,
		&data.Username,
		&data.Role,
		&data.EmailVerified,
	)
	if err != nil {
		return nil
	}
	return &data
}
func (r *accountRepository) FindByID(ctx context.Context, id string) (*model.Account, error) {
	logg := logrus.WithFields(logrus.Fields{
		"data": id,
	})
	row := sq.Select("display_name", "short_bio", "gender", "picture_url").From("accounts").Where(sq.Eq{"id": id}).RunWith(r.db).QueryRowContext(ctx)
	var data model.Account
	err := row.Scan(
		&data.DisplayName,
		&data.ShortBio,
		&data.Gender,
		&data.PictureUrl)
	if err != nil {
		logg.Error(err.Error())
		return nil, err
	}
	return &data, nil
}
func (r *accountRepository) Update(ctx context.Context, data model.Account, id string) (*model.Account, error) {
	timeNow := time.Now().UTC()
	_, err := sq.Update("accounts").
		Set("display_name", data.DisplayName).
		Set("short_bio", data.ShortBio).
		Set("gender", data.Gender).
		Set("picture_url", data.PictureUrl).
		Set("updated_at", timeNow).
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
func (r *accountRepository) FindByIDs(ctx context.Context, ids []string) ([]*model.Account, error) {
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
			id                  int64
			pictureUrl, sortBio sql.NullString
			username            string
			createdAt           time.Time
		)
		if err := rows.Scan(&id, &username, &pictureUrl, &sortBio, &createdAt); err != nil {
			return nil, err
		}
		searchAccounnts := model.SearchModelResponse{
			ID:         id,
			Username:   username,
			PictureUrl: pictureUrl.String,
			SortBio:    sortBio.String,
			CreatedAt:  createdAt,
		}
		accounts = append(accounts, &searchAccounnts)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	logrus.Infof("Searched Username: %s", search.Username)
	return accounts, nil
}
func (r *accountRepository) SetVerify(ctx context.Context, email string) error {
	_, err := sq.Update("accounts").Set("verify", true).Where(sq.Eq{"email": email}).RunWith(r.db).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (r *accountRepository) UpdatePassword(ctx context.Context, data model.ResetPasswordReq) error {
	hashedPassword, err := helper.Hashpassword(data.NewPass)
	_, err = sq.Update("accounts").Set("password", hashedPassword).Where(sq.Eq{"email": data.Email}).RunWith(r.db).ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
