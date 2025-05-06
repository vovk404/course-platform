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

type getListResponseBody struct {
	*service.CreateGetListOutput
}

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
		routerGroup.GET("/teachers_list", authMiddleware(options), wrapHandler(options, router.getListByTeacherId))
		routerGroup.GET("/list", wrapHandler(options, router.getList))
	}
}

// upload course, only for teacher type of the user
func (a *courseRouter) uploadCourse(requestContext *gin.Context) (interface{}, *httpResponseError) {
	logger := a.logger.Named("uploadCourse").WithContext(requestContext)

	var body createUploadCourseRequestBody
	err := requestContext.ShouldBindJSON(&body)
	if err != nil {
		logger.Info("failed to parse request body", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "invalid request body", Details: err.Error()}
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
		return nil, &httpResponseError{Type: ErrorTypeServer, Message: "failed to create course", Details: err.Error()}
	}
	logger = logger.With("uploadedCourse", uploadedCourse)

	logger.Info("course created successfully")
	return &uploadCourseResponseBody{uploadedCourse}, nil
}

// Get list of courses created by particular teacher.
func (a *courseRouter) getListByTeacherId(requestContext *gin.Context) (interface{}, *httpResponseError) {
	logger := a.logger.Named("getListByTeacherId").WithContext(requestContext)
	userId, ok := requestContext.Value("userId").(string)
	if !ok || userId == "" {
		logger.Error("userId is required and must be a string")
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "userId is required and must be a string"}
	}

	list, err := a.services.CourseService.GetTeachersList(requestContext, userId)
	if list == nil || err != nil {
		logger.Error("failed to get teachers course list", "err", err)
		return nil, &httpResponseError{
			Type:    ErrorTypeClient,
			Message: "failed to get course list",
			Details: err.Error(),
		}
	}

	logger.Info("teachers courses served successfully")
	return &getListResponseBody{
		&service.CreateGetListOutput{list},
	}, nil
}

// Get public course list.
func (a *courseRouter) getList(requestContext *gin.Context) (interface{}, *httpResponseError) {
	logger := a.logger.Named("getList")

	list, err := a.services.CourseService.GetList()
	if list == nil || err != nil {
		logger.Error("failed to get course list", "err", err)
		return nil, &httpResponseError{
			Type:    ErrorTypeClient,
			Message: "failed to get course list",
			Details: err.Error(),
		}
	}

	logger.Info("Courses served successfully")
	return &getListResponseBody{
		&service.CreateGetListOutput{list},
	}, nil
}
