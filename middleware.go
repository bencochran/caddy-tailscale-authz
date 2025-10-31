package tailscaleauthz

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(Middleware{})
}

type Middleware struct {
	ResourceName string
	App          *App
}

func (Middleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.tailscale_authz",
		New: func() caddy.Module { return new(Middleware) },
	}
}

func (m *Middleware) Provision(ctx caddy.Context) error {
	appIface, err := ctx.App("tailscale_authz")
	if err != nil {
		return err
	}
	app, ok := appIface.(*App)
	if !ok {
		return fmt.Errorf("tailscale_authz app is not of type *App")
	}
	m.App = app
	return nil
}

func (m *Middleware) Validate() error {
	if m.ResourceName == "" {
		return fmt.Errorf("<resource> is required")
	}
	return nil
}

func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// TODO: Implement
	return next.ServeHTTP(w, r)
}

var (
	_ caddy.Module                = (*Middleware)(nil)
	_ caddy.Provisioner           = (*Middleware)(nil)
	_ caddy.Validator             = (*Middleware)(nil)
	_ caddyhttp.MiddlewareHandler = (*Middleware)(nil)
)
