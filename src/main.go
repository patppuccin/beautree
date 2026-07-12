package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

const appBanner = `
██▄ ██▀ ▄▀▄ █ █ ▀█▀ █▀▄ ██▀ ██▀
█▄█ █▄▄ █▀█ ▀▄█  █  █▀▄ █▄▄ █▄▄
`

const (
	appName = "beautree"
	appDesc = "Pretty directory tree viewer for your terminal"
)

var (
	appVersion  = "dev"
	buildCommit = "none"
	buildDate   = "unknown"
)

var (
	flagHelp    bool
	flagVersion bool
	flagOutput  string
	flagNoColor bool

	flagTreeDepth  int
	flagIncludeAll bool
	flagDirsOnly   bool
	flagShowSize   bool
	flagIgnore     []string
	flagNoIgnore   bool
	flagNoSummary  bool
	flagFormat     string
)

func bannerString(msg string) string {
	return color.New(color.FgWhite, color.Faint).Sprint(appBanner) + "\n" +
		color.New(color.FgBlue).Sprint(msg) + "\n\n"
}

func versionString() string {
	date := buildDate
	if t, err := time.Parse(time.RFC3339, buildDate); err == nil {
		date = t.UTC().Format("02 Jan 2006 15:04 UTC")
	}
	return fmt.Sprintf("%s %s (commit: %s, built: %s)", appName, appVersion, buildCommit, date)
}

func main() {
	os.Unsetenv("LS_COLORS")
	os.Unsetenv("LSCOLORS")
	pflag.CommandLine.SortFlags = false

	pflag.IntVarP(&flagTreeDepth, "depth", "L", 0, "Max depth to recurse (0 = unlimited)")
	pflag.BoolVarP(&flagIncludeAll, "all", "a", false, "Include hidden files and dirs")
	pflag.BoolVarP(&flagDirsOnly, "dirs-only", "d", false, "Show directories only")
	pflag.BoolVar(&flagShowSize, "size", false, "Show size beside each entry")

	pflag.StringArrayVarP(&flagIgnore, "ignore", "I", nil, "Exclude entries matching GLOB (repeatable)")
	pflag.BoolVar(&flagNoIgnore, "no-ignore", false, "Disable .gitignore parsing")

	pflag.StringVarP(&flagFormat, "format", "f", "default", "Output format: default|ascii|json")
	pflag.BoolVar(&flagNoColor, "no-color", false, "Disable color output")
	pflag.BoolVar(&flagNoSummary, "no-summary", false, "Disable summary footer")
	pflag.StringVarP(&flagOutput, "output", "o", "", "Write output to file instead of stdout")

	pflag.BoolVarP(&flagVersion, "version", "v", false, "Print version")
	pflag.BoolVarP(&flagHelp, "help", "h", false, "Show this help")

	pflag.Usage = func() {
		fmt.Print(bannerString(appDesc))
		fmt.Print("Usage: " + appName + " [path] [flags]\n\n")
		fmt.Print("Flags:\n")
		pflag.PrintDefaults()
	}

	pflag.Parse()

	if flagNoColor || os.Getenv("NO_COLOR") != "" {
		color.NoColor = true
	}

	if flagHelp {
		pflag.Usage()
		return
	}

	if flagVersion {
		fmt.Println(versionString())
		return
	}

	cfg := resolveConfig()

	rootPath := "."
	if pflag.NArg() > 0 {
		rootPath = pflag.Arg(0)
	}

	out := io.Writer(os.Stdout)
	if flagOutput != "" {
		f, err := os.Create(flagOutput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot open output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	}

	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot access %q: %v\n", absPath, err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "error: %q is not a directory\n", absPath)
		os.Exit(1)
	}

	var renderr error
	if cfg.Format == "json" {
		renderr = renderJSONTree(absPath, cfg, out)
	} else {
		renderr = renderCharTree(absPath, cfg, out)
	}
	if renderr != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", renderr)
		os.Exit(1)
	}
}
