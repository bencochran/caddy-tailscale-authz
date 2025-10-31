package tailscaleauthz

import (
	"github.com/caddyserver/caddy/v2"
)

func init() {
	caddy.RegisterModule(App{})
}

type App struct {
	AccessList AccessList `json:"access_list,omitempty"`
}

func (App) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "tailscale_authz",
		New: func() caddy.Module { return new(App) },
	}
}

func (a *App) Provision(ctx caddy.Context) error {
	return nil
}

func (a *App) Start() error {
	return nil
}

func (a *App) Stop() error {
	return nil
}

var (
	_ caddy.Module      = (*App)(nil)
	_ caddy.Provisioner = (*App)(nil)
)
