package vcs

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/NishantBansal2003/go-fuzzstream-lnd/core"
	"github.com/go-git/go-git/v5"
	"golang.org/x/sync/errgroup"
)

func FetchRepositories(ctx context.Context, logger *slog.Logger, cfg *core.EngineConfig) error {
	logger.Info("Cloning repositories")
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return cloneRepo(ctx, logger, cfg.SourceRepo, "out/project", "source")
	})

	g.Go(func() error {
		return cloneRepo(ctx, logger, cfg.StorageRepo, "out/corpus", "storage")
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("Repository cloning failed: %w", err)
	}
	return nil
}

func cloneRepo(ctx context.Context, logger *slog.Logger, repoURL, path, repoType string) error {
	sanitized := sanitizeURL(repoURL)
	logger.Info("Cloning repository", "type", repoType, "url", sanitized)

	_, err := git.PlainCloneContext(ctx, path, false, &git.CloneOptions{
		URL: repoURL,
	})
	return err
}

func CommitFuzzData(logger *slog.Logger) error {
	repo, err := git.PlainOpen("out/corpus")
	if err != nil {
		return fmt.Errorf("Failed to access repository: %w", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("Failed to get worktree: %w", err)
	}

	if status, _ := wt.Status(); status.IsClean() {
		logger.Info("No corpus changes detected")
		return nil
	}

	if _, err := wt.Add("."); err != nil {
		return fmt.Errorf("Staging failed: %w", err)
	}

	_, err = wt.Commit("Update fuzz corpus", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "fuzz-bot",
			Email: "fuzz@noreply",
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("Commit failed: %w", err)
	}

	return repo.Push(&git.PushOptions{})
}

func sanitizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if u.User != nil {
		u.User = url.User("*****")
	}
	return u.String()
}
