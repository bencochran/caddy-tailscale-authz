package tailscaleauthz

import (
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func TestUnmarshalCaddyfile_ValidSingleUser(t *testing.T) {
	input := `tailscale_authz {
		user bob resource1 resource2
	}
	`
	expected := &App{
		AccessList: AccessList{
			Users: map[string]*UserConfig{
				"bob": {
					AllowedResources: []string{"resource1", "resource2"},
				},
			},
		},
	}

	app, err := makeAppFromCaddyfile(input)
	if err != nil {
		t.Fatalf("failed to parse Caddyfile: %v", err)
	}
	compareAccessLists(t, app.AccessList, expected.AccessList)
}

func TestUnmarshalCaddyfile_ValidMultipleUsers(t *testing.T) {
	input := `tailscale_authz {
		user bob resource1 resource2
		user alice resource3
	}
	`
	expected := &App{
		AccessList: AccessList{
			Users: map[string]*UserConfig{
				"bob": {
					AllowedResources: []string{"resource1", "resource2"},
				},
				"alice": {
					AllowedResources: []string{"resource3"},
				},
			},
		},
	}

	app, err := makeAppFromCaddyfile(input)
	if err != nil {
		t.Fatalf("failed to parse Caddyfile: %v", err)
	}
	compareAccessLists(t, app.AccessList, expected.AccessList)
}

func TestUnmarshalCaddyfile_ValidWildcard(t *testing.T) {
	input := `tailscale_authz {
		user alice *
	}
	`
	expected := &App{
		AccessList: AccessList{
			Users: map[string]*UserConfig{
				"alice": {
					AllowedResources: []string{"*"},
				},
			},
		},
	}

	app, err := makeAppFromCaddyfile(input)
	if err != nil {
		t.Fatalf("failed to parse Caddyfile: %v", err)
	}
	compareAccessLists(t, app.AccessList, expected.AccessList)
}

func TestUnmarshalCaddyfile_DuplicateUser(t *testing.T) {
	input := `tailscale_authz {
		user bob resource1
		user bob resource2
	}
	`
	expectedErr := "user bob already defined"

	app, err := makeAppFromCaddyfile(input)
	if err == nil {
		t.Fatalf("expected error for duplicate user definition, got nil")
	}
	if app != nil {
		t.Fatalf("expected nil app on error, got %+v", app)
	}
	if !startsWith(err.Error(), expectedErr) {
		t.Fatalf("unexpected error message: got %q, want prefix %q", err.Error(), expectedErr)
	}
}

func TestUnmarshalCaddyfile_UserWithNoResources(t *testing.T) {
	input := `tailscale_authz {
		user bob
	}
	`
	expectedErr := "no resources specified for user bob"

	app, err := makeAppFromCaddyfile(input)
	if err == nil {
		t.Fatalf("expected error for user with no resources, got nil")
	}
	if app != nil {
		t.Fatalf("expected nil app on error, got %+v", app)
	}
	if !startsWith(err.Error(), expectedErr) {
		t.Fatalf("unexpected error message: got %q, want prefix %q", err.Error(), expectedErr)
	}
}

func TestUnmarshalCaddyfile_WildcardWithOtherResources(t *testing.T) {
	input := `tailscale_authz {
		user bob * resource1
	}
	`
	expectedErr := "cannot combine wildcard '*' with other resources"

	app, err := makeAppFromCaddyfile(input)
	if err == nil {
		t.Fatalf("expected error for wildcard with other resources, got nil")
	}
	if app != nil {
		t.Fatalf("expected nil app on error, got %+v", app)
	}
	if !startsWith(err.Error(), expectedErr) {
		t.Fatalf("unexpected error message: got %q, want prefix %q", err.Error(), expectedErr)
	}
}

func TestUnmarshalCaddyfile_UnrecognizedDirective(t *testing.T) {
	input := `tailscale_authz {
		unknown_directive value
	}
	`
	expectedErr := "unrecognized directive: unknown_directive"

	app, err := makeAppFromCaddyfile(input)
	if err == nil {
		t.Fatalf("expected error for unrecognized directive, got nil")
	}
	if app != nil {
		t.Fatalf("expected nil app on error, got %+v", app)
	}
	if !startsWith(err.Error(), expectedErr) {
		t.Fatalf("unexpected error message: got %q, want prefix %q", err.Error(), expectedErr)
	}
}

func TestUnmarshalCaddyfile_EmptyBlock(t *testing.T) {
	input := `tailscale_authz {
	}
	`
	expected := &App{
		AccessList: AccessList{
			Users: map[string]*UserConfig{},
		},
	}

	app, err := makeAppFromCaddyfile(input)
	if err != nil {
		t.Fatalf("failed to parse Caddyfile: %v", err)
	}
	compareAccessLists(t, app.AccessList, expected.AccessList)
}

func TestUnmarshalCaddyfile_NoBlock(t *testing.T) {
	// When there's no block, just the directive name, it's treated the same as an empty block
	// This test verifies that behavior is consistent with EmptyBlock test
	input := `tailscale_authz`
	expected := &App{
		AccessList: AccessList{
			Users: map[string]*UserConfig{},
		},
	}

	app, err := makeAppFromCaddyfile(input)
	if err != nil {
		t.Fatalf("failed to parse Caddyfile: %v", err)
	}
	compareAccessLists(t, app.AccessList, expected.AccessList)
}

// MARK: - Helpers

func makeDispenser(input string) *caddyfile.Dispenser {
	return caddyfile.NewTestDispenser(input)
}

func makeAppFromCaddyfile(input string) (*App, error) {
	dispenser := makeDispenser(input)
	app := new(App)
	err := app.UnmarshalCaddyfile(dispenser)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func compareAccessLists(t *testing.T, got, want AccessList) {
	t.Helper()
	if len(got.Users) != len(want.Users) {
		t.Fatalf("number of users mismatch: got %d, want %d", len(got.Users), len(want.Users))
	}
	for userId, wantConfig := range want.Users {
		gotConfig, exists := got.Users[userId]
		if !exists {
			t.Fatalf("missing user %s in got access list", userId)
		}
		if len(gotConfig.AllowedResources) != len(wantConfig.AllowedResources) {
			t.Fatalf("number of allowed resources mismatch for user %s: got %d, want %d", userId, len(gotConfig.AllowedResources), len(wantConfig.AllowedResources))
		}
		for i, wantRes := range wantConfig.AllowedResources {
			if gotConfig.AllowedResources[i] != wantRes {
				t.Fatalf("allowed resource mismatch for user %s at index %d: got %s, want %s", userId, i, gotConfig.AllowedResources[i], wantRes)
			}
		}
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
