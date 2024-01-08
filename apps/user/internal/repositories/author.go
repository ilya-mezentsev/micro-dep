package repositories

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/shared"
)

const (
	findAuthorByIdQuery    = `SELECT id, account_id, username, registered_at FROM author WHERE id = $1`
	findAuthorByCredsQuery = `SELECT id, account_id, username, registered_at FROM author WHERE username = $1 AND password = $2`
	createAuthorQuery      = `INSERT INTO author VALUES(:id, :account_id, :username, :password, :registered_at)`
	usernameExistsQuery    = `SELECT EXISTS(SELECT 1 FROM author WHERE username = $1)`
)

type (
	Author struct {
		db *sqlx.DB
	}

	authorProxy struct {
		Id           string `db:"id"`
		AccountId    string `db:"account_id"`
		Username     string `db:"username"`
		RegisteredAt int64  `db:"registered_at"`
	}
)

func newAuthor(db *sqlx.DB) Author {
	return Author{db: db}
}

func (a Author) Create(author shared.Author, password string) error {
	_, err := a.db.NamedExec(createAuthorQuery, authorProxy{}.mapFromModelAndPassword(author, password))

	return err
}

func (a Author) UsernameExists(username string) (bool, error) {
	var result bool
	err := a.db.Get(&result, usernameExistsQuery, username)

	return result, err
}

func (a Author) ById(authorId models.Id) (shared.Author, error) {
	var ap authorProxy
	err := a.db.Get(&ap, findAuthorByIdQuery, string(authorId))
	if errors.Is(err, sql.ErrNoRows) {
		err = errs.IdMissingInStorage
	}

	return ap.toModel(), err
}

func (a Author) ByCredentials(creds shared.AuthorCreds) (shared.Author, error) {
	var ap authorProxy
	err := a.db.Get(&ap, findAuthorByCredsQuery, creds.Username, creds.Password)
	if errors.Is(err, sql.ErrNoRows) {
		err = errs.KeyMissingInStorage
	}

	return ap.toModel(), err
}

func (ap authorProxy) toModel() shared.Author {
	return shared.Author{
		Id:           models.Id(ap.Id),
		AccountId:    models.Id(ap.AccountId),
		Username:     ap.Username,
		RegisteredAt: ap.RegisteredAt,
	}
}

func (ap authorProxy) mapFromModelAndPassword(author shared.Author, password string) map[string]interface{} {
	return map[string]interface{}{
		"id":            string(author.Id),
		"account_id":    string(author.AccountId),
		"username":      author.Username,
		"password":      password,
		"registered_at": author.RegisteredAt,
	}
}
