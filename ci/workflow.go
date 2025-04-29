package ci

import (
	"context"
	"log/slog"
	"os"

	"github.com/NishantBansal2003/go-fuzzstream-lnd/core"
	"github.com/NishantBansal2003/go-fuzzstream-lnd/targets"
	"github.com/NishantBansal2003/go-fuzzstream-lnd/vcs"
)

func ExecuteWorkflow(ctx context.Context, logger *slog.Logger, cfg *core.EngineConfig) {
	if err := vcs.FetchRepositories(ctx, logger, cfg); err != nil {
		logger.Error("Repository fetch failed", "error", err)
		core.PurgeWorkspace(logger)
		os.Exit(1)
	}

	if err := targets.ExecuteTestSuite(ctx, logger, cfg); err != nil {
		logger.Error("Fuzzing operation failed", "error", err)
		core.PurgeWorkspace(logger)
		os.Exit(1)
	}
}
