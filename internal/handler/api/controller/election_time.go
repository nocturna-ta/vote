package controller

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nocturna-ta/golib/response/rest"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/vote/internal/infrastructures/custresp"
	"github.com/nocturna-ta/vote/internal/usecases/request"
)

// CreateElectionTime godoc
// @Summary Create Election Time
// @Description Create a new election time
// @Tags ElectionTime
// @Accept json
// @Produce json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param request body request.CreateElectionTimeRequest true "Create Election Time Request"
// @Success 200 {object} jsonResponse{data=response.ElectionTimeResponse}
// @Router /v1/election-time [post]
func (api *API) CreateElectionTime(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeController.CreateElectionTime")
	defer span.End()

	var electionTimeRequest request.CreateElectionTimeRequest
	err := json.Unmarshal(req.RawBody(), &electionTimeRequest)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	err = electionTimeRequest.Validate()
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	res, err := api.electionTimeUc.CreateElectionTime(ctx, &electionTimeRequest)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// UpdateElectionTime godoc
// @Summary Update Election Time
// @Description Update an existing election time
// @Tags ElectionTime
// @Accept json
// @Produce json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param id path string true "Election Time ID"
// @Param request body request.UpdateElectionTimeRequest true "Update Election Time Request"
// @Success 200 {object} jsonResponse{data=response.ElectionTimeResponse}
// @Router /v1/election-time/{id} [put]
func (api *API) UpdateElectionTime(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeController.UpdateElectionTime")
	defer span.End()

	electionTimeId, err := uuid.Parse(req.Params("id"))
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	var electionTimeRequest request.UpdateElectionTimeRequest
	err = json.Unmarshal(req.RawBody(), &electionTimeRequest)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	err = electionTimeRequest.Validate()
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	res, err := api.electionTimeUc.UpdateElectionTime(ctx, electionTimeId, &electionTimeRequest)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// DeleteElectionTime godoc
// @Summary Delete Election Time
// @Description Delete an existing election time
// @Tags ElectionTime
// @Accept json
// @Produce json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param id path string true "Election Time ID"
// @Success 200 {object} jsonResponse
// @Router /v1/election-time/{id} [delete]
func (api *API) DeleteElectionTime(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeController.DeleteElectionTime")
	defer span.End()

	electionTimeId, err := uuid.Parse(req.Params("id"))
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	err = api.electionTimeUc.DeleteElectionTime(ctx, electionTimeId)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetMessage("Election time deleted successfully"), nil
}

// GetElectionStatus godoc
// @Summary Get Current Election Status
// @Description Get the current election status
// @Tags ElectionTime
// @Accept json
// @Produce json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Success 200 {object} jsonResponse{data=response.ElectionStatusResponse}
// @Router /v1/election-time/status [get]
func (api *API) GetElectionStatus(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeController.GetElectionStatus")
	defer span.End()

	res, err := api.electionTimeUc.GetCurrentElectionStatus(ctx)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// GetElectionTimeByID godoc
// @Summary Get Election Time by ID
// @Description Get an election time by its ID
// @Tags ElectionTime
// @Accept json
// @Produce json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param id path string true "Election Time ID"
// @Success 200 {object} jsonResponse{data=response.ElectionTimeResponse}
// @Router /v1/election-time/{id} [get]
func (api *API) GetElectionTimeByID(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeController.GetElectionTimeByID")
	defer span.End()

	electionTimeId, err := uuid.Parse(req.Params("id"))
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	res, err := api.electionTimeUc.GetElectionTimeByID(ctx, electionTimeId)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// ActivateElection godoc
// @Summary Manually Activate Election
// @Description Manually activate an election time
// @Tags ElectionTime
// @Accept json
// @Produce json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Param id path string true "Election Time ID"
// @Success 200 {object} jsonResponse{data=response.ElectionTimeResponse}
// @Router /v1/election-time/{id}/activate [post]
func (api *API) ActivateElection(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeController.ActivateElection")
	defer span.End()

	electionTimeId, err := uuid.Parse(req.Params("id"))
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	res, err := api.electionTimeUc.ActivateElection(ctx, electionTimeId)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}

// SyncElectionStatuses godoc
// @Summary Sync Election Statuses
// @Description Sync the statuses of all election times
// @Tags ElectionTime
// @Accept json
// @Produce json
// @Param X-User-Id header string false "Authorized User"
// @Param X-Address header string false "Authorized Address"
// @Param X-Role header string false "Authorized Role"
// @Success 200 {object} jsonResponse
// @Router /v1/election-time/sync [post]
func (api *API) SyncElectionStatuses(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ElectionTimeController.SyncElectionStatuses")
	defer span.End()

	err := api.electionTimeUc.SyncElectionStatuses(ctx)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetMessage("Election statuses synced successfully"), nil
}
