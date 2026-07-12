package main

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/spf13/pflag"
)

const envPrefix = "BEAUTREE_"

type config struct {
	Depth     int      `toml:"depth"`      // Max depth to recurse, 0 = unlimited
	DirsOnly  bool     `toml:"dirs_only"`  // Show directories only, skip files
	All       bool     `toml:"all"`        // Include hidden entries (dotfiles)
	Size      bool     `toml:"size"`       // Show human-readable size beside each entry
	Ignore    []string `toml:"ignore"`     // Glob patterns to exclude, stacks with .gitignore
	NoIgnore  bool     `toml:"no_ignore"`  // Disable .gitignore parsing
	NoSummary bool     `toml:"no_summary"` // Disable the summary footer (N dirs, N files)
	Format    string   `toml:"format"`     // Output format: unicode | ascii | json
}

func configPath() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "beautree", "config.toml")
}

func loadFileConfig() config {
	var fc config
	path := configPath()
	if path == "" {
		return fc
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fc
	}
	toml.DecodeFile(path, &fc)
	return fc
}

func resolveConfig() config {
	// layer 1: file config as base
	cfg := loadFileConfig()

	// layer 2: env vars override file
	if v := os.Getenv(envPrefix + "DEPTH"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			cfg.Depth = i
		}
	}
	if v := os.Getenv(envPrefix + "DIRS_ONLY"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.DirsOnly = b
		}
	}
	if v := os.Getenv(envPrefix + "ALL"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.All = b
		}
	}
	if v := os.Getenv(envPrefix + "SIZE"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.Size = b
		}
	}
	if v := os.Getenv(envPrefix + "NO_IGNORE"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.NoIgnore = b
		}
	}
	if v := os.Getenv(envPrefix + "NO_SUMMARY"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.NoSummary = b
		}
	}
	if v := os.Getenv(envPrefix + "FORMAT"); v != "" {
		cfg.Format = v
	}

	// layer 3: CLI flags override env, only if explicitly passed
	pflag.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "depth":
			cfg.Depth = flagTreeDepth
		case "dirs-only":
			cfg.DirsOnly = flagDirsOnly
		case "all":
			cfg.All = flagIncludeAll
		case "size":
			cfg.Size = flagShowSize
		case "ignore":
			cfg.Ignore = append(cfg.Ignore, flagIgnore...)
		case "no-ignore":
			cfg.NoIgnore = flagNoIgnore
		case "no-summary":
			cfg.NoSummary = flagNoSummary
		case "format":
			cfg.Format = flagFormat
		}
	})

	return cfg
}
