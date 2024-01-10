package session

import (
	"errors"
	"log/slog"
	"time"

	"github.com/frankenbeanies/uuid4"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/services/auth"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/shared"
)

const sevenDays = 7 * 24 * 60 * 60

type Service struct {
	tokenRepo  TokenRepo
	authorRepo AuthorRepo
	logger     *slog.Logger
}

func New(
	tokenRepo TokenRepo,
	authorRepo AuthorRepo,
	logger *slog.Logger,
) Service {

	return Service{
		tokenRepo:  tokenRepo,
		authorRepo: authorRepo,
		logger:     logger,
	}
}

func (s Service) AuthorizedByToken(value string) (shared.Author, error) {
	ids, err := s.tokenRepo.AuthorizedAccountId(value, time.Now())
	if err != nil {
		if errors.Is(err, errs.IdMissingInStorage) {
			err = auth.AccountNotFoundErr
		} else {
			s.logger.Error("Got an error while trying to authorize author by token", slog.Any("err", err))
			err = errs.Unknown
		}

		return shared.Author{}, err
	}

	author, err := s.authorRepo.ById(ids.AuthorId)
	if err != nil {
		s.logger.Error("Got an error while fetching author by id", slog.Any("err", err))

		// if we have authorized account here we cannot get missing-id error
		err = errs.Unknown
	}

	return author, err
}

func (s Service) AuthorizeByCredentials(auth shared.AuthorCreds) (shared.Author, shared.AuthResult, error) {
	author, err := s.authorRepo.ByCredentials(auth)
	if err != nil {
		if errors.Is(err, errs.KeyMissingInStorage) {
			err = CredentialsNotFound
		} else {
			s.logger.Error("Got an error while trying to authorize author by creds", slog.Any("err", err))
			err = errs.Unknown
		}

		return shared.Author{}, shared.AuthResult{}, err
	}

	token := shared.AuthToken{
		AuthorId:  author.Id,
		Value:     uuid4.New().String(),
		CreatedAt: time.Now().Unix(),
		ExpiredAt: time.Now().Add(sevenDays * time.Second).Unix(),
	}

	err = s.tokenRepo.Create(token)
	if err != nil {
		s.logger.Error("Got an error while creating auth token", slog.Any("err", err))

		return shared.Author{}, shared.AuthResult{}, errs.Unknown
	}

	result := shared.AuthResult{
		Value:     token.Value,
		ExpiredAt: token.ExpiredAt,
	}

	return author, result, nil
}
