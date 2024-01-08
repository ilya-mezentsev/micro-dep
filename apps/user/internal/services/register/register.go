package register

import (
	"time"

	"github.com/frankenbeanies/uuid4"

	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/shared"
)

type Service struct {
	accountRepo AccountRepo
	authorRepo  AuthorRepo
}

func New(
	accountRepo AccountRepo,
	authorRepo AuthorRepo,
) Service {

	return Service{
		accountRepo: accountRepo,
		authorRepo:  authorRepo,
	}
}

func (s Service) Register(creds shared.AuthorCreds) (shared.Author, error) {
	if err := s.validateUsername(creds.Username); err != nil {
		return shared.Author{}, err
	}

	accountId := models.Id(uuid4.New().String())
	err := s.accountRepo.Create(shared.Account{
		Id:           accountId,
		RegisteredAt: time.Now().Unix(),
	})
	if err != nil {
		return shared.Author{}, err
	}

	// NOTE. If we failed here, account created above will remain in DB
	return s.register(accountId, creds)
}

func (s Service) validateUsername(username string) error {
	usernameExists, err := s.authorRepo.UsernameExists(username)
	if err != nil {
		return err
	} else if usernameExists {
		return UsernameExists
	}

	return nil
}

func (s Service) RegisterForAccount(accountId models.Id, creds shared.AuthorCreds) (shared.Author, error) {
	accountExits, err := s.accountRepo.Exists(accountId)
	if err != nil {
		return shared.Author{}, err
	} else if !accountExits {
		return shared.Author{}, AccountNotFound
	}

	if err = s.validateUsername(creds.Username); err != nil {
		return shared.Author{}, err
	}

	return s.register(accountId, creds)
}

func (s Service) register(accountId models.Id, creds shared.AuthorCreds) (shared.Author, error) {
	result := shared.Author{
		Id:           models.Id(uuid4.New().String()),
		AccountId:    accountId,
		Username:     creds.Username,
		RegisteredAt: time.Now().Unix(),
	}
	err := s.authorRepo.Create(result, creds.Password)

	return result, err
}

func (s Service) AccountExists(accountId models.Id) (bool, error) {
	return s.accountRepo.Exists(accountId)
}
