package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vovk404/course-platform/application-api/internal/entity"
	"github.com/vovk404/course-platform/application-api/internal/service"
	"github.com/vovk404/course-platform/application-api/pkg/errs"
)

type accountRouter struct {
	RouterContext
}

func setupAccountRoutes(options RouterOptions) {
	router := &accountRouter{
		RouterContext{
			logger:   options.Logger,
			services: options.Services,
			config:   options.Config,
		},
	}

	routerGroup := options.Handler.Group("/account")
	{
		routerGroup.POST("", authMiddleware(options), wrapHandler(options, router.createAccount))
		routerGroup.GET("/:id", authMiddleware(options), wrapHandler(options, router.getAccount))
		routerGroup.PATCH("/:id", authMiddleware(options), wrapHandler(options, router.updateAccount))
	}
}

type createAccountRequestBody struct {
	*service.CreateAccountOptions
} // @name createAccountRequestBody

type createAccountResponseBody struct {
	*service.CreateAccountOutput
} // @name createAccountResponseBody

type createAccountResponseError struct {
	Message string `json:"message"`
	Code    string `json:"code" enums:"user_not_found"`
} // @name createAccountResponseError

func (e createAccountResponseError) Error() *httpResponseError {
	return &httpResponseError{
		Type:    ErrorTypeClient,
		Message: e.Message,
		Code:    e.Code,
	}
}

// @id           CreateAccount
// @Summary      Creates account.
// @Accept       application/json
// @Produce      application/json
// @Param        fields body createAccountRequestBody true "data"
// @Success      200 {object} createAccountResponseBody
// @Failure      422,500 {object} createAccountResponseError
// @Router       /account [POST]
func (a *accountRouter) createAccount(requestContext *gin.Context) (interface{}, *httpResponseError) {
	logger := a.logger.Named("createAccount").WithContext(requestContext)

	var body createAccountRequestBody
	err := requestContext.ShouldBindJSON(&body)
	if err != nil {
		logger.Info("failed to parse request body", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "invalid request body", Details: err}
	}
	logger = logger.With("body", body)
	logger.Debug("parsed request body")

	createdAccount, err := a.services.AccountService.CreateAccount(requestContext, body.CreateAccountOptions)
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, createAccountResponseError{Message: err.Error(), Code: errs.GetCode(err)}.Error()
		}
		logger.Error("failed to create account", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeServer, Message: "failed to create account", Details: err}
	}
	logger = logger.With("createdAccount", createdAccount)

	logger.Info("account created successfully")
	return &createAccountResponseBody{createdAccount}, nil
}

type getAccountResponseBody struct {
	*entity.Account
} // @name getAccountResponseBody

type getAccountResponseError struct {
	Message string `json:"message"`
	Code    string `json:"code" enums:"user_not_found"`
} // @name getAccountResponseError

func (e getAccountResponseError) Error() *httpResponseError {
	return &httpResponseError{
		Type:    ErrorTypeClient,
		Message: e.Message,
		Code:    e.Code,
	}
}

// @id           GetAccount
// @Summary      Gets account.
// @Accept       application/json
// @Produce      application/json
// @Param        id path string true "Account ID"
// @Success      200 {object} getAccountResponseBody
// @Failure      422,500 {object} getAccountResponseError
// @Router       /account/{id} [GET]
func (a *accountRouter) getAccount(requestContext *gin.Context) (interface{}, *httpResponseError) {
	logger := a.logger.Named("getAccount").WithContext(requestContext)

	accountId := requestContext.Param("id")
	if _, ok := uuid.Parse(accountId); ok != nil {
		logger.Info("invalid account id parameter", "param", accountId)
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "invalid account id parameter"}
	}
	logger = logger.With("accountId", accountId)
	logger.Debug("parsed params")

	requestUserId := requestContext.Value("userId")
	if requestUserId == nil {
		logger.Info("user not found")
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "user not found"}
	}
	logger = logger.With("userId", requestUserId)
	logger.Debug("got userId")

	userId := fmt.Sprint(requestUserId)
	if _, ok := uuid.Parse(userId); ok != nil {
		logger.Info("invalid user id parameter")
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "invalid user id parameter"}
	}
	logger = logger.With("userId", userId)
	logger.Debug("validated uuid userId")

	account, err := a.services.AccountService.GetAccount(requestContext, &service.GetAccountOptions{AccountId: accountId, UserId: userId})
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, getAccountResponseError{Message: err.Error(), Code: errs.GetCode(err)}.Error()
		}
	}
	logger = logger.With("account", account)

	logger.Info("successfully got account")
	return getAccountResponseBody{account}, nil
}

type updateAccountResponseBody struct {
	*entity.Account
} // @name updateAccountResponseBody

// @id           UpdateAccount
// @Summary      Updates account entity account.
// @Accept       application/json
// @Produce      application/json
// @Param        id path string true "Account ID"
// @Success      200 {object} updateAccountResponseBody
// @Failure      422,500 {object} httpResponseError
// @Router       /account/{id} [PATCH]
func (a *accountRouter) updateAccount(requestContext *gin.Context) (interface{}, *httpResponseError) {
	logger := a.logger.Named("updateAccount").WithContext(requestContext)

	accountId := requestContext.Param("id")
	if _, ok := uuid.Parse(accountId); ok != nil {
		logger.Info("invalid account id parameter", "param", accountId)
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "invalid account id parameter"}
	}
	logger = logger.With("accountId", accountId)
	logger.Debug("parsed params")

	requestUserId := requestContext.Value("userId")
	if requestUserId == nil {
		logger.Info("user not found")
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "user not found"}
	}
	logger = logger.With("userId", requestUserId)
	logger.Debug("got userId")

	userId := fmt.Sprint(requestUserId)
	if _, ok := uuid.Parse(userId); ok != nil {
		logger.Info("invalid user id parameter")
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "invalid user id parameter"}
	}
	logger = logger.With("userId", userId)
	logger.Debug("validated uuid userId")

	account, err := a.services.AccountService.GetAccount(requestContext, &service.GetAccountOptions{AccountId: accountId, UserId: userId})
	if err != nil {
		if errs.IsExpected(err) {
			logger.Info(err.Error())
			return nil, getAccountResponseError{Message: err.Error(), Code: errs.GetCode(err)}.Error()
		}
	}
	logger = logger.With("account", account)

	err = requestContext.ShouldBindJSON(&account)
	if err != nil {
		logger.Info("failed to parse request body", "err", err)
		return nil, &httpResponseError{Type: ErrorTypeClient, Message: "invalid request body", Details: err}
	}
	logger = logger.With("account", account)
	logger.Debug("parsed request body")

	updatedAccount, err := a.services.AccountService.UpdateAccount(requestContext, account)
	if err != nil {
		logger.Error("failed to update account: ", err)
		return nil, &httpResponseError{Type: ErrorTypeServer, Message: "failed to update account", Details: err}
	}
	logger = logger.With("updatedAccount", updatedAccount)

	logger.Info("successfully updated account")
	return updateAccountResponseBody{updatedAccount}, nil
}
