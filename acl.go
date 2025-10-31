package tailscaleauthz

type AccessList struct {
	Users map[string]*UserConfig `json:"users,omitempty"`
}

type UserConfig struct {
	AllowedResources []string `json:"allowed_resources,omitempty"`
}

func (al *AccessList) IsUserAllowedResource(user string, resource string) bool {
	userConfig, exists := al.Users[user]
	if !exists {
		return false
	}

	for _, allowedResource := range userConfig.AllowedResources {
		if allowedResource == "*" {
			return true
		}
		if allowedResource == resource {
			return true
		}
	}

	return false
}
