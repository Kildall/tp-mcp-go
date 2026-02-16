package app

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/fx"
	"tp-mcp-go/internal/client"
)

// RegisterLifecycleHooks registers fx lifecycle hooks
func RegisterLifecycleHooks(lc fx.Lifecycle, c client.Client) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Launch cache initialization in background (non-blocking)
			go func() {
				if err := c.InitializeCache(context.Background()); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to initialize entity type cache: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Placeholder for future cleanup
			return nil
		},
	})
}
