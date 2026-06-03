package api

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	adminv1 "github.com/arian-nj/chigame/backend/gen/admin/v1"
)

const testAdminSecret = "test-admin-secret"

func setupTestAppWithAdmin(t *testing.T) *APIApplication {
	t.Helper()
	app := setupTestApp(t)
	app.Config.AdminSecret = testAdminSecret
	return app
}

func TestGetOverview_requiresSecret(t *testing.T) {
	app := setupTestAppWithAdmin(t)
	ctx := context.Background()

	_, err := app.GetOverview(ctx, connect.NewRequest(&adminv1.GetOverviewRequest{}))
	assertConnectCode(t, err, connect.CodePermissionDenied)
}

func TestGetOverview_rejectsWrongSecret(t *testing.T) {
	app := setupTestAppWithAdmin(t)
	ctx := context.Background()

	_, err := app.GetOverview(ctx, connect.NewRequest(&adminv1.GetOverviewRequest{
		AdminSecret: "wrong",
	}))
	assertConnectCode(t, err, connect.CodePermissionDenied)
}

func TestGetOverview_unavailableWhenNotConfigured(t *testing.T) {
	app := setupTestApp(t)
	ctx := context.Background()

	_, err := app.GetOverview(ctx, connect.NewRequest(&adminv1.GetOverviewRequest{
		AdminSecret: "anything",
	}))
	assertConnectCode(t, err, connect.CodeUnavailable)
}

func TestGetOverview_returnsCounts(t *testing.T) {
	app := setupTestAppWithAdmin(t)
	ctx := context.Background()

	_, err := app.GetOrCreateGuestUser(ctx, "admin-test-device")
	if err != nil {
		t.Fatalf("create guest: %v", err)
	}

	resp, err := app.GetOverview(ctx, connect.NewRequest(&adminv1.GetOverviewRequest{
		AdminSecret: testAdminSecret,
	}))
	if err != nil {
		t.Fatalf("GetOverview: %v", err)
	}
	if resp.Msg.GetTotalPersons() < 1 {
		t.Fatalf("total persons = %d, want >= 1", resp.Msg.GetTotalPersons())
	}
	if resp.Msg.GetGuestPersons() < 1 {
		t.Fatalf("guest persons = %d, want >= 1", resp.Msg.GetGuestPersons())
	}
	if len(resp.Msg.GetRecentPersons()) < 1 {
		t.Fatal("expected recent persons")
	}
}
