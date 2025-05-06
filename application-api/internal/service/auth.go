package service

import (
	"context"
	"fmt"
	"github.com/vovk404/course-platform/application-api/internal/entity"
	"github.com/vovk404/course-platform/application-api/internal/storage"
	"github.com/vovk404/course-platform/application-api/pkg/auth"
	"github.com/vovk404/course-platform/application-api/pkg/errs"
	"github.com/vovk404/course-platform/application-api/pkg/hash"
)

type authService struct {
	serviceContext
	hash hash.Hash
	auth auth.Authenticator
}

var _ AuthService = (*authService)(nil)

func NewAuthService(options *Options) AuthService {
	return &authService{
		serviceContext: serviceContext{
			storages: options.Storages,
			config:   options.Config,
			logger:   options.Logger.Named("AuthService"),
		},
		hash: options.Hash,
		auth: options.Auth,
	}
}

func (a *authService) SignIn(ctx context.Context, options *SignInOptions) (*SignInOutput, error) {
	logger := a.logger.
		Named("SignIn").
		WithContext(ctx).
		With("options", options)

	user, err := a.storages.UserStorage.GetUser(ctx, &storage.GetUserFilter{Email: options.Email})
	if err != nil {
		logger.Error("failed to get user: ", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		logger.Info("user not found")
		return nil, ErrSignInUserNotFound
	}
	logger = logger.With("user", user)

	err = a.hash.CompareHash([]byte(user.Password), []byte(options.Password))
	if err != nil {
		logger.Info(err.Error())
		return nil, ErrSignInWrongPassword
	}

	accessToken, err := a.auth.GenerateToken(&auth.GenerateTokenClaimsOptions{UserName: user.Username, UserId: user.Id})
	if err != nil {
		logger.Error("failed to generate token for user: ", err)
		return nil, fmt.Errorf("failed to generate token for user: %w", err)
	}

	logger.Info("successfully signed user")
	return &SignInOutput{AccessToken: accessToken}, nil
}

func (a *authService) SignUp(ctx context.Context, options *SignUpOptions) (*SignUpOutput, error) {
	logger := a.logger.
		Named("SignUp").
		WithContext(ctx).
		With("options", options)

	user, err := a.storages.UserStorage.GetUser(ctx, &storage.GetUserFilter{Email: options.Email})
	if err != nil {
		logger.Error("failed to get user: ", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user != nil {
		logger.Info("user already created")
		return nil, ErrSignUpUserAlreadyCreated
	}

	hashedPassword, err := a.hash.GenerateHash(options.Password)
	if err != nil {
		logger.Error("failed to hash user password: ", err)
		return nil, fmt.Errorf("failed to hash user: %w", err)
	}

	createdUser, err := a.storages.UserStorage.CreateUser(ctx, &entity.User{
		Email:    options.Email,
		Password: hashedPassword,
		Username: options.Username,
		Type:     options.Type,
	})
	if err != nil {
		logger.Error("failed to create user: ", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	accessToken, err := a.auth.GenerateToken(&auth.GenerateTokenClaimsOptions{UserName: createdUser.Username, UserId: createdUser.Id})
	if err != nil {
		logger.Error("failed to generate token for user: ", err)
		return nil, fmt.Errorf("failed to generate token for user: %w", err)
	}

	logger.Info("successfully handled sign up")
	return &SignUpOutput{
		Id:          createdUser.Id,
		Username:    createdUser.Username,
		Email:       createdUser.Email,
		Type:        createdUser.Type,
		AccessToken: accessToken,
	}, nil
}

func (a *authService) VerifyToken(ctx context.Context, options *VerifyTokenOptions) (*VerifyTokenOutput, error) {
	logger := a.logger.
		Named("VerifyToken").
		WithContext(ctx)

	claims, err := a.auth.ParseToken(options.AccessToken)
	if err != nil {
		logger.Info("failed to parse token: ", err)
		return nil, fmt.Errorf("invalid token")
	}

	logger.Info("successfully handled auth token")
	return &VerifyTokenOutput{Username: claims.Username, UserId: claims.UserId}, nil
}

func (s *SignUpOptions) Validate() error {
	fmt.Println("User type: ", s.Type)
	if s.Type > 2 || s.Type < 1 {
		return errs.New("Type must be either 1 or 2, which means student or teacher.", "wrong user type")
	}
	return nil
}
