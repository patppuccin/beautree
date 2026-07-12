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
	depth     int      `toml:"depth"`      // Max depth to recurse, 0 = unlimited
	dirsOnly  bool     `toml:"dirs_only"`  // Show directories only, skip files
	all       bool     `toml:"all"`        // Include hidden entries (dotfiles)
	size      bool     `toml:"size"`       // Show human-readable size beside each entry
	ignore    []string `toml:"ignore"`     // Glob patterns to exclude, stacks with .gitignore
	noIgnore  bool     `toml:"no_ignore"`  // Disable .gitignore parsing
	noSummary bool     `toml:"no_summary"` // Disable the summary footer (N dirs, N files)
	format    string   `toml:"format"`     // Output format: unicode | ascii | json
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
			cfg.depth = i
		}
	}
	if v := os.Getenv(envPrefix + "DIRS_ONLY"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.dirsOnly = b
		}
	}
	if v := os.Getenv(envPrefix + "ALL"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.all = b
		}
	}
	if v := os.Getenv(envPrefix + "SIZE"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.size = b
		}
	}
	if v := os.Getenv(envPrefix + "NO_IGNORE"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.noIgnore = b
		}
	}
	if v := os.Getenv(envPrefix + "NO_SUMMARY"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.noSummary = b
		}
	}
	if v := os.Getenv(envPrefix + "FORMAT"); v != "" {
		cfg.format = v
	}

	// layer 3: CLI flags override env, only if explicitly passed
	pflag.Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "depth":
			cfg.depth = flagTreeDepth
		case "dirs-only":
			cfg.dirsOnly = flagDirsOnly
		case "all":
			cfg.all = flagIncludeAll
		case "size":
			cfg.size = flagShowSize
		case "ignore":
			cfg.ignore = append(cfg.ignore, flagIgnore...)
		case "no-ignore":
			cfg.noIgnore = flagNoIgnore
		case "no-summary":
			cfg.noSummary = flagNoSummary
		case "format":
			cfg.format = flagFormat
		}
	})

	return cfg
}
