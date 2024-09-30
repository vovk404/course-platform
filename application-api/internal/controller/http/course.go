package http

import (
	"github.com/gin-gonic/gin"
	"github.com/vovk404/course-platform/application-api/internal/service"
	"github.com/vovk404/course-platform/application-api/pkg/errs"
)

type courseRouter struct {
	RouterContext
}

type createUploadCourseRequestBody struct {
	*service.UploadCourseOptions
} // @name createUploadCourseRequestBody

type uploadCourseResponseBody struct {
	*service.CreateCourseOutput
} // @name createAccountResponseBody

type uploadCourseResponseError struct {
	Message string `json:"message"`
	Code    string `json:"code" enums:"user_not_found"`
} // @name uploadCourseResponseError

func (e uploadCourseResponseError) Error() *httpResponseError {
	return &httpResponseError{
		Type:    ErrorTypeClient,
		Message: e.Message,
		Code:    e.Code,
	}
}

func setupCourseRoutes(options RouterOptions) {
	router := &courseRouter{
		RouterContext{
			logger:   options.Logger,
			services: options.Services,
			config:   options.Config,
		},
	}
	routerGroup := options.Handler.Group("/course")
	{
		routerGroup.POST("/new", authMiddleware(options), wrapHandler(options, router.uploadCourse))
	}
}

func (a *courseRouter) uploadCourse(requestContext *gin.Context) (interface{}, *httpResponseError) {
	logger := a.logger.Named("uploadCourse").WithContext(requestContext)

	var body createUploadCourseRequestBody
	err := requestContext.ShouldBindJSON(&body)
	if err != nil {
		logger.Info("failed to parse request body", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "invalid request body", Details: err}
	}
	logger = logger.With("body", body)
	logger.Debug("parsed request body")

	uploadedCourse, err := a.services.CourseService.UploadCourse(requestContext, body.UploadCourseOptions)
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, uploadCourseResponseError{Message: err.Error(), Code: errs.GetCode(err)}.Error()
		}
		logger.Error("failed to create course", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeServer, Message: "failed to create course", Details: err}
	}
	logger = logger.With("uploadedCourse", uploadedCourse)

	logger.Info("course created successfully")
	return &uploadCourseResponseBody{uploadedCourse}, nil
}
