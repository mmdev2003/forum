package status

import (
	"forum-status/internal/model"
	"forum-status/pkg/api/forum-authorization"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func CreatePaymentForStatus(statusService model.IStatusService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		var body CreatePaymentForStatusBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		paymentData, err := statusService.CreatePaymentForStatus(ctx,
			body.StatusID,
			authorizationData.AccountID,
			body.Duration,
			body.Currency,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, CreatePaymentForStatusResponse{
			PaymentID: paymentData.PaymentID,
			Address:   paymentData.Address,
			Amount:    paymentData.Amount,
		})
	}
}

func ConfirmPaymentForStatus(statusService model.IStatusService, interServerSecretKey string) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body ConfirmPaymentForStatusBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if body.InterServerSecretKey != interServerSecretKey {
			return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect inter server secret key")
		}

		err := statusService.ConfirmPaymentForStatus(ctx, body.PaymentID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "OK")
	}
}

func RevokeStatusFromAccount(statusService model.IStatusService, interServerSecretKey string) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body RevokeStatusFromAccountBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if body.InterServerSecretKey != interServerSecretKey {
			return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect inter server secret key")
		}

		err := statusService.RevokeStatus(ctx, body.StatusID, body.AccountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, "ok")
	}
}

func AssignStatusToAccount(statusService model.IStatusService, interServerSecretKey string) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body AssignStatusToAccountBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if body.InterServerSecretKey != interServerSecretKey {
			return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect inter server secret key")
		}

		err := statusService.AssignStatusToAccount(ctx, body.StatusID, body.AccountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, "OK")
	}
}

func GetStatusByAccountID(statusService model.IStatusService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		accountIDStr := request.Param("accountID")
		accountID, _ := strconv.Atoi(accountIDStr)

		accountStatuses, err := statusService.StatusesByAccountID(ctx, accountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, GetStatusByAccountIDResponse{
			Statuses: accountStatuses,
		})
	}
}

func GetAllStatus(statusService model.IStatusService) echo.HandlerFunc {
	var statuses []*model.StatusConst
	for _, status := range model.Statuses {
		statuses = append(statuses, &status)
	}
	return func(request echo.Context) error {
		return request.JSON(http.StatusOK, GetAllStatusResponse{
			Statuses: statuses,
		})
	}
}
