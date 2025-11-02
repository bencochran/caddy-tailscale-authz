package tailscaleauthz

import (
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func TestMiddlewareUnmarshalCaddyfile_ValidServiceName(t *testing.T) {
	input := `tailscale_authz my_service`
	expected := &Middleware{
		ResourceName: "my_service",
	}

	middleware, err := makeMiddlewareFromCaddyfile(input)
	if err != nil {
		t.Fatalf("failed to parse Caddyfile: %v", err)
	}
	if middleware.ResourceName != expected.ResourceName {
		t.Fatalf("ResourceName mismatch: got %q, want %q", middleware.ResourceName, expected.ResourceName)
	}
}

func TestMiddlewareUnmarshalCaddyfile_NoServiceName(t *testing.T) {
	input := `tailscale_authz`

	middleware, err := makeMiddlewareFromCaddyfile(input)
	if err == nil {
		t.Fatalf("expected error for missing service name, got nil")
	}
	if middleware != nil {
		t.Fatalf("expected nil middleware on error, got %+v", middleware)
	}
	// Expected ArgErr
}

func TestMiddlewareUnmarshalCaddyfile_TooManyArgs(t *testing.T) {
	input := `tailscale_authz service1 service2`

	middleware, err := makeMiddlewareFromCaddyfile(input)
	if err == nil {
		t.Fatalf("expected error for too many arguments, got nil")
	}
	if middleware != nil {
		t.Fatalf("expected nil middleware on error, got %+v", middleware)
	}
	// Expected ArgErr
}

func TestMiddlewareUnmarshalCaddyfile_SpecialCharactersInName(t *testing.T) {
	input := `tailscale_authz my-service_123`
	expected := &Middleware{
		ResourceName: "my-service_123",
	}

	middleware, err := makeMiddlewareFromCaddyfile(input)
	if err != nil {
		t.Fatalf("failed to parse Caddyfile: %v", err)
	}
	if middleware.ResourceName != expected.ResourceName {
		t.Fatalf("ResourceName mismatch: got %q, want %q", middleware.ResourceName, expected.ResourceName)
	}
}

func TestMiddlewareUnmarshalCaddyfile_WildcardServiceName(t *testing.T) {
	input := `tailscale_authz *`
	expectedErr := "wildcard '*' not allowed as resource name"

	middleware, err := makeMiddlewareFromCaddyfile(input)
	if err == nil {
		t.Fatalf("expected error for wildcard resource name, got nil")
	}
	if middleware != nil {
		t.Fatalf("expected nil middleware on error, got %+v", middleware)
	}
	if !startsWith(err.Error(), expectedErr) {
		t.Fatalf("unexpected error message: got %q, want prefix %q", err.Error(), expectedErr)
	}
}

// MARK: - Helpers

func makeMiddlewareFromCaddyfile(input string) (*Middleware, error) {
	dispenser := caddyfile.NewTestDispenser(input)
	middleware := new(Middleware)
	err := middleware.UnmarshalCaddyfile(dispenser)
	if err != nil {
		return nil, err
	}
	return middleware, nil
}
