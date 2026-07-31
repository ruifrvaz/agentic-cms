package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var wantFiles = []string{
	"CONTENT.md",
	"CLAUDE.md",
	"raw/README.md",
	"wiki/index.md",
	"wiki/log.md",
	".content-gen/VERSION",
	".content-gen/templates/doc.md",
	".content-gen/templates/entity.md",
	".content-gen/templates/concept.md",
	".content-gen/templates/source.md",
	".content-gen/templates/topic.md",
	".claude/skills/content-new/SKILL.md",
	".claude/skills/content-new-item/SKILL.md",
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
	res, err := Install(dir)
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
	// Templates must keep their placeholders.
	b, _ = os.ReadFile(filepath.Join(dir, ".content-gen", "templates", "doc.md"))
	if !strings.Contains(string(b), "{{TITLE}}") {
		t.Error("template lost {{TITLE}} placeholder")
	}
}

func TestInstallIdempotent(t *testing.T) {
	dir := t.TempDir()
	if _, err := Install(dir); err != nil {
		t.Fatal(err)
	}
	// Mutate a file to prove re-init never overwrites.
	marker := filepath.Join(dir, "wiki", "index.md")
	if err := os.WriteFile(marker, []byte("user content"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := Install(dir)
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

func TestInstallBrownfieldClaudeMD(t *testing.T) {
	dir := t.TempDir()
	existing := "# My project\n\nExisting instructions.\n"
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := Install(dir)
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
	if _, err := Install(dir); err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if strings.Count(string(b), markerBegin) != 1 {
		t.Error("managed block duplicated on re-init")
	}
}
