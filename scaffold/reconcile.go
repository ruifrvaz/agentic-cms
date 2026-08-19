package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UpstreamSidecar is where the shipped CONTENT.md is written for manual merge
// when the installed copy is missing upstream schema sections.
const UpstreamSidecar = ".agentic-cms/CONTENT.upstream.md"

// ContentReport describes upstream schema sections absent from a project's
// customized CONTENT.md, discovered by ReconcileContentMD.
type ContentReport struct {
	MissingSections  []string
	InstalledVersion string // scaffold version recorded before this install ("" if unknown)
	ShippedVersion   string // version of the binary performing the install
	Sidecar          string // repo-relative path the upstream copy was written to
}

// ReconcileContentMD compares the installed CONTENT.md against the shipped
// one, keyed on top-level "## " schema section headings. When upstream
// sections are missing it writes the current upstream copy to UpstreamSidecar
// and returns a report naming the absent sections. It NEVER edits the user's
// CONTENT.md — the schema file is user-co-evolved by design, and unlike
// CLAUDE.md it has no managed-block convention to merge into.
//
// Returns (nil, nil) when the installed CONTENT.md is absent (fresh install:
// Install just wrote the current copy) or already carries every upstream
// section heading.
func ReconcileContentMD(dir, installedVersion, shippedVersion string) (*ContentReport, error) {
	installed, err := os.ReadFile(filepath.Join(dir, "CONTENT.md"))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	shipped, err := Tree.ReadFile(treeRoot + "/CONTENT.md")
	if err != nil {
		return nil, err
	}
	shippedText := strings.ReplaceAll(string(shipped), "{{DATE}}", time.Now().Format("2006-01-02"))

	installedHeadings := map[string]bool{}
	for _, h := range sectionHeadings(string(installed)) {
		installedHeadings[h] = true
	}
	var missing []string
	for _, h := range sectionHeadings(shippedText) {
		if !installedHeadings[h] {
			missing = append(missing, h)
		}
	}
	if len(missing) == 0 {
		return nil, nil
	}

	sidecar := filepath.Join(dir, filepath.FromSlash(UpstreamSidecar))
	if err := os.MkdirAll(filepath.Dir(sidecar), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(sidecar, []byte(shippedText), 0o644); err != nil {
		return nil, err
	}
	return &ContentReport{
		MissingSections:  missing,
		InstalledVersion: installedVersion,
		ShippedVersion:   shippedVersion,
		Sidecar:          UpstreamSidecar,
	}, nil
}

// sectionHeadings returns the text of every top-level "## " heading, in order.
func sectionHeadings(text string) []string {
	var out []string
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "## ") {
			out = append(out, strings.TrimSpace(strings.TrimPrefix(line, "## ")))
		}
	}
	return out
}

// Print writes a human-readable reconciliation report.
func (r *ContentReport) Print() {
	fmt.Println()
	fmt.Printf("CONTENT.md reconciliation: your CONTENT.md is missing %d upstream schema section(s):\n", len(r.MissingSections))
	for _, s := range r.MissingSections {
		fmt.Printf("  - %s\n", s)
	}
	fmt.Printf("The current upstream schema was written to %s for manual merge\n", r.Sidecar)
	fmt.Println("(agentic-cms never edits your CONTENT.md).")
	installed := r.InstalledVersion
	if installed == "" {
		installed = "unknown (pre-versioned scaffold)"
	}
	fmt.Printf("Installed scaffold: %s → shipping: %s\n", installed, r.ShippedVersion)
}
