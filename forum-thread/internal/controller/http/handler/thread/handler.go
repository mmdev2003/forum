package thread

import (
	"forum-thread/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateThread(threadService model.IThreadService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body CreateThreadBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		threadID, err := threadService.CreateThread(ctx,
			body.ThreadName,
			body.ThreadDescription,
			body.ThreadColor,
			body.AllowedStatuses,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, CreateThreadResponse{
			ThreadID: threadID,
		})
	}
}

func AllThread(threadService model.IThreadService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		threads, err := threadService.AllThreads(ctx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, AllThreadResponse{
			Threads: threads,
		})
	}
}
