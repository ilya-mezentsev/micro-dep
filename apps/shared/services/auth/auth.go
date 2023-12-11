package auth

import (
	"errors"
	"time"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

var AccountNotFoundErr = errors.New("account-not-found")

type Service struct {
	tokenRepo TokenReaderRepo
}

func NewService(tokenRepo TokenReaderRepo) Service {
	return Service{tokenRepo: tokenRepo}
}

func (s Service) IsAuthenticated(value string) (models.Id, error) {
	ids, err := s.tokenRepo.AuthorizedAccountId(value, time.Now())
	if errors.Is(err, errs.IdMissingInStorage) {
		err = AccountNotFoundErr
	}

	return ids.AccountId, err
}
