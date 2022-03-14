package app

import (
	"context"
	"time"

	"wisdom/internal/wisdom"

	"wisdom/pkg/config"
	"wisdom/pkg/logger"

	"github.com/pkg/errors"
)

// Run runs app. If returned error is not nil, program exited
// unexpectedly and non-zero code should be returned (os.Exit(1) or log.Fatal(...)).
func Run(ctx context.Context, log logger.Logger) error {
	log.Info("staring app")

	var cfg appConfig

	err := config.ParseEnv(&cfg)
	if err != nil {
		return errors.Wrap(err, "can't parse env")
	}

	client := wisdom.NewClient(cfg.WoWHost, cfg.MaxComplexity, log)

	start := time.Now()

	quote, err := client.GetQuote(ctx)
	if err != nil {
		return errors.Wrap(err, "can't get quote")
	}

	log.Infof("after %v we got wise quote: %q", time.Since(start), quote)

	log.Info("app finished")

	return nil
}
