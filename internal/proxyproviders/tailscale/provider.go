// SPDX-FileCopyrightText: 2024 Paulo Almeida <almeidapaulopt@gmail.com>
// SPDX-License-Identifier: MIT
package tailscale

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"tailscale.com/tsnet"

	"github.com/almeidapaulopt/tsdproxy/internal/config"
	"github.com/almeidapaulopt/tsdproxy/internal/proxyconfig"
	"github.com/almeidapaulopt/tsdproxy/internal/proxyproviders"
)

// Client struct implements proxyprovider for tailscale
type Client struct {
	log zerolog.Logger

	Hostname   string
	OAuthKey   string
	OAuthTags  []string
	AuthKey    string
	controlURL string
	datadir    string
}

func New(log zerolog.Logger, name string, provider *config.TailscaleServerConfig) (*Client, error) {
	datadir := filepath.Join(config.Config.Tailscale.DataDir, name)

	return &Client{
		log:        log.With().Str("tailscale", name).Logger(),
		Hostname:   name,
		// make sure the keys are trimmed
		OAuthKey:   strings.TrimSpace(provider.OAuthKey),
		OAuthTags:  provider.OAuthTags,
		AuthKey:    strings.TrimSpace(provider.AuthKey),
		datadir:    datadir,
		controlURL: provider.ControlURL,
	}, nil
}

// NewProxy method implements proxyprovider NewProxy method
func (c *Client) NewProxy(config *proxyconfig.Config) (proxyproviders.ProxyInterface, error) {
	c.log.Debug().
		Str("hostname", config.Hostname).
		Msg("Setting up tailscale server")

	log := c.log.With().Str("Hostname", config.Hostname).Logger()

	// Determine which key to use with the following priority:
	// 1. OAuth key from config (per-proxy)
	// 2. OAuth key from provider (global)
	// 3. Auth key from config (per-proxy)
	// 4. Auth key from provider (global)
	var authKey string
	var advertiseTags []string
	if config.Tailscale.OAuthKey != "" {
		authKey = config.Tailscale.OAuthKey
		advertiseTags = c.OAuthTags
	} else if c.OAuthKey != "" {
		authKey = c.OAuthKey
		advertiseTags = c.OAuthTags
	} else if config.Tailscale.AuthKey != "" {
		authKey = config.Tailscale.AuthKey
	} else {
		authKey = c.AuthKey
	}

	datadir := path.Join(c.datadir, config.Hostname)

	tserver := &tsnet.Server{
		Hostname:      config.Hostname,
		AuthKey:       authKey,
		AdvertiseTags: advertiseTags,
		Dir:           datadir,
		Ephemeral:     config.Tailscale.Ephemeral,
		RunWebClient:  config.Tailscale.RunWebClient,
		UserLogf: func(format string, args ...any) {
			log.Info().Msgf(format, args...)
		},
		Logf: func(format string, args ...any) {
			log.Trace().Msgf(format, args...)
		},

		ControlURL: c.getControlURL(),
	}

	// if verbose is set, use the info log level
	if config.Tailscale.Verbose {
		tserver.Logf = func(format string, args ...any) {
			log.Info().Msgf(format, args...)
		}
	}

	return &Proxy{
		log:      log,
		config:   config,
		tsServer: tserver,
		events:   make(chan proxyproviders.ProxyEvent),
	}, nil
}

// getControlURL method returns the control URL
func (c *Client) getControlURL() string {
	if c.controlURL == "" {
		return proxyconfig.DefaultTailscaleControlURL
	}
	return c.controlURL
}
