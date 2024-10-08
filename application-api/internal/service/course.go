package service

import (
	"context"
	"fmt"
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

// TODO - need to test this api, all done except testing.
func (a courseService) UploadCourse(ctx context.Context, options *UploadCourseOptions) (*CreateCourseOutput, error) {
	logger := a.logger.
		Named("UploadCourse").
		WithContext(ctx).
		With("options", options)

	// TODO get user and check his user type,
	//userId := requestContext.Get("userId")

	course, err := a.storages.CourseStorage.GetCourse(ctx, &GetCourseFilter{Name: options.Name, Author: options.Author})
	if err != nil {
		logger.Error("failed to get course: ", course)
		return nil, fmt.Errorf("failed to get course: %w", err)
	}
	if course != nil {
		logger.Error("course with such name and author already created: ", course)
		return nil, fmt.Errorf("failed to create a course")
	}

	//TODO set teacher id to course

	//create course
	createdCourse, err := a.storages.CourseStorage.CreateCourse(ctx, course)
	if err != nil {
		logger.Error("failed to create new course: %w", err)
		return nil, fmt.Errorf("failed to create new course: %w", err)
	}
	logger = logger.With("createdCourse", createdCourse)
	logger.Info("successfully created course")
	return &CreateCourseOutput{
		Id:     createdCourse.Id,
		Name:   createdCourse.Name,
		Author: createdCourse.Author,
	}, nil
}
