package httphandler

import (
	"net/http"

	"github.com/CelticAlreadyUse/Article-accountservices/internal/model"
	"github.com/labstack/echo/v4"
)

func (handler *AccountHandler) requestEmailOTP(e echo.Context) error {
	var body model.OTPRequestGenerateAndSend
	err := e.Bind(&body)
	if err != nil {
	}
	err = validate.Struct(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = handler.otpUsecase.GenerateAndSendOTP(body)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, "Sucessfully sent OTP")
}
func (handler *AccountHandler) confirmEmailOTP(e echo.Context) error {
	var body model.OTPRequestValidate
	err := e.Bind(&body)
	if err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid request body")
	}
	err = validate.Struct(body)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}
	ok, err := handler.otpUsecase.ValidateOTP(&body)
	if !ok {
		return e.JSON(http.StatusUnauthorized, "Invalid OTP")
	}
	if err != nil {
		return e.JSON(http.StatusInternalServerError, "Looks like there have some error")
	}
	err = handler.accountUsecase.SetVerify(e.Request().Context(), body.Email)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, "Verify email failed")
	}
	return e.JSON(http.StatusAccepted, "Sucess your email has been verified")
}
func (handler *AccountHandler) sendOTPPass(e echo.Context) error {
	var data model.OTPRequestGenerateAndSend
	err := e.Bind(&data)
	if err != nil {
		return echo.ErrInternalServerError
	}
	err = handler.otpUsecase.SendOTPPass(data)
	if err != nil {
		return echo.ErrBadRequest
	}
	return e.JSON(http.StatusAccepted, "Sucessfully Send Email")
}
func (handler *AccountHandler) verifyOTPPass(e echo.Context) error {
	var data model.OTPRequestValidate
	err := e.Bind(&data)
	token, err := handler.otpUsecase.ValidateOTPGenerateToken(&data)
	if err != nil {
		return echo.ErrBadRequest
	}
	return e.JSON(http.StatusAccepted, Response{
		AccesToken: token,
	})
}
func (handler *AccountHandler) resetPassword(e echo.Context) error {
	var data model.ResetPasswordReq
	err := e.Bind(&data)
	if err != nil {
		return echo.ErrBadRequest
	}
	err = handler.otpUsecase.ChangePassword(e.Request().Context(), data)
	if err != nil {
		return err
	}
	return e.JSON(http.StatusAccepted, Response{
		Data: "Reset Password Success",
	})
}
