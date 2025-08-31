package notification

import (
	"forum-notification/internal/model"
	forumauthorization "forum-notification/pkg/api/forum-authorization"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetFiltersByAccountID(filterService model.INotificationFilter) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forumauthorization.AuthorizationData)

		filters, err := filterService.GetFilters(ctx, authorizationData.AccountID)
		if filters == nil && err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, GetFiltersResponse{EnabledFilters: filters})
	}
}

func UpdateFiltersByAccountID(filterService model.INotificationFilter) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forumauthorization.AuthorizationData)
		var filters UpdateFiltersRequest
		if err := request.Bind(&filters); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		err := filterService.UpsertFilters(ctx, authorizationData.AccountID, filters.EnabledFilters)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.NoContent(http.StatusOK)
	}
}
