package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAvailableTypesIncludesCandidateInterview(t *testing.T) {
	types, err := AvailableTypes()
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, ty := range types {
		if ty == "candidate-interview" {
			found = true
		}
	}
	if !found {
		t.Errorf("AvailableTypes() = %v, want it to include candidate-interview", types)
	}
}

func TestHasType(t *testing.T) {
	if !HasType("candidate-interview") {
		t.Error("HasType(candidate-interview) = false, want true")
	}
	if HasType("nonexistent-type") {
		t.Error("HasType(nonexistent-type) = true, want false")
	}
}

func TestInstalledTypeUntypedProject(t *testing.T) {
	dir := t.TempDir()
	if _, err := Install(dir, testVersion); err != nil {
		t.Fatal(err)
	}
	if got := InstalledType(dir); got != "" {
		t.Errorf("InstalledType() on an untyped project = %q, want empty", got)
	}
}

func TestFrontmatterField(t *testing.T) {
	content := "---\nname: candidate-interview\nversion: 0.1.0\nbase: \"v0.7.0\"\n---\n\nbody\n"
	if got := frontmatterField(content, "name"); got != "candidate-interview" {
		t.Errorf("frontmatterField(name) = %q, want candidate-interview", got)
	}
	if got := frontmatterField(content, "base"); got != "v0.7.0" {
		t.Errorf("frontmatterField(base) = %q, want v0.7.0 (quotes stripped)", got)
	}
	if got := frontmatterField(content, "missing"); got != "" {
		t.Errorf("frontmatterField(missing) = %q, want empty", got)
	}
	if got := frontmatterField("no frontmatter here", "name"); got != "" {
		t.Errorf("frontmatterField on a file with no frontmatter = %q, want empty", got)
	}
}

func TestUntypedInstallUnaffectedByTypesPackage(t *testing.T) {
	// A plain Install (no type involved at all) must produce exactly the base
	// tree — the types.go/types_install.go additions must be fully inert
	// unless InstallType is explicitly invoked.
	dir := t.TempDir()
	if _, err := Install(dir, testVersion); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "exercises")); !os.IsNotExist(err) {
		t.Error("plain Install created exercises/ — type overlay leaked into an untyped install")
	}
	if _, err := os.Stat(filepath.Join(dir, ".agentic-cms", "TYPE.md")); !os.IsNotExist(err) {
		t.Error("plain Install created .agentic-cms/TYPE.md — type overlay leaked into an untyped install")
	}
}

func TestInstallTypeOverlay(t *testing.T) {
	dir := t.TempDir()
	res, err := Install(dir, testVersion)
	if err != nil {
		t.Fatal(err)
	}
	if err := InstallType(dir, "candidate-interview", testVersion, res); err != nil {
		t.Fatal(err)
	}

	// A type template lands alongside the base ones.
	if _, err := os.Stat(filepath.Join(dir, ".agentic-cms/templates/round-plan.md")); err != nil {
		t.Errorf("type template not installed: %v", err)
	}
	// A type skill lands under .claude/skills/.
	if _, err := os.Stat(filepath.Join(dir, ".claude/skills/interview-setup/SKILL.md")); err != nil {
		t.Errorf("type skill not installed: %v", err)
	}
	// The exercises payload, including its go.mod (embedded as
	// go.mod.embedded to dodge the go:embed nested-module exclusion), lands
	// under its real name.
	modPath := filepath.Join(dir, "exercises/001-rate-limiter/go.mod")
	b, err := os.ReadFile(modPath)
	if err != nil {
		t.Fatalf("exercises go.mod not installed: %v", err)
	}
	if !strings.HasPrefix(string(b), "module ") {
		t.Errorf("exercises go.mod content looks wrong: %q", string(b))
	}
	if _, err := os.Stat(filepath.Join(dir, "exercises/001-rate-limiter/go.mod.embedded")); !os.IsNotExist(err) {
		t.Error("go.mod.embedded leaked onto disk under its embedded name")
	}

	// TYPE.md's {{VERSION}} placeholder is stamped.
	typeMD, err := os.ReadFile(filepath.Join(dir, ".agentic-cms/TYPE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(typeMD), "{{VERSION}}") {
		t.Error("TYPE.md still contains the {{VERSION}} placeholder")
	}
	if !strings.Contains(string(typeMD), testVersion) {
		t.Errorf("TYPE.md does not contain the stamped version %q", testVersion)
	}
	if got := InstalledType(dir); got != "candidate-interview" {
		t.Errorf("InstalledType() after InstallType = %q, want candidate-interview", got)
	}

	// Type templates keep {{DATE}}/{{TITLE}} placeholders live (page-creation
	// time, not install time) — same contract as the base templates.
	roundPlan, err := os.ReadFile(filepath.Join(dir, ".agentic-cms/templates/round-plan.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(roundPlan), "{{TITLE}}") {
		t.Error("type template lost its {{TITLE}} placeholder on install")
	}

	// Re-running InstallType always overwrites (type files are framework-owned).
	if err := os.WriteFile(modPath, []byte("locally edited\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	res2, err := Install(dir, testVersion)
	if err != nil {
		t.Fatal(err)
	}
	if err := InstallType(dir, "candidate-interview", testVersion, res2); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(modPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(after), "locally edited") {
		t.Error("InstallType did not refresh a locally-edited type file — type paths must be framework-owned")
	}
}

func TestInstallTypeUnknownType(t *testing.T) {
	dir := t.TempDir()
	res := &Result{}
	if err := InstallType(dir, "nonexistent-type", testVersion, res); err == nil {
		t.Error("InstallType with an unknown type name returned nil error, want an error")
	}
}

func TestInstallTypeClaudeMDAppendOnce(t *testing.T) {
	dir := t.TempDir()
	res, err := Install(dir, testVersion)
	if err != nil {
		t.Fatal(err)
	}
	if err := InstallTypeClaudeMD(dir, "candidate-interview", res); err != nil {
		t.Fatal(err)
	}
	first, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(first), "<!-- agentic-cms:type:candidate-interview:begin -->") {
		t.Error("type CLAUDE.md block not appended")
	}
	if !strings.Contains(string(first), "<!-- agentic-cms:begin -->") {
		t.Error("base CLAUDE.md block missing after type block append")
	}

	// Re-running must not duplicate the block.
	if err := InstallTypeClaudeMD(dir, "candidate-interview", res); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(second), "<!-- agentic-cms:type:candidate-interview:begin -->") != 1 {
		t.Error("type CLAUDE.md block duplicated on a second run")
	}
}

func TestComposeTypeContentMD(t *testing.T) {
	dir := t.TempDir()
	if _, err := Install(dir, testVersion); err != nil {
		t.Fatal(err)
	}
	if err := ComposeTypeContentMD(dir, "candidate-interview"); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(dir, "CONTENT.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(content)
	for _, want := range []string{
		"## Type: candidate-interview",
		"type: doc | entity | concept | source | note | rehearsal",
		"exercises/            type-specific: runnable coding drills",
		"TYPE.md             the installed content type and what it owns",
		"**Type operations**",
		"run `interview-setup` first",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("composed CONTENT.md missing expected text: %q", want)
		}
	}
	// The base schema must still be intact — composition is additive.
	if !strings.Contains(text, "## Classification") {
		t.Error("composed CONTENT.md lost the base Classification section")
	}
}

func TestApplyTypeFragmentMissingAnchorErrors(t *testing.T) {
	base := "line one\nline two\n"
	frag := `<!-- ac-type-anchor: after "line does not exist" -->
inserted
`
	if _, err := applyTypeFragment(base, frag); err == nil {
		t.Error("applyTypeFragment with a missing anchor returned nil error, want an error")
	}
}

func TestApplyTypeFragmentModes(t *testing.T) {
	base := "alpha\nbeta\ngamma\n"
	frag := `<!-- ac-type-anchor: after "alpha" -->
after-alpha
<!-- ac-type-anchor: before "gamma" -->
before-gamma
<!-- ac-type-anchor: replace "beta" -->
replaced-beta
`
	got, err := applyTypeFragment(base, frag)
	if err != nil {
		t.Fatal(err)
	}
	want := "alpha\nafter-alpha\nreplaced-beta\nbefore-gamma\ngamma\n"
	if got != want {
		t.Errorf("applyTypeFragment() =\n%q\nwant\n%q", got, want)
	}
}

func TestReconcileContentMDTypeSectionMissing(t *testing.T) {
	dir := t.TempDir()
	if _, err := Install(dir, testVersion); err != nil {
		t.Fatal(err)
	}
	if err := ComposeTypeContentMD(dir, "candidate-interview"); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "CONTENT.md")
	b, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	// Simulate a stale typed install predating the Type section.
	start := strings.Index(string(b), "## Type: candidate-interview")
	end := strings.Index(string(b), "## Operations")
	if start == -1 || end == -1 {
		t.Fatal("test setup: composed CONTENT.md missing expected anchors")
	}
	stripped := string(b)[:start] + string(b)[end:]
	if err := os.WriteFile(target, []byte(stripped), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := ReconcileContentMD(dir, "v0.7.0", testVersion, "candidate-interview")
	if err != nil {
		t.Fatal(err)
	}
	if report == nil {
		t.Fatal("expected a reconciliation report for a missing Type section, got nil")
	}
	found := false
	for _, s := range report.MissingSections {
		if s == "Type: candidate-interview" {
			found = true
		}
	}
	if !found {
		t.Errorf("missing sections %v do not include Type: candidate-interview", report.MissingSections)
	}
	sidecar, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(UpstreamSidecar)))
	if err != nil {
		t.Fatalf("sidecar not written: %v", err)
	}
	if !strings.Contains(string(sidecar), "## Type: candidate-interview") {
		t.Error("sidecar does not carry the composed type section")
	}
	after, _ := os.ReadFile(target)
	if string(after) != stripped {
		t.Error("reconciliation modified the user's CONTENT.md")
	}
}

func TestReconcileContentMDTypedCleanInstall(t *testing.T) {
	dir := t.TempDir()
	res, err := Install(dir, testVersion)
	if err != nil {
		t.Fatal(err)
	}
	if err := InstallType(dir, "candidate-interview", testVersion, res); err != nil {
		t.Fatal(err)
	}
	if err := ComposeTypeContentMD(dir, "candidate-interview"); err != nil {
		t.Fatal(err)
	}
	report, err := ReconcileContentMD(dir, testVersion, testVersion, "candidate-interview")
	if err != nil {
		t.Fatal(err)
	}
	if report != nil {
		t.Errorf("freshly composed typed install reported missing sections: %v", report.MissingSections)
	}
}
