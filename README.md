# caddy-tailscale-authz

> [!CAUTION]
> This module is very experimental. It may not work. It may not provide any security. It may break all on its own or because I make breaking changes. You almost certainly shouldn’t use it.

A caddy middleware that authorization checks for Tailscale-authenticated requests. It works in conjunction with [tailscale-nginx-auth](https://caddyserver.com/docs/caddyfile/directives/forward_auth#tailscale) forward auth proxy or [caddy-tailscale](https://github.com/tailscale/caddy-tailscale), adding a minimal ACL layer on top of the identity information provided by those.

It doesn’t go groups or anything fancy. Just the minimum to restrict access to certain tailscale users.

## Why?

Scenario: you have a homelab with a bunch of apps behind a reverse proxy. You want to share a couple of those apps with someone. You can invite them to your tailnet and ACL them out of accessing anything but your reverse proxy. Or share your reverse proxy machine with them on on their own tailnet. But in either scenario they have defacto access to every app served by that reverse proxy. And thus, this plugin.

## Example

```caddyfile
{
    tailscale_authz {
        user alice@example.com *
        user bob@example.com dashboard
    }

    # A snipped using the nginx-auth forward proxy
    (tailscale_authn) {
        forward_auth unix//run/tailscale.nginx-auth.sock {
            uri /auth
            header_up Remote-Addr {remote_host}
            header_up Remote-Port {remote_port}
            copy_headers Tailscale-User
        }
    }

    dashboard.example.com {
        import tailscale_authn
        tailscale_authz dashboard

        reverse_proxy http://dashboard.internal:8080
    }

    admin.example.com {
        import tailscale_authn
        tailscale_authz admin

        reverse_proxy http://admin.internal:3000
    }
}

```

In this example, both `dashboard.example.com` and `admin.example.com` would require clients to be connecting via Tailscale. Alice would have access to both sites, and Bob would only have access to `dashboard.example.com`. Any other users on this tailnet would not have access to either.

## Building

Use [xcaddy](https://github.com/caddyserver/xcaddy) to build Caddy with caddy-tailscale-authz:

```bash
xcaddy build v2.10.2 --with github.com/bencochran/caddy-tailscale-authz
```

If using Docker, see [Adding custom Caddy modules](https://hub.docker.com/_/caddy#adding-custom-caddy-modules) from the Caddy docker docs.

## Syntax

#### Global option

```caddyfile
{
    tailscale_authz {
        user <identifier> <resources...>
    }
}
```

* **identifier** is the target user’s Tailscale identifier (typically their email, [but occasionally `@github` or `@passkey`](https://tailscale.com/kb/1337/policy-syntax#reference-users0)). This can be a full user of your tailnet, or a user who has access to this machine through [sharing](https://tailscale.com/kb/1084/sharing).
* **<resource...>** is a list of resources accessible to this user. These are simple string identifiers internal to your Caddyfile (just something to reference in a site directive). To allow a user access to all resources use `*` instead of a list of resources.

#### Site directive

```caddyfile
tailscale_authz <resource>
```

* `resource` is a simple string identifying this resource, as referenced in the above global option

### Usage

caddy-tailscale-authz requires a `Tailscale-User` header to be populated by an authentication layer prior to `tailscale_authz` taking a pass at a request. The two primary ways to do that are via `forward_auth` to [tailscale-nginx-auth](https://tailscale.com/blog/tailscale-auth-nginx) or the [caddy-tailscale](https://github.com/tailscale/caddy-tailscale) module.

#### Using with [tailscale-nginx-auth](https://tailscale.com/blog/tailscale-auth-nginx)

[As described in the Caddy docs](https://caddyserver.com/docs/caddyfile/directives/forward_auth#tailscale), `forward_auth` can be used to authenticate via an existing Tailscale connection.

```caddyfile
example.com {
    forward_auth unix//run/tailscale.nginx-auth.sock {
        uri /auth
        header_up Remote-Addr {remote_host}
        header_up Remote-Port {remote_port}
        copy_headers Tailscale-User
    }

    tailscale_authz example

    # ...
}
```


#### Using with [caddy-tailscale](https://github.com/tailscale/caddy-tailscale)

The `tailscale_auth` directive from caddy-tailscale doesn’t populate the `Tailscale-User` header so we must hoist the `user.tailscale_user` value to a header ourselves.

```caddyfile
# TODO: Add caddy-tailscale example
```


## HTTP Status Codes

- **401 Unauthorized**: Returned when Tailscale authentication headers are missing
- **403 Forbidden**: Returned when the user is authenticated but not authorized for the service
- **Continues normally**: When the user is both authenticated and authorized
