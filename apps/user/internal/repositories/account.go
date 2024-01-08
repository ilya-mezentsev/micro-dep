package repositories

import (
	"github.com/jmoiron/sqlx"

	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/shared"
)

const (
	createAccountQuery = `INSERT INTO account VALUES($1, $2)`
	accountExistsQuery = `SELECT EXISTS(SELECT 1 FROM account WHERE id = $1)`
)

type Account struct {
	db *sqlx.DB
}

func newAccount(db *sqlx.DB) Account {
	return Account{db: db}
}

func (a Account) Create(account shared.Account) error {
	_, err := a.db.Exec(createAccountQuery, account.Id, account.RegisteredAt)

	return err
}

func (a Account) Exists(id models.Id) (bool, error) {
	var result bool
	err := a.db.Get(&result, accountExistsQuery, string(id))

	return result, err
}
