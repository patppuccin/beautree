package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type ignorePattern struct {
	pattern  string // cleaned pattern
	anchored bool   // true if pattern contains /
	dirOnly  bool   // true if pattern has trailing /
	negate   bool   // true if pattern starts with !
	root     string // directory this pattern is relative to
}

type ignoreMatcher struct {
	cfg      config
	patterns []ignorePattern            // cli patterns, pre-compiled once
	cache    map[string][]ignorePattern // gitignore patterns keyed by dir
}

func newIgnoreMatcher(cfg config) *ignoreMatcher {
	m := &ignoreMatcher{
		cfg:   cfg,
		cache: make(map[string][]ignorePattern),
	}
	for _, glob := range cfg.ignore {
		m.patterns = append(m.patterns, ignorePattern{
			pattern:  glob,
			anchored: false,
			dirOnly:  false,
			negate:   false,
			root:     "",
		})
	}
	return m
}

func (m *ignoreMatcher) gitignorePatterns(dir string) []ignorePattern {
	if patterns, ok := m.cache[dir]; ok {
		return patterns
	}
	patterns := loadGitignore(dir)
	m.cache[dir] = patterns
	return patterns
}

func loadGitignore(dir string) []ignorePattern {
	f, err := os.Open(filepath.Join(dir, ".gitignore"))
	if err != nil {
		return nil
	}
	defer f.Close()

	var patterns []ignorePattern
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), " \t\r")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		p := ignorePattern{root: dir}

		if strings.HasPrefix(line, "!") {
			p.negate = true
			line = line[1:]
		}
		if strings.HasSuffix(line, "/") {
			p.dirOnly = true
			line = strings.TrimSuffix(line, "/")
		}
		if strings.Contains(line, "/") {
			p.anchored = true
			line = strings.TrimPrefix(line, "/")
		}

		p.pattern = line
		patterns = append(patterns, p)
	}

	return patterns
}

func (m *ignoreMatcher) match(dir string, entry os.DirEntry) bool {
	if m.cfg.noIgnore {
		return false
	}

	name := entry.Name()
	entryPath := filepath.Join(dir, name)
	isDir := entry.IsDir()

	gitPatterns := m.gitignorePatterns(dir)
	allPatterns := append(gitPatterns, m.patterns...)

	ignored := false
	for _, p := range allPatterns {
		if p.dirOnly && !isDir {
			continue
		}

		matched := false
		if p.anchored {
			rel, err := filepath.Rel(p.root, entryPath)
			if err != nil {
				continue
			}
			matched, _ = filepath.Match(p.pattern, filepath.ToSlash(rel))
		} else {
			matched, _ = filepath.Match(p.pattern, name)
			if !matched && p.root != "" {
				rel, err := filepath.Rel(p.root, entryPath)
				if err == nil {
					matched, _ = filepath.Match(p.pattern, filepath.ToSlash(rel))
				}
			}
		}

		if matched {
			if p.negate {
				ignored = false
			} else {
				ignored = true
			}
		}
	}

	return ignored
}
