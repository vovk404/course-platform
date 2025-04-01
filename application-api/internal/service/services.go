package service

import (
	"context"
	"github.com/vovk404/course-platform/application-api/config"
	"github.com/vovk404/course-platform/application-api/internal/entity"
	"github.com/vovk404/course-platform/application-api/pkg/auth"
	"github.com/vovk404/course-platform/application-api/pkg/errs"
	"github.com/vovk404/course-platform/application-api/pkg/hash"
	"github.com/vovk404/course-platform/application-api/pkg/logger"
)

type Services struct {
	AuthService    AuthService
	AccountService AccountService
	NodeService    NodeService
	CourseService  CourseService
}

type Options struct {
	Storages *Storages
	Config   *config.Config
	Logger   logger.Logger
	Hash     hash.Hash
	Auth     auth.Authenticator
}

type serviceContext struct {
	storages *Storages
	config   *config.Config
	logger   logger.Logger
}

type AuthService interface {
	// SignIn provides logic of authentication of clients and returns access and refresh tokens.
	SignIn(ctx context.Context, options *SignInOptions) (*SignInOutput, error)
	// SignUp provides logic of creating the clients and returns access and refresh tokens.
	SignUp(ctx context.Context, options *SignUpOptions) (*SignUpOutput, error)
	// VerifyToken provides logic of validating provided authorization token.
	VerifyToken(ctx context.Context, options *VerifyTokenOptions) (*VerifyTokenOutput, error)
}

type SignInOptions struct {
	Email    string
	Password string
}

type SignInOutput struct {
	AccessToken string
}

type SignUpOptions struct {
	Username   string
	Email      string
	Password   string
	Type       int
	MacAddress string
}

type SignUpOutput struct {
	Id          string
	Username    string
	Email       string
	Type        int
	AccessToken string
}

type VerifyTokenOptions struct {
	AccessToken string
}

type VerifyTokenOutput struct {
	Username string
	UserId   string
}

var (
	ErrSignUpUserAlreadyCreated = errs.New("user already created", "user_already_created")
	ErrSignInUserNotFound       = errs.New("user not found", "user_not_found")
	ErrSignInWrongPassword      = errs.New("wrong password", "wrong_password")
)

type AccountService interface {
	// CreateAccount provides logic of creating account for clients.
	CreateAccount(ctx context.Context, options *CreateAccountOptions) (*CreateAccountOutput, error)
	// GetAccount provides logic of getting account via accountId.
	GetAccount(ctx context.Context, options *GetAccountOptions) (*entity.Account, error)
	// UpdateAccount provides logic of updating existing account.
	UpdateAccount(ctx context.Context, account *entity.Account) (*entity.Account, error)
}

type CreateAccountOptions struct {
	UserId           string `json:"userId"`
	DeviceName       string `json:"deviceName"`
	DeviceOS         string `json:"deviceOs"`
	DeviceMacAddress string `json:"deviceMacAddress"`
	Active           bool   `json:"active"`
	AccountLanguage  string `json:"accountLanguage"`
}

type CreateAccountOutput struct {
	Id     string `json:"id"`
	UserId string `json:"userId"`
}

type GetAccountOptions struct {
	AccountId string `json:"accountId"`
	UserId    string `json:"userId"`
}

var (
	ErrCreateAccountUserNotFound = errs.New("user not found", "user_not_found")
	ErrGetAccountAccountNotFound = errs.New("account not found", "account_not_found")
)

type NodeService interface {
	// CreateNode provides logic of creating new node.
	CreateNode(ctx context.Context, options *CreateNodeOptions) (*CreateNodeOutput, error)
}

type CreateNodeOptions struct {
	SenderEmail   string `json:"senderEmail"`
	ReceiverEmail string `json:"receiverEmail"`
}

type CreateNodeOutput struct {
	Id string
}

type CourseService interface {
	// UploadCourse provides logic of creating course for selling.
	UploadCourse(ctx context.Context, options *UploadCourseOptions) (*CreateCourseOutput, error)
	GetTeachersList(ctx context.Context, teacherId string) ([]*entity.Course, error)
}

type UploadCourseOptions struct {
	Author         string  `json:"author"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Price          float32 `json:"price"`
	CourseLanguage string  `json:"courseLanguage"`
}

type CreateCourseOutput struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

type CreateGetTeachersListOutput struct {
	Courses []*entity.Course `json:"courses"`
}
