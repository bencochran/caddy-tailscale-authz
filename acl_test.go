package tailscaleauthz

import (
	"testing"
)

func TestIsUserAllowedResource(t *testing.T) {
	accessList := AccessList{
		Users: map[string]*UserConfig{
			"bob":     {AllowedResources: []string{"resource1", "resource2"}},
			"alice":   {AllowedResources: []string{"*"}},
			"charlie": {AllowedResources: []string{"resource3"}},
			"amanda":  {AllowedResources: []string{}},
		},
	}
	tests := []struct {
		user     string
		resource string
		allowed  bool
	}{
		{"bob", "resource1", true},
		{"bob", "resource2", true},
		{"bob", "resource3", false},
		{"bob", "unknownresource", false},
		{"alice", "resource1", true},
		{"alice", "resource2", true},
		{"alice", "resource3", true},
		{"alice", "unknownresource", true},
		{"charlie", "resource1", false},
		{"charlie", "resource2", false},
		{"charlie", "resource3", true},
		{"charlie", "unknownresource", false},
		{"amanda", "resource1", false},
		{"amanda", "resource2", false},
		{"amanda", "resource3", false},
		{"amanda", "unknownresource", false},
		{"dave", "resource1", false},
		{"dave", "resource2", false},
		{"dave", "resource3", false},
		{"dave", "unknownresource", false},
	}

	for _, tt := range tests {
		t.Run(tt.user+"_"+tt.resource, func(t *testing.T) {
			result := accessList.IsUserAllowedResource(tt.user, tt.resource)
			if result != tt.allowed {
				t.Errorf("IsUserAllowedResource(%q, %q) = %v; want %v", tt.user, tt.resource, result, tt.allowed)
			}
		})
	}
}

func TestIsUserAllowedResource_EmptyAccessList(t *testing.T) {
	accessList := AccessList{
		Users: map[string]*UserConfig{},
	}
	tests := []struct {
		user     string
		resource string
		allowed  bool
	}{
		{"bob", "resource1", false},
		{"alice", "resource2", false},
		{"charlie", "resource3", false},
		{"amanda", "unknownresource", false},
		{"dave", "resource1", false},
	}

	for _, tt := range tests {
		t.Run(tt.user+"_"+tt.resource, func(t *testing.T) {
			result := accessList.IsUserAllowedResource(tt.user, tt.resource)
			if result != tt.allowed {
				t.Errorf("IsUserAllowedResource(%q, %q) = %v; want %v", tt.user, tt.resource, result, tt.allowed)
			}
		})
	}
}
