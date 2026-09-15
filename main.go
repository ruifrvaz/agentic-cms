// agentic-cms — installer for a markdown-based agentic content management
// system. `agentic-cms init` lays the scaffolding (raw/, docs/, wiki/, skills,
// subagents, templates, schema) on top of the current project folder.
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/ruifrvaz/agentic-cms/scaffold"
)

// version is overridden at build time via -X main.version=vX.Y.Z (see the
// Makefile and the release workflow). It is left at "dev" for `go build .`
// and `go run .` — including `go install github.com/ruifrvaz/agentic-cms@latest`,
// which does not run our ldflags; resolvedVersion() falls back to the Go
// module version embedded by the toolchain in that case.
var version = "dev"

// resolvedVersion returns the ldflags-injected version if set, otherwise
// falls back to the module version recorded by `go install pkg@version`.
func resolvedVersion() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}

const usage = `agentic-cms %s — agentic content management scaffolding

Usage:
  agentic-cms init [dir] [--type <name>]
                            install the scaffolding into dir (default: .).
                            --type layers a content type on top of the base
                            scaffold; omitting it on a project that already
                            has one installed re-applies that type.
  agentic-cms init --type list
                            list the content types embedded in this binary
  agentic-cms update       update the binary to the latest release, then
                            re-run init in the current directory if it looks
                            like an installed project
  agentic-cms version      print version
  agentic-cms help         show this help

init never overwrites your content (raw/, docs/, wiki/, CONTENT.md);
framework files (skills, agents, templates, scripts, hooks) are refreshed
to the installed release on every run. An existing CLAUDE.md is extended
with a managed agentic-cms block instead of replaced.
`

func main() {
	v := resolvedVersion()
	if len(os.Args) < 2 {
		fmt.Printf(usage, v)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "init":
		dir, typeName, err := parseInitArgs(os.Args[2:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "agentic-cms: %v\n\n", err)
			fmt.Printf(usage, v)
			os.Exit(2)
		}
		if typeName == "list" {
			printAvailableTypes()
			return
		}
		if err := runInit(dir, v, typeName); err != nil {
			fmt.Fprintf(os.Stderr, "agentic-cms: %v\n", err)
			os.Exit(1)
		}
	case "update":
		runUpdate(v)
	case "version", "--version", "-v":
		fmt.Println("agentic-cms " + v)
	case "help", "--help", "-h":
		fmt.Printf(usage, v)
	default:
		fmt.Fprintf(os.Stderr, "agentic-cms: unknown command %q\n\n", os.Args[1])
		fmt.Printf(usage, v)
		os.Exit(2)
	}
}

// parseInitArgs splits init's arguments into a target directory (default ".")
// and an optional --type value, in either order ("init --type x dir" and
// "init dir --type x" both work). "--type list" is recognized by the caller
// as a request to print available types rather than install anything.
func parseInitArgs(args []string) (dir, typeName string, err error) {
	dir = "."
	var positional []string
	for i := 0; i < len(args); i++ {
		if args[i] == "--type" {
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("--type requires a value")
			}
			typeName = args[i+1]
			i++
			continue
		}
		positional = append(positional, args[i])
	}
	if len(positional) > 0 {
		dir = positional[0]
	}
	return dir, typeName, nil
}

func printAvailableTypes() {
	types, err := scaffold.AvailableTypes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "agentic-cms: listing content types: %v\n", err)
		os.Exit(1)
	}
	if len(types) == 0 {
		fmt.Println("No content types embedded in this binary.")
		return
	}
	fmt.Println("Available content types:")
	for _, t := range types {
		fmt.Printf("  %s\n", t)
	}
}

// runInit installs the scaffolding into dir. typeName selects a content type
// to layer on top of the base scaffold; when empty, an already-installed
// type recorded in dir's .agentic-cms/TYPE.md is honored automatically, so a
// bare re-init or `update` never silently drops back to a typeless install.
func runInit(dir, version, typeName string) error {
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

	if typeName == "" {
		typeName = scaffold.InstalledType(dir)
	}
	if typeName != "" && !scaffold.HasType(typeName) {
		available, _ := scaffold.AvailableTypes()
		return fmt.Errorf("unknown content type %q (available: %v)", typeName, available)
	}

	abs, _ := os.Getwd()
	if dir != "." {
		abs = dir
	}
	if typeName != "" {
		fmt.Printf("Initializing agentic-cms scaffolding (type: %s) in %s\n", typeName, abs)
	} else {
		fmt.Printf("Initializing agentic-cms scaffolding in %s\n", abs)
	}

	// Read the previously installed scaffold generation before Install
	// restamps .agentic-cms/VERSION with this binary's version.
	prevVersion := scaffold.InstalledVersion(dir)

	res, err := scaffold.Install(dir, version)
	if err != nil {
		return err
	}
	if err := scaffold.InstallGitHook(dir, res); err != nil {
		return err
	}

	if typeName != "" {
		contentJustCreated := false
		for _, f := range res.Created {
			if f == "CONTENT.md" {
				contentJustCreated = true
				break
			}
		}
		if err := scaffold.InstallType(dir, typeName, version, res); err != nil {
			return err
		}
		if contentJustCreated {
			if err := scaffold.ComposeTypeContentMD(dir, typeName); err != nil {
				return err
			}
		}
		if err := scaffold.InstallTypeClaudeMD(dir, typeName, res); err != nil {
			return err
		}
	}
	res.Print()

	report, err := scaffold.ReconcileContentMD(dir, prevVersion, version, typeName)
	if err != nil {
		return err
	}
	if report != nil {
		report.Print()
	}

	fmt.Println()
	fmt.Printf("Done: %d created, %d updated, %d merged, %d skipped.\n",
		len(res.Created), len(res.Updated), len(res.Merged), len(res.Skipped))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  - Read CONTENT.md — the schema your agent follows")
	fmt.Println("  - Greenfield: open your agent and run the content-new skill")
	fmt.Println("  - Brownfield: drop sources into raw/ (or point at an existing")
	fmt.Println("    folder) and run the content-import skill")
	return nil
}
