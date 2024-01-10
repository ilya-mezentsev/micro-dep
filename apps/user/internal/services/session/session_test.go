package session

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/frankenbeanies/uuid4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ilya-mezentsev/micro-dep/shared/errs"
	"github.com/ilya-mezentsev/micro-dep/shared/services/auth"
	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	sessionMocks "github.com/ilya-mezentsev/micro-dep/user/internal/services/session/mocks"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/shared"
	sharedMocks "github.com/ilya-mezentsev/micro-dep/user/internal/services/shared/mocks"
)

const (
	justToken = "just-token"
	sixDays   = 6 * 24 * 60 * 60
)

var (
	justAccountId = models.Id(uuid4.New().String())
	justAuthorId  = models.Id(uuid4.New().String())

	justAuthorizeIds = auth.AuthorizedIds{
		AuthorId:  justAuthorId,
		AccountId: justAccountId,
	}

	justAuthor = shared.Author{
		Id:        justAuthorId,
		AccountId: justAccountId,
	}

	justCreds = shared.AuthorCreds{
		Username: "username",
		Password: "password",
	}
)

func TestService_AuthorizedByToken(t *testing.T) {
	tests := []struct {
		name             string
		token            string
		mocksConstructor func() (TokenRepo, AuthorRepo)
		expectedAuthor   shared.Author
		expectedErr      error
	}{
		{
			name:  "success",
			token: justToken,
			mocksConstructor: func() (TokenRepo, AuthorRepo) {
				tokenRepo := sessionMocks.NewMockTokenRepo(t)
				tokenRepo.EXPECT().AuthorizedAccountId(justToken, mock.Anything).Return(justAuthorizeIds, nil)

				authorRepo := sessionMocks.NewMockAuthorRepo(t)
				authorRepo.EXPECT().ById(justAuthorId).Return(justAuthor, nil)

				return tokenRepo, authorRepo
			},
			expectedAuthor: justAuthor,
			expectedErr:    nil,
		},

		{
			name:  "missed token",
			token: justToken,
			mocksConstructor: func() (TokenRepo, AuthorRepo) {
				tokenRepo := sessionMocks.NewMockTokenRepo(t)
				tokenRepo.EXPECT().AuthorizedAccountId(justToken, mock.Anything).Return(justAuthorizeIds, errs.IdMissingInStorage)

				return tokenRepo, nil
			},
			expectedAuthor: shared.Author{},
			expectedErr:    auth.AccountNotFoundErr,
		},

		{
			name:  "some error from token repo",
			token: justToken,
			mocksConstructor: func() (TokenRepo, AuthorRepo) {
				tokenRepo := sessionMocks.NewMockTokenRepo(t)
				tokenRepo.EXPECT().AuthorizedAccountId(justToken, mock.Anything).Return(justAuthorizeIds, sharedMocks.SomeError)

				return tokenRepo, nil
			},
			expectedAuthor: shared.Author{},
			expectedErr:    errs.Unknown,
		},

		{
			name:  "some error from author repo",
			token: justToken,
			mocksConstructor: func() (TokenRepo, AuthorRepo) {
				tokenRepo := sessionMocks.NewMockTokenRepo(t)
				tokenRepo.EXPECT().AuthorizedAccountId(justToken, mock.Anything).Return(justAuthorizeIds, nil)

				authorRepo := sessionMocks.NewMockAuthorRepo(t)
				authorRepo.EXPECT().ById(justAuthorId).Return(shared.Author{}, sharedMocks.SomeError)

				return tokenRepo, authorRepo
			},
			expectedAuthor: shared.Author{},
			expectedErr:    errs.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, ar := tt.mocksConstructor()
			s := New(tr, ar, slog.New(slog.NewTextHandler(io.Discard, nil)))

			author, err := s.AuthorizedByToken(tt.token)

			require.Equal(t, tt.expectedAuthor, author)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestService_AuthorizeByCredentials(t *testing.T) {
	tests := []struct {
		name             string
		creds            shared.AuthorCreds
		mocksConstructor func() (TokenRepo, AuthorRepo)
		expectedAuthor   shared.Author
		expectedErr      error
	}{
		{
			name:  "success",
			creds: justCreds,
			mocksConstructor: func() (TokenRepo, AuthorRepo) {
				authorRepo := sessionMocks.NewMockAuthorRepo(t)
				authorRepo.EXPECT().ByCredentials(justCreds).Return(justAuthor, nil)

				tokenRepo := sessionMocks.NewMockTokenRepo(t)
				tokenRepo.EXPECT().Create(mock.Anything).Return(nil)

				return tokenRepo, authorRepo
			},
			expectedAuthor: justAuthor,
			expectedErr:    nil,
		},

		{
			name:  "credentials not found error",
			creds: justCreds,
			mocksConstructor: func() (TokenRepo, AuthorRepo) {
				authorRepo := sessionMocks.NewMockAuthorRepo(t)
				authorRepo.EXPECT().ByCredentials(justCreds).Return(justAuthor, errs.KeyMissingInStorage)

				return nil, authorRepo
			},
			expectedAuthor: shared.Author{},
			expectedErr:    CredentialsNotFound,
		},

		{
			name:  "some error from author repo",
			creds: justCreds,
			mocksConstructor: func() (TokenRepo, AuthorRepo) {
				authorRepo := sessionMocks.NewMockAuthorRepo(t)
				authorRepo.EXPECT().ByCredentials(justCreds).Return(justAuthor, sharedMocks.SomeError)

				return nil, authorRepo
			},
			expectedAuthor: shared.Author{},
			expectedErr:    errs.Unknown,
		},

		{
			name:  "some error from token repo",
			creds: justCreds,
			mocksConstructor: func() (TokenRepo, AuthorRepo) {
				authorRepo := sessionMocks.NewMockAuthorRepo(t)
				authorRepo.EXPECT().ByCredentials(justCreds).Return(justAuthor, nil)

				tokenRepo := sessionMocks.NewMockTokenRepo(t)
				tokenRepo.EXPECT().Create(mock.Anything).Return(sharedMocks.SomeError)

				return tokenRepo, authorRepo
			},
			expectedAuthor: shared.Author{},
			expectedErr:    errs.Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, ar := tt.mocksConstructor()
			s := New(tr, ar, slog.New(slog.NewTextHandler(io.Discard, nil)))

			author, authResult, err := s.AuthorizeByCredentials(tt.creds)

			require.Equal(t, tt.expectedAuthor, author)
			require.Equal(t, tt.expectedErr, err)

			if tt.expectedErr == nil {
				require.True(t, authResult.ExpiredAt > time.Now().Add(sixDays*time.Second).Unix())
			}
		})
	}
}
