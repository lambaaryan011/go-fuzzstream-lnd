package core

import (
	"context"
	"log/slog"
	"os"
)

const OperationGuide = `Usage: go run main.go [command]

Commands:
  help      Show usage information

Environment Variables:
  JOB_COUNT       - Concurrent fuzzing processes
  REPO_SOURCE     - Target repository URL
  REPO_STORAGE    - Results storage repo
  FUZZ_DURATION   - Fuzzing duration in seconds
  TARGET_PKGS     - Packages to fuzz

Example:
  go run main.go`

func PurgeWorkspace(logger *slog.Logger) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := os.RemoveAll("out"); err != nil {
		logger.Error("Workspace cleanup failed", "error", err)
	}
}
