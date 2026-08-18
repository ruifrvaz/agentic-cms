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

// Result summarizes what Install did.
type Result struct {
	Created []string
	Skipped []string
	Merged  []string
}

// Install writes the scaffold into dir. It is non-destructive: existing files
// are never overwritten (they are reported as skipped). CLAUDE.md gets special
// treatment: if it already exists and does not contain the agentic-cms block,
// the block is appended; otherwise it is left alone.
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

		if _, statErr := os.Stat(target); statErr == nil {
			res.Skipped = append(res.Skipped, rel+" (exists)")
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		// The ac toolkit must be executable.
		mode := os.FileMode(0o644)
		if strings.HasPrefix(rel, filepath.FromSlash(".agentic-cms/scripts/")) && !strings.HasSuffix(rel, ".md") {
			mode = 0o755
		}
		if err := os.WriteFile(target, []byte(content), mode); err != nil {
			return err
		}
		res.Created = append(res.Created, rel)
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
	for _, f := range r.Merged {
		fmt.Printf("  merged   %s\n", f)
	}
	for _, f := range r.Skipped {
		fmt.Printf("  skipped  %s\n", f)
	}
}
