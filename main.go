// agentic-cms — installer for a markdown-based agentic content management
// system. `agentic-cms init` lays the scaffolding (raw/, docs/, wiki/, skills,
// subagents, templates, schema) on top of the current project folder.
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ruifrvaz/agentic-cms/scaffold"
)

var version = "0.1.0"

const usage = `agentic-cms %s — agentic content management scaffolding

Usage:
  agentic-cms init [dir]   install the scaffolding into dir (default: .)
  agentic-cms version      print version
  agentic-cms help         show this help

init is non-destructive: existing files are never overwritten. An existing
CLAUDE.md is extended with a managed agentic-cms block instead of replaced.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(usage, version)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "init":
		dir := "."
		if len(os.Args) > 2 {
			dir = os.Args[2]
		}
		if err := runInit(dir); err != nil {
			fmt.Fprintf(os.Stderr, "agentic-cms: %v\n", err)
			os.Exit(1)
		}
	case "version", "--version", "-v":
		fmt.Println("agentic-cms " + version)
	case "help", "--help", "-h":
		fmt.Printf(usage, version)
	default:
		fmt.Fprintf(os.Stderr, "agentic-cms: unknown command %q\n\n", os.Args[1])
		fmt.Printf(usage, version)
		os.Exit(2)
	}
}

func runInit(dir string) error {
	if runtime.GOOS != "linux" {
		fmt.Fprintln(os.Stderr, "warning: agentic-cms currently targets Linux; proceeding anyway")
	}

	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("target directory %q: %w", dir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("target %q is not a directory", dir)
	}

	abs, _ := os.Getwd()
	if dir != "." {
		abs = dir
	}
	fmt.Printf("Initializing agentic-cms scaffolding in %s\n", abs)

	res, err := scaffold.Install(dir)
	if err != nil {
		return err
	}
	res.Print()

	fmt.Println()
	fmt.Printf("Done: %d created, %d merged, %d skipped.\n",
		len(res.Created), len(res.Merged), len(res.Skipped))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  - Read CONTENT.md — the schema your agent follows")
	fmt.Println("  - Greenfield: open your agent and run the content-new skill")
	fmt.Println("  - Brownfield: drop sources into raw/ (or point at an existing")
	fmt.Println("    folder) and run the content-import skill")
	return nil
}
