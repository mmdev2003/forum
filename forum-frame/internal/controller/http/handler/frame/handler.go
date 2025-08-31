package frame

import (
	"forum-frame/internal/model"
	"forum-frame/pkg/api/forum-authorization"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strconv"
)

func CreatePaymentForFrame(frameService model.IFrameService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		var body CreatePaymentForFrameBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		paymentData, err := frameService.CreatePaymentForFrame(ctx,
			body.FrameID,
			authorizationData.AccountID,
			body.Duration,
			body.Currency,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, CreatePaymentForFrameResponse{
			PaymentID: paymentData.PaymentID,
			Address:   paymentData.Address,
			Amount:    paymentData.Amount,
		})
	}
}

func AddNewFrame(frameService model.IFrameService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		file, err := request.FormFile("frameFile")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		src, err := file.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer src.Close()

		frameFile, err := io.ReadAll(src)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		monthlyPriceStr := request.FormValue("monthlyPrice")
		monthlyPrice, err := strconv.ParseFloat(monthlyPriceStr, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		foreverPriceStr := request.FormValue("foreverPrice")
		foreverPrice, err := strconv.ParseFloat(foreverPriceStr, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		name := request.FormValue("name")

		err = frameService.AddNewFrame(ctx,
			frameFile,
			float32(monthlyPrice),
			float32(foreverPrice),
			name,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusCreated, "ok")
	}
}

func ChangeCurrentFrame(frameService model.IFrameService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		dbFrameIDStr := request.Param("dbFrameID")
		dbFrameID, err := strconv.Atoi(dbFrameIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		authorizationData := request.Get(model.AuthorizationDataKey).(*forum_authorization.AuthorizationData)

		err = frameService.ChangeCurrentFrame(ctx, dbFrameID, authorizationData.AccountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, "ok")
	}
}

func ConfirmPaymentForFrame(frameService model.IFrameService, interServerSecretKey string) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body ConfirmPaymentForFrameBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		if body.InterServerSecretKey != interServerSecretKey {
			return echo.NewHTTPError(http.StatusUnauthorized, "inter server secret key is not valid")
		}

		err := frameService.ConfirmPaymentForFrame(ctx, body.PaymentID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "OK")
	}
}

func FramesByAccountID(frameService model.IFrameService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		accountIDStr := request.Param("accountID")
		accountID, _ := strconv.Atoi(accountIDStr)

		frames, currentFrame, err := frameService.FramesByAccountID(ctx, accountID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, FramesByAccountIDResponse{
			Frames:       frames,
			CurrentFrame: currentFrame,
		})
	}
}

func AllFrame(frameService model.IFrameService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		frames, err := frameService.AllFrame(ctx)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, AllFrameResponse{
			Frames: frames,
		})
	}
}

func DownloadFrame(frameService model.IFrameService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		frameIDStr := request.Param("frameID")
		frameID, _ := strconv.Atoi(frameIDStr)

		frameFile, err := frameService.DownloadFrame(ctx, frameID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.Blob(http.StatusOK, "application/octet-stream", frameFile)
	}
}
