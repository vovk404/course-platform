package storage

import (
	"context"
	"github.com/vovk404/course-platform/application-api/internal/entity"
	"github.com/vovk404/course-platform/application-api/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userStorage struct {
	*database.PostgreSQL
}

var _ UserStorage = (*userStorage)(nil)

func NewUserStorage(postgresql *database.PostgreSQL) UserStorage {
	return &userStorage{postgresql}
}

func (u *userStorage) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := u.DB.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userStorage) GetUser(ctx context.Context, filter *GetUserFilter) (*entity.User, error) {
	stmt := u.DB.Preload(clause.Associations)

	if filter.Email != "" {
		stmt = stmt.Where(entity.User{Email: filter.Email})
	}

	if filter.UserId != "" {
		stmt = stmt.Where(entity.User{Id: filter.UserId})
	}

	var user entity.User
	err := stmt.
		WithContext(ctx).
		First(&user).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
