package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"
	adminv1 "github.com/arian-nj/chigame/backend/gen/admin/v1"
)

func (app *APIApplication) GetOverview(
	ctx context.Context,
	req *connect.Request[adminv1.GetOverviewRequest],
) (*connect.Response[adminv1.GetOverviewResponse], error) {
	if err := app.requireAdminSecret(req.Msg.GetAdminSecret()); err != nil {
		return nil, err
	}

	total, err := app.Queries.CountPersons(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("count persons: %w", err))
	}

	guests, err := app.Queries.CountGuestPersons(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("count guests: %w", err))
	}

	recent, err := app.Queries.ListRecentPersons(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("list persons: %w", err))
	}

	summaries := make([]*adminv1.PersonSummary, 0, len(recent))
	for _, p := range recent {
		createdAt := ""
		if p.CreatedAt.Valid {
			createdAt = p.CreatedAt.Time.UTC().Format(time.RFC3339)
		}
		summaries = append(summaries, &adminv1.PersonSummary{
			Id:        p.ID,
			Username:  p.Username,
			IsGuest:   p.IsGuest,
			CreatedAt: createdAt,
		})
	}

	return connect.NewResponse(&adminv1.GetOverviewResponse{
		TotalPersons:   total,
		GuestPersons:   guests,
		RecentPersons:  summaries,
	}), nil
}

func (app *APIApplication) requireAdminSecret(got string) error {
	if app.Config.AdminSecret == "" {
		return connect.NewError(connect.CodeUnavailable, errors.New("admin is not configured"))
	}
	if got != app.Config.AdminSecret {
		return connect.NewError(connect.CodePermissionDenied, errors.New("invalid admin secret"))
	}
	return nil
}
