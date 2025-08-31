package account_statistic

import (
	"forum-thread/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func CreateAccountStatistic(accountStatisticService model.IAccountStatisticService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body CreateAccountStatisticBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		_, err := accountStatisticService.CreateAccountStatistic(ctx,
			body.AccountID,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func StatisticByAccountID(accountStatisticService model.IAccountStatisticService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		accountIDStr := request.Param("accountID")
		accountID, err := strconv.Atoi(accountIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		accountStatistic, err := accountStatisticService.StatisticByAccountID(ctx, accountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, StatisticByAccountIDResponse{
			AccountStatistic: accountStatistic,
		})
	}
}
