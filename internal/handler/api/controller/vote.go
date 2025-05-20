package controller

import (
	"context"
	"encoding/json"
	"github.com/nocturna-ta/golib/response/rest"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/vote/internal/infrastructures/custresp"
	"github.com/nocturna-ta/vote/internal/usecases/request"
)

// CastVote godoc
// @Summary Cast a vote
// @Description Cast a vote
// @Tags Vote
// @Accept json
// @Produce json
// @Param request body request.CastVoteRequest true "Cast vote request"
// @Success 200 {object} jsonResponse{data=response.CastVoterResponse}
// @Router /v1/vote/cast [post]
func (api *API) CastVote(ctx context.Context, req *router.Request) (*rest.JSONResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteController.CastVote")
	defer span.End()

	var voteRequest request.CastVoteRequest
	err := json.Unmarshal(req.RawBody(), &voteRequest)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	res, err := api.voteUc.CastVote(ctx, &voteRequest)
	if err != nil {
		return custresp.CustomErrorResponse(err)
	}

	return rest.NewJSONResponse().SetData(res), nil
}
