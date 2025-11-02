package tailscaleauthz

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func TestMiddleware_ServeHTTP(t *testing.T) {
	tests := []struct {
		name           string
		resourceName   string
		accessList     AccessList
		tailscaleUser  string
		expectedStatus int
	}{
		{
			name:         "authorized user with wildcard",
			resourceName: "service1",
			accessList: AccessList{
				Users: map[string]*UserConfig{
					"user@example.com": {
						AllowedResources: []string{"*"},
					},
				},
			},
			tailscaleUser:  "user@example.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:         "authorized user with specific resource",
			resourceName: "service1",
			accessList: AccessList{
				Users: map[string]*UserConfig{
					"user@example.com": {
						AllowedResources: []string{"service1", "service2"},
					},
				},
			},
			tailscaleUser:  "user@example.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:         "unauthorized user - wrong resource",
			resourceName: "service1",
			accessList: AccessList{
				Users: map[string]*UserConfig{
					"user@example.com": {
						AllowedResources: []string{"service2"},
					},
				},
			},
			tailscaleUser:  "user@example.com",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:         "unauthorized user - not in access list",
			resourceName: "service1",
			accessList: AccessList{
				Users: map[string]*UserConfig{
					"otheruser@example.com": {
						AllowedResources: []string{"service1"},
					},
				},
			},
			tailscaleUser:  "user@example.com",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:         "missing Tailscale-User header",
			resourceName: "service1",
			accessList: AccessList{
				Users: map[string]*UserConfig{
					"user@example.com": {
						AllowedResources: []string{"service1"},
					},
				},
			},
			tailscaleUser:  "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := &Middleware{
				ResourceName: tt.resourceName,
				App: &App{
					AccessList: tt.accessList,
				},
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.tailscaleUser != "" {
				req.Header.Set("Tailscale-User", tt.tailscaleUser)
			}

			rec := httptest.NewRecorder()

			nextHandler := caddyhttp.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
				w.WriteHeader(http.StatusOK)
				return nil
			})

			err := middleware.ServeHTTP(rec, req, nextHandler)

			if err != nil {
				if httpErr, ok := err.(caddyhttp.HandlerError); ok {
					if httpErr.StatusCode != tt.expectedStatus {
						t.Fatalf("expected status %d, got %d", tt.expectedStatus, httpErr.StatusCode)
					}
				} else {
					t.Fatalf("unexpected error type: %v", err)
				}
			} else {
				if tt.expectedStatus != http.StatusOK {
					t.Fatalf("expected status %d, but handler succeeded", tt.expectedStatus)
				}
			}
		})
	}
}
