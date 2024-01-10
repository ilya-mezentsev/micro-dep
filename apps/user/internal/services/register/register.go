package register

import (
	"log/slog"
	"time"

	"github.com/frankenbeanies/uuid4"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/shared"
)

type Service struct {
	accountRepo AccountRepo
	authorRepo  AuthorRepo
	logger      *slog.Logger
}

func New(
	accountRepo AccountRepo,
	authorRepo AuthorRepo,
	logger *slog.Logger,
) Service {

	return Service{
		accountRepo: accountRepo,
		authorRepo:  authorRepo,
		logger:      logger,
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
		s.logger.Error("Got an error while creating account", slog.Any("err", err))

		return shared.Author{}, errs.Unknown
	}

	// NOTE. If we failed here, account created above will remain in DB
	return s.register(accountId, creds)
}

func (s Service) validateUsername(username string) error {
	usernameExists, err := s.authorRepo.UsernameExists(username)
	if err != nil {
		s.logger.Error("Got an error while checking username existence", slog.Any("err", err))

		return errs.Unknown
	} else if usernameExists {
		return UsernameExists
	}

	return nil
}

func (s Service) RegisterForAccount(accountId models.Id, creds shared.AuthorCreds) (shared.Author, error) {
	accountExits, err := s.accountRepo.Exists(accountId)
	if err != nil {
		s.logger.Error("Got an error while checking account existence", slog.Any("err", err))

		return shared.Author{}, errs.Unknown
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
	if err != nil {
		s.logger.Error("Got an error while creating author", slog.Any("err", err))
		err = errs.Unknown
	}

	return result, err
}

func (s Service) AccountExists(accountId models.Id) (bool, error) {
	exists, err := s.accountRepo.Exists(accountId)
	if err != nil {
		s.logger.Error("Got an error while checking account existence", slog.Any("err", err))
		err = errs.Unknown
	}

	return exists, err
}
