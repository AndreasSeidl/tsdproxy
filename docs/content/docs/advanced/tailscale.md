---
title: Tailscale
next: /docs/scenarios
---

## Authentication

### OAuth Client (Recommended for Production)

OAuth clients provide the best experience for hands-off deployment. The OAuth client secret never expires and can automatically handle authentication without manual intervention.

{{% steps %}}

#### Create OAuth Client

1. Go to [https://login.tailscale.com/admin/settings/trust-credentials](https://login.tailscale.com/admin/settings/trust-credentials)
2. Click "Credential" button → "OAuth"
3. Select the "Auth keys" scope
4. **Select the tags** you need for your devices (e.g., `tag:server`, `tag:prod`)
5. Click "Generate credential"
6. Copy the Client Secret (you'll only see it once)

>[!IMPORTANT]
> **Tags are required** when using OAuth client secrets. Make note of which tags you selected - you'll need to specify them in the configuration.

>[!TIP]
> OAuth client secrets never expire, making them perfect for production deployments.
> You can create multiple OAuth clients with different tag combinations for different use cases.

#### Add to Configuration

Add the OAuth client secret **and the tags** to your configuration:

```yaml {filename="/config/tsdproxy.yaml"}
tailscale:
  providers:
    default: 
      oauthKey: "tskey-client-xxxxxxxxxxxxx-yyyyyyyyyyyyyyyy"
      oauthTags: ["tag:server"]  # Must match tags selected when creating OAuth client
      oauthKeyFile: "" # alternatively, load from file
```

For multiple tags, add them to the list:

```yaml {filename="/config/tsdproxy.yaml"}
tailscale:
  providers:
    default: 
      oauthKey: "tskey-client-xxxxxxxxxxxxx-yyyyyyyyyyyyyyyy"
      oauthTags: ["tag:server", "tag:prod"]  # Multiple tags
```

Or load from a file for better security:

```yaml {filename="/config/tsdproxy.yaml"}
tailscale:
  providers:
    default: 
      oauthKey: ""
      oauthKeyFile: "/run/secrets/tailscale_oauth" # recommended for Docker secrets
      oauthTags: ["tag:server"]
```

#### Restart

Restart TSDProxy. The OAuth client will automatically handle authentication.

>[!NOTE]
> If both `oauthKey` and `authKey` are set, `oauthKey` takes precedence.

{{% /steps %}}

### Manual OAuth (Browser Authentication)

Manual OAuth authentication mode is enabled when no OAuth key or Auth key is set in the configuration.

{{% steps %}}

#### Configure

Leave both keys empty:

```yaml {filename="/config/tsdproxy.yaml"}
tailscale:
  providers:
    default: 
      oauthKey: ""
      oauthTags: []
      authKey: ""
      authKeyFile: ""
```

#### Authenticate

Go to TSDProxy Dashboard and click on the Proxy that shows "Authentication" status.

>[!TIP]
> Set "Ephemeral" to false in the Tailscale provider to avoid the need of
authentication next time. See [docker Ephemeral label](../../docker/#tsdproxyephemeral)
or [Proxy List configuration](../../list/#proxy-list-file-options)

{{% /steps %}}

### AuthKey

{{% steps %}}

#### Generate Authkey

1. Go to [https://login.tailscale.com/admin/settings/keys](https://login.tailscale.com/admin/settings/keys)
2. Click in "Generate auth key"
3. Add a Description
4. Enable Reusable
5. Enable Ephemeral
6. Add Tags if you need
7. Click in "Generate key"

>[!WARNING]
> If tags were added to the key, all proxies initialized with the same authkey
> will get the same tags.
> Add a new Tailscale provider to the configuration if
> you need to use different)

#### Add to configuration

Add you key to the configuration as follow:

```yaml {filename="/config/tsdproxy.yaml"}
tailscale:
  providers:
    default: 
      oauthKey: ""
      authKey: "GENERATED KEY HERE"
      authKeyFile: ""
```

#### Restart

Restart TSDProxy
gg
{{% /steps %}}

## Funnel

Beside adding the TSDProxy configuration to activate Funnel to a proxy, you also
should give permissions on Tailscale ACL. See [here](../../troubleshooting/#funnel-doesnt-work) to more detail.
