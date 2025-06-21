package controller

import (
	"context"
	"encoding/json"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/response/rest"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/vote/internal/infrastructures/custresp"
	"github.com/nocturna-ta/vote/internal/usecases/request"
)

// GenerateOTP godoc
// @Summary Generate OTP
// @Description Generate a new OTP for the specified voter and purpose
// @Tags OTP
// @Accept json
// @Produce json
// @Param request body request.GenerateOTPRequest true "Generate OTP request"
// @Success 200 {object} jsonResponse{data=response.GenerateOTPResponse}
// @Router /v1/otp/generate [post]
func (api *API) GenerateOTP(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "OTPController.GenerateOTP")
	defer span.End()

	var otpRequest request.GenerateOTPRequest
	err := json.Unmarshal(req.RawBody(), &otpRequest)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	res, err := api.otpUc.GenerateOTP(ctx, &otpRequest)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// VerifyOTP godoc
// @Summary Verify OTP
// @Description Verify the provided OTP code
// @Tags OTP
// @Accept json
// @Produce json
// @Param request body request.VerifyOTPRequest true "Verify OTP request"
// @Success 200 {object} jsonResponse{data=response.VerifyOTPResponse}
// @Router /v1/otp/verify [post]
func (api *API) VerifyOTP(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "OTPController.VerifyOTP")
	defer span.End()

	var otpRequest request.VerifyOTPRequest
	err := json.Unmarshal(req.RawBody(), &otpRequest)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	res, err := api.otpUc.VerifyOTP(ctx, &otpRequest)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// ResendOTP godoc
// @Summary Resend OTP
// @Description Resend OTP for the specified voter and purpose
// @Tags OTP
// @Accept json
// @Produce json
// @Param request body request.ResendOTPRequest true "Resend OTP request"
// @Success 200 {object} jsonResponse{data=response.ResendOTPResponse}
// @Router /v1/otp/resend [post]
func (api *API) ResendOTP(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "OTPController.ResendOTP")
	defer span.End()

	var otpRequest request.ResendOTPRequest
	err := json.Unmarshal(req.RawBody(), &otpRequest)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	res, err := api.otpUc.ResendOTP(ctx, &otpRequest)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetOTPStatus godoc
// @Summary Get OTP Status
// @Description Get the current status of OTP for a voter and purpose
// @Tags OTP
// @Accept json
// @Produce json
// @Param voter_id query string true "Voter ID"
// @Param purpose query string true "OTP Purpose"
// @Success 200 {object} jsonResponse{data=response.OTPStatusResponse}
// @Router /v1/otp/status [get]
func (api *API) GetOTPStatus(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "OTPController.GetOTPStatus")
	defer span.End()

	voterID := req.Query("voter_id")
	purpose := req.Query("purpose")

	if voterID == "" || purpose == "" {
		return custresp.CustomErrorResponse(custerr.ErrChain{
			Message: "voter_id and purpose are required",
			Code:    400,
			Type:    response.ErrBadRequest,
		})
	}

	res, err := api.otpUc.GetOTPStatus(ctx, voterID, purpose)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}
