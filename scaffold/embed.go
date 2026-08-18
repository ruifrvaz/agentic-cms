// Package scaffold embeds the agentic-cms scaffolding tree and installs it
// into a target project directory.
package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Tree holds the embedded scaffolding. The all: prefix is required so that
// dot-directories like .claude and .agentic-cms are included.
//
//go:embed all:tree
var Tree embed.FS

const treeRoot = "tree"

// markers delimiting the block agentic-cms manages inside CLAUDE.md.
const (
	markerBegin = "<!-- agentic-cms:begin -->"
	markerEnd   = "<!-- agentic-cms:end -->"
)

// frameworkOwnedPrefixes lists scaffold paths that are versioned framework
// logic — skills, agents, templates, scripts, hooks — rather than user-owned
// project content. Install always overwrites these on every run so a project
// stays in sync with the installed agentic-cms release, even if the file was
// locally edited.
var frameworkOwnedPrefixes = []string{
	filepath.FromSlash(".claude/skills/"),
	filepath.FromSlash(".claude/agents/"),
	filepath.FromSlash(".agentic-cms/templates/"),
	filepath.FromSlash(".agentic-cms/scripts/"),
	filepath.FromSlash(".agentic-cms/hooks/"),
	filepath.FromSlash(".codex/"),
}

func isFrameworkOwned(rel string) bool {
	for _, prefix := range frameworkOwnedPrefixes {
		if strings.HasPrefix(rel, prefix) {
			return true
		}
	}
	return false
}

// Result summarizes what Install did.
type Result struct {
	Created []string
	Updated []string
	Skipped []string
	Merged  []string
}

// Install writes the scaffold into dir. User-owned content files are
// non-destructive: an existing file is never overwritten (it is reported as
// skipped). Framework-owned scaffolding logic — skills, agents, templates,
// scripts, hooks — is always overwritten with the embedded version, even if
// it already exists, so that init/update actually refreshes it. CLAUDE.md
// gets special treatment: if it already exists and does not contain the
// agentic-cms block, the block is appended; otherwise it is left alone.
func Install(dir string) (*Result, error) {
	res := &Result{}
	date := time.Now().Format("2006-01-02")

	err := fs.WalkDir(Tree, treeRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(treeRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dir, filepath.FromSlash(rel))

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		data, err := Tree.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(data)
		// Templates keep {{DATE}} as a live placeholder — ac-page fills it in at
		// page-creation time, not at scaffold-install time. .agentic-cms/scripts/ is
		// source code: it references the literal string "{{DATE}}" as a
		// placeholder key, which a blanket substitution would corrupt.
		if !strings.HasPrefix(rel, filepath.FromSlash(".agentic-cms/templates/")) &&
			!strings.HasPrefix(rel, filepath.FromSlash(".agentic-cms/scripts/")) {
			content = strings.ReplaceAll(content, "{{DATE}}", date)
		}

		if rel == "CLAUDE.md" {
			return installClaudeMD(target, content, res)
		}

		_, statErr := os.Stat(target)
		exists := statErr == nil
		if exists && !isFrameworkOwned(rel) {
			res.Skipped = append(res.Skipped, rel+" (exists)")
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		// The ac toolkit and the versioned git hook scripts must be executable.
		mode := os.FileMode(0o644)
		isScriptDir := strings.HasPrefix(rel, filepath.FromSlash(".agentic-cms/scripts/")) ||
			strings.HasPrefix(rel, filepath.FromSlash(".agentic-cms/hooks/"))
		if isScriptDir && !strings.HasSuffix(rel, ".md") {
			mode = 0o755
		}
		if err := os.WriteFile(target, []byte(content), mode); err != nil {
			return err
		}
		if exists {
			res.Updated = append(res.Updated, rel)
		} else {
			res.Created = append(res.Created, rel)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// installClaudeMD creates CLAUDE.md if absent, or appends the managed block to
// an existing CLAUDE.md that does not contain it yet.
func installClaudeMD(target, block string, res *Result) error {
	existing, err := os.ReadFile(target)
	if os.IsNotExist(err) {
		if werr := os.WriteFile(target, []byte(block), 0o644); werr != nil {
			return werr
		}
		res.Created = append(res.Created, "CLAUDE.md")
		return nil
	}
	if err != nil {
		return err
	}
	if strings.Contains(string(existing), markerBegin) {
		res.Skipped = append(res.Skipped, "CLAUDE.md (block already present)")
		return nil
	}
	out := strings.TrimRight(string(existing), "\n") + "\n\n" + block
	if werr := os.WriteFile(target, []byte(out), 0o644); werr != nil {
		return werr
	}
	res.Merged = append(res.Merged, "CLAUDE.md (block appended)")
	return nil
}

// Print writes a human-readable summary of the install result.
func (r *Result) Print() {
	for _, f := range r.Created {
		fmt.Printf("  created  %s\n", f)
	}
	for _, f := range r.Updated {
		fmt.Printf("  updated  %s\n", f)
	}
	for _, f := range r.Merged {
		fmt.Printf("  merged   %s\n", f)
	}
	for _, f := range r.Skipped {
		fmt.Printf("  skipped  %s\n", f)
	}
}
