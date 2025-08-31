package support_request

import (
	"forum-support/internal/model"
	forumauthorization "forum-support/pkg/api/forum-authorization"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CreateRequest(supportRequestService model.ISupportRequestService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forumauthorization.AuthorizationData)

		var requestData CreateSupportRequestRequest
		if err := request.Bind(&requestData); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		if len(requestData.Title) > model.MaxTitleLength {
			return echo.NewHTTPError(http.StatusBadRequest, model.ErrTitleTooLong)
		}

		if len(requestData.Description) > model.MaxDescriptionLength {
			return echo.NewHTTPError(http.StatusBadRequest, model.ErrDescriptionTooLong)
		}

		supportRequestID, err := supportRequestService.CreateRequest(
			ctx,
			authorizationData.AccountID,
			requestData.Title,
			requestData.Description,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, CreateSupportRequestResponse{
			SupportRequestID: supportRequestID,
		})
	}
}

func OpenRequest(supportRequestService model.ISupportRequestService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forumauthorization.AuthorizationData)

		requestIDStr := request.Param("requestID")
		requestID, err := strconv.Atoi(requestIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = supportRequestService.OpenRequest(ctx, requestID, authorizationData.Role)
		if err != nil {
			if err == model.ErrActionNotAllowed {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.NoContent(http.StatusOK)
	}
}

func CloseRequest(supportRequestService model.ISupportRequestService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forumauthorization.AuthorizationData)

		requestIDStr := request.Param("requestID")
		requestID, err := strconv.Atoi(requestIDStr)
		if err != nil {
			if err == model.ErrActionNotAllowed {
				return echo.NewHTTPError(http.StatusForbidden, err.Error())
			}
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = supportRequestService.CloseRequest(ctx, requestID, authorizationData.AccountID, authorizationData.Role)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.NoContent(http.StatusOK)
	}
}

func GetRequest(supportRequestService model.ISupportRequestService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		requestIDStr := request.Param("requestID")
		requestID, err := strconv.Atoi(requestIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		supportRequest, err := supportRequestService.GetRequestById(ctx, requestID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, supportRequest)
	}
}

func GetRequests(supportRequestService model.ISupportRequestService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()
		rawStatus := request.QueryParam("status")
		status := model.RequestStatus(rawStatus)

		var err error
		var requests []model.SupportRequest
		if len(status) == 0 {
			requests, err = supportRequestService.GetRequests(ctx)
		} else {
			requests, err = supportRequestService.GetRequestsWithStatus(ctx, status)
		}

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, requests)
	}
}
