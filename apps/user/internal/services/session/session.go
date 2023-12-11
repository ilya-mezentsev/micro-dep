package session

import (
	"errors"
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
}

func New(tokenRepo TokenRepo, authorRepo AuthorRepo) Service {
	return Service{
		tokenRepo:  tokenRepo,
		authorRepo: authorRepo,
	}
}

func (s Service) AuthorizedByToken(value string) (shared.Author, error) {
	ids, err := s.tokenRepo.AuthorizedAccountId(value, time.Now())
	if err != nil {
		if errors.Is(err, errs.IdMissingInStorage) {
			err = auth.AccountNotFoundErr
		}

		return shared.Author{}, err
	}

	return s.authorRepo.ById(ids.AuthorId)
}

func (s Service) AuthorizeByCredentials(auth shared.AuthorCreds) (shared.Author, shared.AuthResult, error) {
	author, err := s.authorRepo.ByCredentials(auth)
	if err != nil {
		if errors.Is(err, errs.KeyMissingInStorage) {
			err = CredentialsNotFound
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
		return shared.Author{}, shared.AuthResult{}, err
	}

	result := shared.AuthResult{
		Value:     token.Value,
		ExpiredAt: token.ExpiredAt,
	}

	return author, result, nil
}
