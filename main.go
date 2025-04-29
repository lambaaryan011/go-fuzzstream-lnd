package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"log/slog"
	
	"github.com/NishantBansal2003/go-fuzzstream-lnd/core"
	"github.com/NishantBansal2003/go-fuzzstream-lnd/vcs"
	"github.com/NishantBansal2003/go-fuzzstream-lnd/ci"
)

func manageFuzzCycles(ctx context.Context, logger *slog.Logger, cfg *core.EngineConfig, duration time.Duration) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("Terminating fuzzing operations")
			return
		default:
			cycleCtx, cancel := context.WithCancel(ctx)
			go executeFuzzCycle(cycleCtx, logger, cfg)
			
			select {
			case <-time.After(duration):
				logger.Info("Cycle completion initiated")
				cancel()
				time.Sleep(5 * time.Second)
				performPostCycleTasks(logger)
			case <-ctx.Done():
				cancel()
				logger.Info("Shutdown during active cycle")
				time.Sleep(5 * time.Second)
				performPostCycleTasks(logger)
				return
			}
		}
	}
}

func executeFuzzCycle(ctx context.Context, logger *slog.Logger, cfg *core.EngineConfig) {
	logger.Info("Initiating fuzz cycle", "timestamp", time.Now().Format(time.RFC1123))
	ci.ExecuteWorkflow(ctx, logger, cfg)
}

func performPostCycleTasks(logger *slog.Logger) {
	defer core.PurgeWorkspace(logger)
	if err := vcs.CommitFuzzData(logger); err != nil {
		logger.Error("Repository update failure", "error", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "help" {
		fmt.Println(core.OperationGuide)
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		logger.Info("Termination signal received")
		cancel()
	}()

	if err := core.LoadEnvironment(); err != nil {
		logger.Error("Environment setup failed", "error", err)
		os.Exit(1)
	}

	cfg, err := core.InitEngine()
	if err != nil {
		logger.Error("Configuration initialization error", "error", err)
		os.Exit(1)
	}

	cycleDuration, err := time.ParseDuration(cfg.ExecutionWindow)
	if err != nil {
		logger.Error("Duration parsing error", "input", cfg.ExecutionWindow, "error", err)
		os.Exit(1)
	}

	manageFuzzCycles(ctx, logger, cfg, cycleDuration)
	logger.Info("Process completed")
}
