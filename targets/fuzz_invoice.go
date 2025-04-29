package targets

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/NishantBansal2003/go-fuzzstream-lnd/core"
)

func ExecuteTestSuite(ctx context.Context, logger *slog.Logger, cfg *core.EngineConfig) error {
	for _, pkg := range cfg.TargetPackages {
		select {
		case <-ctx.Done():
			return nil
		default:
			targets, err := locateFuzzTargets(ctx, logger, pkg)
			if err != nil {
				return fmt.Errorf("Target detection failed: %w", err)
			}

			for _, target := range targets {
				if err := runFuzzTest(ctx, logger, pkg, target, cfg); err != nil {
					return fmt.Errorf("Fuzz test failure: %w", err)
				}
			}
		}
	}
	return nil
}

func locateFuzzTargets(ctx context.Context, logger *slog.Logger, pkg string) ([]string, error) {
	pkgPath := filepath.Join("out/project", pkg)
	cmd := exec.CommandContext(ctx, "go", "test", "-list=^Fuzz", ".")
	cmd.Dir = pkgPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Test listing failed: %w", err)
	}

	var targets []string
	for _, line := range strings.Split(string(output), "\n") {
		if trimmed := strings.TrimSpace(line); strings.HasPrefix(trimmed, "Fuzz") {
			targets = append(targets, trimmed)
		}
	}
	return targets, nil
}

func runFuzzTest(ctx context.Context, logger *slog.Logger, pkg, target string, cfg *core.EngineConfig) error {
	cmd := exec.CommandContext(ctx, "go", "test",
		fmt.Sprintf("-fuzz=^%s$", target),
		fmt.Sprintf("-test.fuzzcachedir=%s", filepath.Join("out/corpus", pkg)),
		fmt.Sprintf("-fuzztime=%s", cfg.ExecutionWindow),
		fmt.Sprintf("-parallel=%d", cfg.ParallelJobs),
	)
	cmd.Dir = filepath.Join("out/project", pkg)

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Process start failed: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go streamLogs(logger.With("target", target), "stdout", &wg, stdout)
	go streamLogs(logger.With("target", target), "stderr", &wg, stderr)
	wg.Wait()

	return cmd.Wait()
}

func streamLogs(logger *slog.Logger, stream string, wg *sync.WaitGroup, r io.Reader) {
	defer wg.Done()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		logger.Info("Fuzzer output", "stream", stream, "data", scanner.Text())
	}
}
