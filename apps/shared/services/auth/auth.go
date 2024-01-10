package auth

import (
	"errors"
	"log/slog"
	"time"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
)

var AccountNotFoundErr = errors.New("account-not-found")

type Service struct {
	tokenRepo TokenReaderRepo
	logger    *slog.Logger
}

func NewService(tokenRepo TokenReaderRepo, logger *slog.Logger) Service {
	return Service{
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}

func (s Service) IsAuthenticated(value string) (models.Id, error) {
	ids, err := s.tokenRepo.AuthorizedAccountId(value, time.Now())
	if errors.Is(err, errs.IdMissingInStorage) {
		err = AccountNotFoundErr
	} else if err != nil {
		s.logger.Error("Got an error while checking account authorization", slog.Any("err", err))
		err = errs.Unknown
	}

	return ids.AccountId, err
}
