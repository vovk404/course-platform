package service

import (
	"context"
	"github.com/vovk404/course-platform/application-api/internal/entity"
)

type Storages struct {
	UserStorage    UserStorage
	AccountStorage AccountStorage
	NodeStorage    NodeStorage
	CourseStorage  CourseStorage
}

type UserStorage interface {
	// GetUser provides getting user from storage via requested filters.
	GetUser(ctx context.Context, filter *GetUserFilter) (*entity.User, error)
	// CreateUser provides creating user in the system.
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
}

type GetUserFilter struct {
	Email  string
	UserId string
}

type AccountStorage interface {
	// CreateAccount provides creating account in the system.
	CreateAccount(ctx context.Context, account *entity.Account) (*entity.Account, error)
	// GetAccount provides logic of getting account from storage.
	GetAccount(ctx context.Context, filter *GetAccountFilter) (*entity.Account, error)
	// UpdateAccount provides logic of updating account in storage.
	UpdateAccount(ctx context.Context, account *entity.Account) (*entity.Account, error)
}

type GetAccountFilter struct {
	AccountId string
	UserId    string
}

type NodeStorage interface {
	// CreateNode provides creating new node in system.
	CreateNode(ctx context.Context, node *entity.Node) (*entity.Node, error)
}

type CourseStorage interface {
	GetCourse(ctx context.Context, filter *GetCourseFilter) (*entity.Course, error)
	// CreateCourse provides creating course in the system.
	CreateCourse(ctx context.Context, course *entity.Course) (*entity.Course, error)
}

type GetCourseFilter struct {
	Name   string
	Author string
}
