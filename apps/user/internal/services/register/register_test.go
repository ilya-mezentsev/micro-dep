package register

import (
	"testing"

	"github.com/frankenbeanies/uuid4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ilya-mezentsev/micro-dep/shared/types/models"
	registerMocks "github.com/ilya-mezentsev/micro-dep/user/internal/services/register/mocks"
	"github.com/ilya-mezentsev/micro-dep/user/internal/services/shared"
	sharedMocks "github.com/ilya-mezentsev/micro-dep/user/internal/services/shared/mocks"
)

var (
	someAccountId = models.Id(uuid4.New().String())

	someCreds = shared.AuthorCreds{
		Username: "my-user",
		Password: "my-password",
	}
)

func TestService_AccountExists(t *testing.T) {
	tests := []struct {
		name             string
		accountId        models.Id
		mocksConstructor func() (AccountRepo, AuthorRepo)
		expected         bool
		expectedErr      error
	}{
		{
			name:      "ok (true)",
			accountId: someAccountId,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				accountMock := registerMocks.NewMockAccountRepo(t)
				accountMock.EXPECT().Exists(someAccountId).Return(true, nil)

				return accountMock, nil
			},
			expected:    true,
			expectedErr: nil,
		},

		{
			name:      "ok (false)",
			accountId: someAccountId,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				accountMock := registerMocks.NewMockAccountRepo(t)
				accountMock.EXPECT().Exists(someAccountId).Return(false, nil)

				return accountMock, nil
			},
			expected:    false,
			expectedErr: nil,
		},

		{
			name:      "error",
			accountId: someAccountId,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				accountMock := registerMocks.NewMockAccountRepo(t)
				accountMock.EXPECT().Exists(someAccountId).Return(false, sharedMocks.SomeError)

				return accountMock, nil
			},
			expected:    false,
			expectedErr: sharedMocks.SomeError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, author := tt.mocksConstructor()
			service := New(acc, author)

			exists, err := service.AccountExists(tt.accountId)

			require.Equal(t, tt.expected, exists)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestService_Register(t *testing.T) {
	tests := []struct {
		name             string
		creds            shared.AuthorCreds
		mocksConstructor func() (AccountRepo, AuthorRepo)
		expectedErr      error
	}{
		{
			name:  "ok",
			creds: someCreds,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				accountMock := registerMocks.NewMockAccountRepo(t)
				accountMock.EXPECT().Create(mock.Anything).Return(nil)

				authorMock := registerMocks.NewMockAuthorRepo(t)
				authorMock.EXPECT().UsernameExists(someCreds.Username).Return(false, nil)
				authorMock.EXPECT().Create(mock.Anything, someCreds.Password).Return(nil)

				return accountMock, authorMock
			},
			expectedErr: nil,
		},

		{
			name:  "username exists",
			creds: someCreds,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				authorMock := registerMocks.NewMockAuthorRepo(t)
				authorMock.EXPECT().UsernameExists(someCreds.Username).Return(true, nil)

				return nil, authorMock
			},
			expectedErr: UsernameExists,
		},

		{
			name:  "failed to check username existence",
			creds: someCreds,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				authorMock := registerMocks.NewMockAuthorRepo(t)
				authorMock.EXPECT().UsernameExists(someCreds.Username).Return(false, sharedMocks.SomeError)

				return nil, authorMock
			},
			expectedErr: sharedMocks.SomeError,
		},

		{
			name:  "failed to create account",
			creds: someCreds,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				accountMock := registerMocks.NewMockAccountRepo(t)
				accountMock.EXPECT().Create(mock.Anything).Return(sharedMocks.SomeError)

				authorMock := registerMocks.NewMockAuthorRepo(t)
				authorMock.EXPECT().UsernameExists(someCreds.Username).Return(false, nil)

				return accountMock, authorMock
			},
			expectedErr: sharedMocks.SomeError,
		},

		{
			name:  "failed to create author",
			creds: someCreds,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				accountMock := registerMocks.NewMockAccountRepo(t)
				accountMock.EXPECT().Create(mock.Anything).Return(nil)

				authorMock := registerMocks.NewMockAuthorRepo(t)
				authorMock.EXPECT().UsernameExists(someCreds.Username).Return(false, nil)
				authorMock.EXPECT().Create(mock.Anything, someCreds.Password).Return(sharedMocks.SomeError)

				return accountMock, authorMock
			},
			expectedErr: sharedMocks.SomeError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, author := tt.mocksConstructor()
			service := New(acc, author)

			a, err := service.Register(someCreds)

			require.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				require.NotZero(t, a)
			}
		})
	}
}

func TestService_RegisterForAccount(t *testing.T) {
	tests := []struct {
		name             string
		accountId        models.Id
		creds            shared.AuthorCreds
		mocksConstructor func() (AccountRepo, AuthorRepo)
		expectedErr      error
	}{
		{
			name:      "ok",
			accountId: someAccountId,
			creds:     someCreds,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				accountMock := registerMocks.NewMockAccountRepo(t)
				accountMock.EXPECT().Exists(someAccountId).Return(true, nil)

				authorMock := registerMocks.NewMockAuthorRepo(t)
				authorMock.EXPECT().UsernameExists(someCreds.Username).Return(false, nil)
				authorMock.EXPECT().Create(mock.Anything, someCreds.Password).Return(nil)

				return accountMock, authorMock
			},
			expectedErr: nil,
		},

		{
			name:      "account not exists",
			accountId: someAccountId,
			creds:     someCreds,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				accountMock := registerMocks.NewMockAccountRepo(t)
				accountMock.EXPECT().Exists(someAccountId).Return(false, nil)

				return accountMock, nil
			},
			expectedErr: AccountNotFound,
		},

		{
			name:      "failed to check account existence",
			accountId: someAccountId,
			creds:     someCreds,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				accountMock := registerMocks.NewMockAccountRepo(t)
				accountMock.EXPECT().Exists(someAccountId).Return(false, sharedMocks.SomeError)

				return accountMock, nil
			},
			expectedErr: sharedMocks.SomeError,
		},

		{
			name:      "username exists",
			accountId: someAccountId,
			creds:     someCreds,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				accountMock := registerMocks.NewMockAccountRepo(t)
				accountMock.EXPECT().Exists(someAccountId).Return(true, nil)

				authorMock := registerMocks.NewMockAuthorRepo(t)
				authorMock.EXPECT().UsernameExists(someCreds.Username).Return(true, nil)

				return accountMock, authorMock
			},
			expectedErr: UsernameExists,
		},

		{
			name:      "failed to check username existence",
			accountId: someAccountId,
			creds:     someCreds,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				accountMock := registerMocks.NewMockAccountRepo(t)
				accountMock.EXPECT().Exists(someAccountId).Return(true, nil)

				authorMock := registerMocks.NewMockAuthorRepo(t)
				authorMock.EXPECT().UsernameExists(someCreds.Username).Return(false, sharedMocks.SomeError)

				return accountMock, authorMock
			},
			expectedErr: sharedMocks.SomeError,
		},

		{
			name:      "failed to create author",
			accountId: someAccountId,
			creds:     someCreds,
			mocksConstructor: func() (AccountRepo, AuthorRepo) {
				accountMock := registerMocks.NewMockAccountRepo(t)
				accountMock.EXPECT().Exists(someAccountId).Return(true, nil)

				authorMock := registerMocks.NewMockAuthorRepo(t)
				authorMock.EXPECT().UsernameExists(someCreds.Username).Return(false, nil)
				authorMock.EXPECT().Create(mock.Anything, someCreds.Password).Return(sharedMocks.SomeError)

				return accountMock, authorMock
			},
			expectedErr: sharedMocks.SomeError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, author := tt.mocksConstructor()
			service := New(acc, author)

			a, err := service.RegisterForAccount(tt.accountId, someCreds)

			require.Equal(t, tt.expectedErr, err)
			if tt.expectedErr == nil {
				require.NotZero(t, a)
			}
		})
	}
}
