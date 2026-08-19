package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// testVersion is the stand-in binary version tests install with.
const testVersion = "v9.9.9-test"

var wantFiles = []string{
	"CONTENT.md",
	"CLAUDE.md",
	"raw/README.md",
	"wiki/index.md",
	"wiki/log.md",
	".agentic-cms/VERSION",
	".agentic-cms/scripts/README.md",
	".agentic-cms/scripts/ac-page",
	".agentic-cms/scripts/ac-index",
	".agentic-cms/scripts/ac-log",
	".agentic-cms/scripts/ac-links",
	".agentic-cms/scripts/ac-inventory",
	".agentic-cms/scripts/ac-search",
	".agentic-cms/scripts/ac-classify",
	".agentic-cms/hooks/pre-commit",
	".agentic-cms/templates/doc.md",
	".agentic-cms/templates/entity.md",
	".agentic-cms/templates/concept.md",
	".agentic-cms/templates/source.md",
	".agentic-cms/templates/topic.md",
	".claude/settings.json",
	".codex/hooks.json",
	".claude/skills/content-new/SKILL.md",
	".claude/skills/content-manage-item/SKILL.md",
	".claude/skills/content-query/SKILL.md",
	".claude/skills/content-research/SKILL.md",
	".claude/skills/content-import/SKILL.md",
	".claude/skills/content-add-notes/SKILL.md",
	".claude/skills/content-list/SKILL.md",
	".claude/skills/content-lint/SKILL.md",
	".claude/skills/content-export/SKILL.md",
	".claude/agents/content-researcher.md",
	".claude/agents/content-importer.md",
	".claude/agents/content-exporter.md",
}

func TestInstallGreenfield(t *testing.T) {
	dir := t.TempDir()
	res, err := Install(dir, testVersion)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range wantFiles {
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(f))); err != nil {
			t.Errorf("missing %s: %v", f, err)
		}
	}
	if len(res.Skipped) != 0 {
		t.Errorf("greenfield install skipped files: %v", res.Skipped)
	}
	// {{DATE}} must be substituted.
	b, _ := os.ReadFile(filepath.Join(dir, "wiki", "log.md"))
	if strings.Contains(string(b), "{{DATE}}") {
		t.Error("wiki/log.md still contains {{DATE}} placeholder")
	}
	if !strings.Contains(string(b), time.Now().Format("2006-01-02")) {
		t.Error("wiki/log.md missing today's date")
	}
	// Templates must keep all placeholders, including {{DATE}} — ac-page fills
	// it in at page-creation time, not at scaffold-install time.
	b, _ = os.ReadFile(filepath.Join(dir, ".agentic-cms", "templates", "doc.md"))
	if !strings.Contains(string(b), "{{TITLE}}") {
		t.Error("template lost {{TITLE}} placeholder")
	}
	if !strings.Contains(string(b), "{{DATE}}") {
		t.Error("template lost {{DATE}} placeholder — it must stay live for ac-page new")
	}
	// .agentic-cms/scripts/ is source code, not a page: ac-page's own source uses
	// the literal string "{{DATE}}" as a placeholder key, so a blanket
	// {{DATE}} substitution at install time would corrupt the script itself.
	b, _ = os.ReadFile(filepath.Join(dir, ".agentic-cms", "scripts", "ac-page"))
	if !strings.Contains(string(b), "{{DATE}}") {
		t.Error("ac-page source corrupted — installer substituted {{DATE}} inside the script itself")
	}
	// Toolkit scripts must be executable.
	fi, err := os.Stat(filepath.Join(dir, ".agentic-cms", "scripts", "ac-index"))
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode()&0o111 == 0 {
		t.Error(".agentic-cms/scripts/ac-index is not executable")
	}
	// The versioned git hook script must be executable too — it lives outside
	// .agentic-cms/scripts/, a distinct directory the installer's mode logic must
	// also cover (this regressed once: it shipped 0644 and .git/hooks/pre-commit's
	// `[ -x ... ]` guard silently no-op'd the whole gate).
	fi, err = os.Stat(filepath.Join(dir, ".agentic-cms", "hooks", "pre-commit"))
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode()&0o111 == 0 {
		t.Error(".agentic-cms/hooks/pre-commit is not executable")
	}
}

func TestInstallIdempotent(t *testing.T) {
	dir := t.TempDir()
	if _, err := Install(dir, testVersion); err != nil {
		t.Fatal(err)
	}
	// Mutate a file to prove re-init never overwrites.
	marker := filepath.Join(dir, "wiki", "index.md")
	if err := os.WriteFile(marker, []byte("user content"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := Install(dir, testVersion)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Created) != 0 {
		t.Errorf("second install created files: %v", res.Created)
	}
	b, _ := os.ReadFile(marker)
	if string(b) != "user content" {
		t.Error("re-init overwrote an existing file")
	}
}

func TestInstallOverwritesFrameworkFiles(t *testing.T) {
	dir := t.TempDir()
	if _, err := Install(dir, testVersion); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(dir, ".claude", "skills", "content-lint", "SKILL.md")
	original, err := os.ReadFile(marker)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(marker, []byte("locally edited content"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := Install(dir, testVersion)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.FromSlash(".claude/skills/content-lint/SKILL.md")
	found := false
	for _, f := range res.Updated {
		if f == want {
			found = true
		}
	}
	if !found {
		t.Errorf("expected %s in Updated, got Updated=%v Skipped=%v", want, res.Updated, res.Skipped)
	}
	for _, f := range res.Skipped {
		if f == want+" (exists)" {
			t.Errorf("%s reported as skipped, should have been overwritten", want)
		}
	}
	b, err := os.ReadFile(marker)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != string(original) {
		t.Error("framework-owned file was not reset to the embedded version on re-install")
	}
}

func TestInstallBrownfieldClaudeMD(t *testing.T) {
	dir := t.TempDir()
	existing := "# My project\n\nExisting instructions.\n"
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := Install(dir, testVersion)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Merged) != 1 {
		t.Fatalf("expected 1 merged file, got %v", res.Merged)
	}
	b, _ := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	s := string(b)
	if !strings.HasPrefix(s, "# My project") {
		t.Error("existing CLAUDE.md content lost")
	}
	if !strings.Contains(s, markerBegin) || !strings.Contains(s, markerEnd) {
		t.Error("managed block not appended")
	}
	// Third run: block present, must be skipped, not duplicated.
	if _, err := Install(dir, testVersion); err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if strings.Count(string(b), markerBegin) != 1 {
		t.Error("managed block duplicated on re-init")
	}
}
