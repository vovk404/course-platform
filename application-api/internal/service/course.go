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
	userId := ctx.Value("userId").(string)
	user, err := a.storages.UserStorage.GetUser(ctx, &GetUserFilter{UserId: userId})
	if err != nil || user == nil {
		logger.Error("can`t find user with this id", err)
		return nil, fmt.Errorf("can`t find user with this id: %w , error: %w", userId, err)
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
		TeacherId:      user.Id,
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

func (a courseService) GetTeachersList(ctx context.Context, teacherId string) ([]*entity.Course, error) {
	logger := a.logger.
		Named("GetTeachersList").
		WithContext(ctx).
		With("teacherId", teacherId)
	courses, err := a.storages.CourseStorage.GetListByTeacherId(ctx, teacherId)
	if err != nil {
		logger.Error("failed to get courses by teacherId:", err)
	}

	return courses, nil
}
