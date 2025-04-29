package core

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type EngineConfig struct {
	SourceRepo      string
	StorageRepo     string
	TargetPackages  []string
	ExecutionWindow string
	ParallelJobs    int
}

func LoadEnvironment() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("Environment load error: %v", err)
	}
	return nil
}

func InitEngine() (*EngineConfig, error) {
	cfg := &EngineConfig{
		SourceRepo:     os.Getenv("REPO_SOURCE"),
		StorageRepo:    os.Getenv("REPO_STORAGE"),
		ExecutionWindow: "120s",
	}

	if cfg.SourceRepo == "" || cfg.StorageRepo == "" {
		return nil, errors.New("Required repository paths not configured")
	}

	if window := os.Getenv("FUZZ_DURATION"); window != "" {
		sec, err := strconv.Atoi(window)
		if err != nil {
			return nil, fmt.Errorf("Invalid duration format: %s", window)
		}
		cfg.ExecutionWindow = fmt.Sprintf("%ds", sec)
	}

	cfg.ParallelJobs = calculateConcurrency()
	
	if pkgList := os.Getenv("TARGET_PKGS"); pkgList != "" {
		cfg.TargetPackages = strings.Fields(pkgList)
	} else {
		return nil, errors.New("Target packages not specified")
	}

	return cfg, nil
}

func calculateConcurrency() int {
	if envVal := os.Getenv("JOB_COUNT"); envVal != "" {
		if count, err := strconv.Atoi(envVal); err == nil && count > 0 {
			if count > runtime.NumCPU() {
				return runtime.NumCPU()
			}
			return count
		}
	}
	return runtime.NumCPU()
}
