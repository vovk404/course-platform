package storage

import (
	"context"
	"github.com/vovk404/course-platform/application-api/internal/entity"
	"github.com/vovk404/course-platform/application-api/internal/service"
	"github.com/vovk404/course-platform/application-api/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type courseStorage struct {
	*database.PostgreSQL
}

var _ service.CourseStorage = (*courseStorage)(nil)

func NewCourseStorage(postgresql *database.PostgreSQL) service.CourseStorage {
	return &courseStorage{postgresql}
}

func (u courseStorage) CreateCourse(ctx context.Context, course *entity.Course) (*entity.Course, error) {
	//TODO somewhy without pointer it throws an error
	err := u.DB.WithContext(ctx).Create(course).Error
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (u courseStorage) GetCourse(ctx context.Context, filter *service.GetCourseFilter) (*entity.Course, error) {
	stmt := u.DB.Preload(clause.Associations)

	if filter.Name != "" {
		stmt = stmt.Where(entity.Course{Name: filter.Name})
	}

	if filter.Author != "" {
		stmt = stmt.Where(entity.Course{Author: filter.Author})
	}

	var course entity.Course
	err := stmt.
		WithContext(ctx).
		First(&course).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &course, nil
}
