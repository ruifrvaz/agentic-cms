package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// This file mirrors types_test.go's candidate-interview coverage for the
// second type, recruiter-interview — the direct test of whether the task
// 013 mechanism (InstallType/ComposeTypeContentMD/InstallTypeClaudeMD/
// applyTypeFragment) is genuinely type-name-parameterized, or only happened
// to work for one type.

func TestAvailableTypesIncludesBothTypes(t *testing.T) {
	types, err := AvailableTypes()
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"candidate-interview": false, "recruiter-interview": false}
	for _, ty := range types {
		if _, ok := want[ty]; ok {
			want[ty] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("AvailableTypes() = %v, missing %q", types, name)
		}
	}
}

func TestHasTypeRecruiterInterview(t *testing.T) {
	if !HasType("recruiter-interview") {
		t.Error("HasType(recruiter-interview) = false, want true")
	}
}

func TestInstallTypeRecruiterInterviewOverlay(t *testing.T) {
	dir := t.TempDir()
	res, err := Install(dir, testVersion)
	if err != nil {
		t.Fatal(err)
	}
	if err := InstallType(dir, "recruiter-interview", testVersion, res); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dir, ".agentic-cms/templates/scorecard.md")); err != nil {
		t.Errorf("type template not installed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, ".claude/skills/recruiter-setup/SKILL.md")); err != nil {
		t.Errorf("type skill not installed: %v", err)
	}

	// The take-home exercise carries an interviewer-only GRADING.md
	// alongside the same go.mod-embedding workaround as candidate-interview.
	modPath := filepath.Join(dir, "exercises/001-rate-limiter/go.mod")
	b, err := os.ReadFile(modPath)
	if err != nil {
		t.Fatalf("exercises go.mod not installed: %v", err)
	}
	if !strings.HasPrefix(string(b), "module ") {
		t.Errorf("exercises go.mod content looks wrong: %q", string(b))
	}
	if _, err := os.Stat(filepath.Join(dir, "exercises/001-rate-limiter/GRADING.md")); err != nil {
		t.Errorf("GRADING.md not installed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "exercises/001-rate-limiter/go.mod.embedded")); !os.IsNotExist(err) {
		t.Error("go.mod.embedded leaked onto disk under its embedded name")
	}

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
	// Same task-014 regression class, for a second type: the {{DATE}}
	// exclusion in InstallType is keyed by relative path (typeManifestFile),
	// not type name, so it must already cover this type too.
	if !strings.Contains(string(typeMD), "{{DATE}}") {
		t.Error("TYPE.md's placeholder-name documentation lost its literal {{DATE}} mention")
	}
	if got := InstalledType(dir); got != "recruiter-interview" {
		t.Errorf("InstalledType() after InstallType = %q, want recruiter-interview", got)
	}
}

func TestComposeTypeContentMDRecruiterInterview(t *testing.T) {
	dir := t.TempDir()
	if _, err := Install(dir, testVersion); err != nil {
		t.Fatal(err)
	}
	if err := ComposeTypeContentMD(dir, "recruiter-interview"); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(dir, "CONTENT.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(content)
	for _, want := range []string{
		"## Type: recruiter-interview",
		"exercises/            type-specific: take-home assignments as runnable code",
		"TYPE.md             the installed content type and what it owns",
		"**Type operations**",
		"only the organization's own",
		"cNNN-profile.md",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("composed CONTENT.md missing expected text: %q", want)
		}
	}
	// No new page type for this type — unlike candidate-interview, the
	// frontmatter type: enum must stay at its base value.
	if strings.Contains(text, "type: doc | entity | concept | source | note | rehearsal") {
		t.Error("composed CONTENT.md unexpectedly carries candidate-interview's rehearsal page type")
	}
	if !strings.Contains(text, "type: doc | entity | concept | source | note\n") {
		t.Error("composed CONTENT.md's frontmatter type: enum should be unchanged from base (no new page type for this type)")
	}
	if !strings.Contains(text, "## Classification") {
		t.Error("composed CONTENT.md lost the base Classification section")
	}
}

func TestTwoTypesDoNotLeakIntoEachOther(t *testing.T) {
	candidateDir := t.TempDir()
	res, err := Install(candidateDir, testVersion)
	if err != nil {
		t.Fatal(err)
	}
	if err := InstallType(candidateDir, "candidate-interview", testVersion, res); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(candidateDir, ".claude/skills/recruiter-setup")); !os.IsNotExist(err) {
		t.Error("candidate-interview install leaked a recruiter-interview skill")
	}
	if _, err := os.Stat(filepath.Join(candidateDir, "exercises/001-rate-limiter/GRADING.md")); !os.IsNotExist(err) {
		t.Error("candidate-interview install leaked recruiter-interview's GRADING.md")
	}

	recruiterDir := t.TempDir()
	res2, err := Install(recruiterDir, testVersion)
	if err != nil {
		t.Fatal(err)
	}
	if err := InstallType(recruiterDir, "recruiter-interview", testVersion, res2); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(recruiterDir, ".claude/skills/interview-setup")); !os.IsNotExist(err) {
		t.Error("recruiter-interview install leaked a candidate-interview skill")
	}
	if _, err := os.Stat(filepath.Join(recruiterDir, ".agentic-cms/templates/round-plan.md")); !os.IsNotExist(err) {
		t.Error("recruiter-interview install leaked candidate-interview's round-plan.md template")
	}
}

func TestReconcileContentMDRecruiterInterviewTypeSectionMissing(t *testing.T) {
	dir := t.TempDir()
	if _, err := Install(dir, testVersion); err != nil {
		t.Fatal(err)
	}
	if err := ComposeTypeContentMD(dir, "recruiter-interview"); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "CONTENT.md")
	b, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	start := strings.Index(string(b), "## Type: recruiter-interview")
	end := strings.Index(string(b), "## Operations")
	if start == -1 || end == -1 {
		t.Fatal("test setup: composed CONTENT.md missing expected anchors")
	}
	stripped := string(b)[:start] + string(b)[end:]
	if err := os.WriteFile(target, []byte(stripped), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := ReconcileContentMD(dir, "v0.8.2", testVersion, "recruiter-interview")
	if err != nil {
		t.Fatal(err)
	}
	if report == nil {
		t.Fatal("expected a reconciliation report for a missing Type section, got nil")
	}
	found := false
	for _, s := range report.MissingSections {
		if s == "Type: recruiter-interview" {
			found = true
		}
	}
	if !found {
		t.Errorf("missing sections %v do not include Type: recruiter-interview", report.MissingSections)
	}
	after, _ := os.ReadFile(target)
	if string(after) != stripped {
		t.Error("reconciliation modified the user's CONTENT.md")
	}
}
