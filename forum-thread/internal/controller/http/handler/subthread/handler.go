package subthread

import (
	"forum-thread/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func CreateSubthread(subthreadService model.ISubthreadService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body CreateSubthreadBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		subthreadID, err := subthreadService.CreateSubthread(ctx,
			body.ThreadID,
			body.ThreadName,
			body.SubthreadName,
			body.SubthreadDescription,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, CreateSubthreadResponse{
			SubthreadID: subthreadID,
		})
	}
}

func AddViewToSubthread(subthreadService model.ISubthreadService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body AddViewToSubthreadBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		err := subthreadService.AddViewToSubthread(ctx,
			body.SubthreadID,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func SubthreadsByThreadID(subthreadService model.ISubthreadService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		threadIDStr := request.Param("threadID")
		threadID, err := strconv.Atoi(threadIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		subthreads, err := subthreadService.SubthreadsByThreadID(ctx,
			threadID,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, SubthreadsByThreadIDResponse{
			Subthreads: subthreads,
		})
	}
}
