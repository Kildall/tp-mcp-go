package app

import (
	"go.uber.org/fx"
	"tp-mcp-go/internal/client"
	"tp-mcp-go/internal/client/auth"
	"tp-mcp-go/internal/config"
)

var Module = fx.Options(
	fx.Provide(config.Load),
	fx.Provide(func(cfg *config.Config) auth.Strategy {
		return auth.NewAccessTokenStrategy(cfg.AccessToken)
	}),
	fx.Provide(func(cfg *config.Config, authStrategy auth.Strategy) client.Client {
		return client.NewHTTPClient(cfg, authStrategy)
	}),
	fx.Invoke(RegisterLifecycleHooks),
)
