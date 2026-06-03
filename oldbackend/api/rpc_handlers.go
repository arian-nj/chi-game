package api

import (
	"context"

	"connectrpc.com/connect"
	healthcheckv1 "github.com/arian-nj/chigame/backend/gen/healthcheck/v1"
)

// HealthCheck implements [healthcheckv1connect.HealthcheckServiceHandler].
func (app *APIApplication) HealthCheck(ctx context.Context,req *connect.Request[healthcheckv1.HealthRequest]) (*connect.Response[healthcheckv1.HealthResponse], error) {
	return connect.NewResponse(&healthcheckv1.HealthResponse{
		HealthType: healthcheckv1.HealthType_HEALTH_TYPE_OK,
	}), nil
}
