package tailscaleauthz

import (
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
)

func init() {
	httpcaddyfile.RegisterGlobalOption("tailscale_authz", parseApp)
}

func parseApp(d *caddyfile.Dispenser, _ any) (any, error) {
	app := new(App)
	err := app.UnmarshalCaddyfile(d)
	if err != nil {
		return nil, err
	}
	return httpcaddyfile.App{
		Name:  "tailscale_authz",
		Value: caddyconfig.JSON(app, nil),
	}, nil
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler. Syntax:
//
//	tailscale_authz {
//	    user <identifier> <resource...>
//	    ...
//	}
func (a *App) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	accessList := AccessList{}
	accessList.Users = make(map[string]*UserConfig)

	if !d.Next() {
		return d.ArgErr()
	}

	for d.NextBlock(0) {
		switch d.Val() {
		case "user":
			if !d.NextArg() {
				return d.ArgErr()
			}
			userId := d.Val()

			if accessList.Users[userId] != nil {
				return d.Errf("user %s already defined", userId)
			}

			allowedResources := d.RemainingArgs()
			if len(allowedResources) == 0 {
				return d.Errf("no resources specified for user %s", userId)
			}

			hasWildcard := false
			for _, res := range allowedResources {
				if res == "*" {
					hasWildcard = true
					break
				}
			}

			if hasWildcard && len(allowedResources) > 1 {
				return d.Err("cannot combine wildcard '*' with other resources")
			}

			accessList.Users[userId] = &UserConfig{
				AllowedResources: allowedResources,
			}
		default:
			return d.Errf("unrecognized directive: %s", d.Val())
		}
	}

	a.AccessList = accessList
	return nil
}

var (
	_ caddyfile.Unmarshaler = (*App)(nil)
)
