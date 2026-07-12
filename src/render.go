package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type treeChars struct {
	VLine            string // vertical line
	HLine            string // horizontal line
	RootTop          string // top of root dir
	BranchMid        string // branch with trailing siblings
	BranchEnd        string // last branch, no trailing siblings
	Indent           string // empty indent when parent was last
	IndentPipe       string // indent with pipe when parent has siblings
	DirIndicator     string // shown before directory name
	FileIndicator    string // shown before file name
	SymlinkIndicator string // shown before symlink name
	SymlinkPointer   string // shown between symlink name and target
	SymlinkBroken    string // shown when symlink target is missing
	SocketIndicator  string // shown after socket name
	PipeIndicator    string // shown after named pipe name
}

func newTreeCharsUnicode() treeChars {
	return treeChars{
		VLine:            "│",
		HLine:            "─",
		RootTop:          "╭─ ",
		BranchMid:        "├─ ",
		BranchEnd:        "╰─ ",
		Indent:           "   ",
		IndentPipe:       "│  ",
		DirIndicator:     "▸ ",
		FileIndicator:    "• ",
		SymlinkIndicator: "» ",
		SocketIndicator:  "◦ ",
		PipeIndicator:    "⌇ ",
		SymlinkPointer:   " ⇾ ",
		SymlinkBroken:    " ⇢ ",
	}
}

func newTreeCharsASCII() treeChars {
	return treeChars{
		VLine:            "|",
		HLine:            "-",
		RootTop:          ",- ",
		BranchMid:        "|- ",
		BranchEnd:        "`- ",
		Indent:           "   ",
		IndentPipe:       "|  ",
		DirIndicator:     "# ",
		FileIndicator:    "+ ",
		SymlinkIndicator: "@ ",
		SocketIndicator:  "o ",
		PipeIndicator:    "| ",
		SymlinkPointer:   " -> ",
		SymlinkBroken:    " ~> ",
	}
}

func shouldUseASCII(cfg config, out io.Writer) bool {
	if cfg.format == "ascii" {
		return true
	}
	if os.Getenv("TERM") == "dumb" {
		return true
	}
	f, ok := out.(*os.File)
	if !ok {
		return true
	}
	fi, err := f.Stat()
	if err != nil {
		return true
	}
	return fi.Mode()&os.ModeCharDevice == 0
}

var (
	styleHeader          = color.New(color.FgBlue)
	styleDir             = color.New(color.FgBlue)
	styleDirEmpty        = color.New(color.FgBlue, color.Faint)
	styleFile            = color.New(color.FgWhite)
	styleSymlink         = color.New(color.FgCyan)
	styleSymBroken       = color.New(color.FgRed)
	styleBranch          = color.New(color.FgWhite, color.Faint)
	styleSize            = color.New(color.FgWhite, color.Faint)
	styleFooterHighlight = color.New(color.FgBlue)
)

type walkState struct {
	cfg       config
	chars     treeChars
	out       *bufio.Writer
	matcher   *ignoreMatcher
	dirCount  int
	fileCount int
}

type jsonEntry struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Type     string      `json:"type"`
	Size     int64       `json:"size,omitempty"`
	Children []jsonEntry `json:"children,omitempty"`
}

func renderCharTree(path string, cfg config, out io.Writer) error {
	chars := newTreeCharsUnicode()
	if shouldUseASCII(cfg, out) {
		chars = newTreeCharsASCII()
	}

	matcher := newIgnoreMatcher(cfg)

	bw := bufio.NewWriter(out)
	defer bw.Flush()

	state := &walkState{
		cfg:     cfg,
		chars:   chars,
		out:     bw,
		matcher: matcher,
	}

	fmt.Fprintf(bw, "%s%s\n",
		styleBranch.Sprint(chars.RootTop),
		styleHeader.Sprint(path),
	)
	fmt.Fprintln(bw, styleBranch.Sprint(chars.VLine))

	if err := walk(path, 0, []bool{}, state); err != nil {
		return err
	}

	if !cfg.noSummary {
		fmt.Fprintf(bw, "\nfound %s and %s\n",
			styleFooterHighlight.Sprintf("%d directories", state.dirCount),
			styleFooterHighlight.Sprintf("%d files", state.fileCount),
		)
	}

	return nil
}

func walk(path string, depth int, isLastStack []bool, state *walkState) error {
	if state.cfg.depth > 0 && depth >= state.cfg.depth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	filtered := make([]os.DirEntry, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if name == ".git" {
			continue
		}
		if !state.cfg.all && strings.HasPrefix(name, ".") {
			continue
		}
		if state.cfg.dirsOnly && !entry.IsDir() {
			continue
		}
		if state.matcher.match(path, entry) {
			continue
		}
		filtered = append(filtered, entry)
	}

	for i, entry := range filtered {
		isLast := i == len(filtered)-1
		name := entry.Name()
		entryPath := filepath.Join(path, name)

		prefix := buildPrefix(isLastStack, state.chars)
		connector := state.chars.BranchMid
		if isLast {
			connector = state.chars.BranchEnd
		}

		linfo, err := os.Lstat(entryPath)
		if err != nil {
			continue
		}

		isSymlink := linfo.Mode()&os.ModeSymlink != 0
		isSocket := linfo.Mode()&os.ModeSocket != 0
		isNamedPipe := linfo.Mode()&os.ModeNamedPipe != 0
		isBrokenSymlink := false
		symlinkTarget := ""

		if isSymlink {
			if target, err := os.Readlink(entryPath); err == nil {
				symlinkTarget = target
				if _, err := os.Stat(entryPath); err != nil {
					isBrokenSymlink = true
				}
			}
		}

		var formattedName string
		switch {
		case isBrokenSymlink:
			formattedName = styleSymBroken.Sprint(state.chars.SymlinkIndicator + name)
		case isSymlink:
			formattedName = styleSymlink.Sprint(state.chars.SymlinkIndicator + name)
		case entry.IsDir():
			if dirHasVisibleChildren(entryPath, state) {
				formattedName = styleDir.Sprint(state.chars.DirIndicator + name)
			} else {
				formattedName = styleDirEmpty.Sprint(state.chars.DirIndicator + name)
			}
		default:
			formattedName = styleFile.Sprint(state.chars.FileIndicator + name)
		}

		line := styleBranch.Sprint(prefix+connector) + formattedName

		if isSocket {
			line += styleBranch.Sprint(state.chars.SocketIndicator)
		}
		if isNamedPipe {
			line += styleBranch.Sprint(state.chars.PipeIndicator)
		}
		if isSymlink && symlinkTarget != "" {
			if isBrokenSymlink {
				line += styleSymBroken.Sprint(state.chars.SymlinkBroken + symlinkTarget)
			} else {
				line += styleSymlink.Sprint(state.chars.SymlinkPointer + symlinkTarget)
			}
		}
		if state.cfg.size && !entry.IsDir() {
			line += styleSize.Sprintf(" (%s)", humanSize(linfo.Size()))
		}

		fmt.Fprintln(state.out, line)

		if entry.IsDir() {
			state.dirCount++
		} else {
			state.fileCount++
		}

		if entry.IsDir() && !isSymlink {
			if err := walk(entryPath, depth+1, append(isLastStack, isLast), state); err != nil {
				return err
			}
		}
	}

	return nil
}

func buildPrefix(isLastStack []bool, chars treeChars) string {
	var sb strings.Builder
	for _, isLast := range isLastStack {
		if isLast {
			sb.WriteString(chars.Indent)
		} else {
			sb.WriteString(chars.IndentPipe)
		}
	}
	return sb.String()
}

func dirHasVisibleChildren(path string, state *walkState) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.Name() == ".git" {
			continue
		}
		if !state.cfg.all && strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if state.cfg.dirsOnly && !e.IsDir() {
			continue
		}
		if state.matcher.match(path, e) {
			continue
		}
		return true
	}
	return false
}

func humanSize(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func renderJSONTree(path string, cfg config, out io.Writer) error {
	root, err := buildJSONTree(path, 0, cfg, newIgnoreMatcher(cfg))
	if err != nil {
		return err
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(root)
}

func buildJSONTree(path string, depth int, cfg config, matcher *ignoreMatcher) (jsonEntry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return jsonEntry{}, err
	}

	entry := jsonEntry{
		Name: filepath.Base(path),
		Path: path,
		Type: "file",
		Size: info.Size(),
	}

	if info.IsDir() {
		entry.Type = "dir"
		entry.Size = 0

		if cfg.depth > 0 && depth >= cfg.depth {
			return entry, nil
		}

		dirEntries, err := os.ReadDir(path)
		if err != nil {
			return entry, nil
		}

		for _, de := range dirEntries {
			name := de.Name()
			if name == ".git" {
				continue
			}
			if !cfg.all && strings.HasPrefix(name, ".") {
				continue
			}
			if cfg.dirsOnly && !de.IsDir() {
				continue
			}
			if matcher.match(path, de) {
				continue
			}
			child, err := buildJSONTree(filepath.Join(path, name), depth+1, cfg, matcher)
			if err != nil {
				continue
			}
			entry.Children = append(entry.Children, child)
		}
	}

	return entry, nil
}
