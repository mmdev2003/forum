package payment

import (
	"forum-payment/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func CreatePayment(threadService model.IPaymentService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		var body CreatePaymentBody
		if err := request.Bind(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
		if err := request.Validate(&body); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		amount, address, paymentID, err := threadService.CreatePayment(
			ctx,
			body.AccountID,
			body.ProductType,
			body.Currency,
			body.AmountUSD,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return request.JSON(http.StatusOK, CreatePaymentResponse{
			PaymentID: paymentID,
			Amount:    amount,
			Address:   address,
			Currency:  body.Currency,
		})
	}
}

func PaidPayment(threadService model.IPaymentService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		paymentIDStr := request.Param("paymentID")
		paymentID, _ := strconv.Atoi(paymentIDStr)

		err := threadService.PaidPayment(ctx, paymentID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "OK")
	}
}

func CancelPayment(threadService model.IPaymentService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		paymentIDStr := request.Param("paymentID")
		paymentID, _ := strconv.Atoi(paymentIDStr)

		err := threadService.CancelPayment(ctx, paymentID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, "OK")
	}
}

func StatusPayment(threadService model.IPaymentService) echo.HandlerFunc {
	return func(request echo.Context) error {
		ctx := request.Request().Context()

		paymentIDStr := request.Param("paymentID")
		paymentID, _ := strconv.Atoi(paymentIDStr)

		paymentStatus, err := threadService.StatusPayment(ctx, paymentID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return request.JSON(http.StatusOK, StatusPaymentResponse{
			PaymentStatus: paymentStatus,
		})
	}
}
