package tailscaleauthz

import (
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	httpcaddyfile.RegisterHandlerDirective("tailscale_authz", parseHandler)
	httpcaddyfile.RegisterDirectiveOrder("tailscale_authz", httpcaddyfile.After, "forward_auth")
}

func parseHandler(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Middleware
	err := m.UnmarshalCaddyfile(h.Dispenser)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler. Syntax:
//
// tailscale_authz <service_name>
func (m *Middleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if !d.Next() {
		return d.ArgErr()
	}

	if !d.NextArg() {
		return d.ArgErr()
	}
	m.ResourceName = d.Val()

	if m.ResourceName == "*" {
		return d.Err("wildcard '*' not allowed as resource name")
	}

	if d.NextArg() {
		return d.ArgErr()
	}

	return nil
}

var (
	_ caddyfile.Unmarshaler = (*Middleware)(nil)
)
