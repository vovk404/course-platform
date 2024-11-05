package service

import (
	"context"
	"fmt"
	"github.com/vovk404/course-platform/application-api/internal/entity"
)

type courseService struct {
	serviceContext
}

var _ CourseService = (*courseService)(nil)

func NewCourseService(options *Options) CourseService {
	return &courseService{
		serviceContext: serviceContext{
			storages: options.Storages,
			config:   options.Config,
			logger:   options.Logger.Named("CourseService"),
		},
	}
}

func (a courseService) UploadCourse(ctx context.Context, options *UploadCourseOptions) (*CreateCourseOutput, error) {
	logger := a.logger.
		Named("UploadCourse").
		WithContext(ctx).
		With("options", options)

	course, err := a.storages.CourseStorage.GetCourse(ctx, &GetCourseFilter{Name: options.Name, Author: options.Author})
	if err != nil {
		logger.Error("failed to get course: ", course)
		return nil, fmt.Errorf("failed to get course: %w", err)
	}
	if course != nil {
		logger.Error("course with such name and author already created: ", course)
		return nil, fmt.Errorf("course with such name and author already created")
	}
	//get user
	user, err := a.storages.UserStorage.GetUser(ctx, &GetUserFilter{UserId: options.TeacherId})
	if err != nil || user == nil {
		logger.Error("can`t find user with this id", err)
		return nil, fmt.Errorf("can`t find user with this id: %w , error: %w", options.TeacherId, err)
	}
	if user.Type != entity.Teacher {
		return nil, fmt.Errorf("user`s type can`t allow creating a course")
	}

	insertCourse := entity.Course{
		Name:           options.Name,
		Author:         options.Author,
		Description:    options.Description,
		Price:          options.Price,
		CourseLanguage: options.CourseLanguage,
		TeacherId:      options.TeacherId,
	}
	//create course
	createdCourse, err := a.storages.CourseStorage.CreateCourse(ctx, &insertCourse)
	if err != nil {
		logger.Error("failed to create a new course: %w", err)
		return nil, fmt.Errorf("failed to create a new course: %w", err)
	}
	logger = logger.With("createdCourse", createdCourse)
	logger.Info("successfully created course")
	return &CreateCourseOutput{
		Id:     createdCourse.Id,
		Name:   createdCourse.Name,
		Author: createdCourse.Author,
	}, nil
}
