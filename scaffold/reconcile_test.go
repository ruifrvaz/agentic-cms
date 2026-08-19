package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallStampsVersion(t *testing.T) {
	dir := t.TempDir()
	if _, err := Install(dir, "v0.7.0"); err != nil {
		t.Fatal(err)
	}
	if got := InstalledVersion(dir); got != "v0.7.0" {
		t.Errorf("InstalledVersion = %q, want v0.7.0", got)
	}
	// A later install with a newer binary must restamp, not skip.
	res, err := Install(dir, "v0.8.0")
	if err != nil {
		t.Fatal(err)
	}
	if got := InstalledVersion(dir); got != "v0.8.0" {
		t.Errorf("re-install left VERSION at %q, want v0.8.0", got)
	}
	want := filepath.FromSlash(".agentic-cms/VERSION")
	for _, f := range res.Skipped {
		if f == want+" (exists)" {
			t.Errorf("%s reported as skipped, should have been restamped", want)
		}
	}
	// The placeholder must never leak into the installed file.
	b, _ := os.ReadFile(filepath.Join(dir, want))
	if strings.Contains(string(b), "{{VERSION}}") {
		t.Error("installed VERSION still contains the {{VERSION}} placeholder")
	}
}

func TestReconcileContentMDFreshInstallClean(t *testing.T) {
	dir := t.TempDir()
	if _, err := Install(dir, testVersion); err != nil {
		t.Fatal(err)
	}
	report, err := ReconcileContentMD(dir, "", testVersion)
	if err != nil {
		t.Fatal(err)
	}
	if report != nil {
		t.Errorf("fresh install reported missing sections: %v", report.MissingSections)
	}
	if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(UpstreamSidecar))); err == nil {
		t.Error("sidecar written even though nothing was missing")
	}
}

func TestReconcileContentMDMissingSection(t *testing.T) {
	dir := t.TempDir()
	if _, err := Install(dir, testVersion); err != nil {
		t.Fatal(err)
	}
	// Simulate a pre-classification install: the user's customized CONTENT.md
	// lacks the upstream Classification section.
	target := filepath.Join(dir, "CONTENT.md")
	b, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	stripped := strings.Replace(string(b), "## Classification", "## My local notes", 1)
	if stripped == string(b) {
		t.Fatal("test setup: CONTENT.md has no '## Classification' heading to strip")
	}
	if err := os.WriteFile(target, []byte(stripped), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := ReconcileContentMD(dir, "v0.6.2", "v0.7.0")
	if err != nil {
		t.Fatal(err)
	}
	if report == nil {
		t.Fatal("expected a reconciliation report, got nil")
	}
	found := false
	for _, s := range report.MissingSections {
		if s == "Classification" {
			found = true
		}
	}
	if !found {
		t.Errorf("missing sections %v do not include Classification", report.MissingSections)
	}
	if report.InstalledVersion != "v0.6.2" || report.ShippedVersion != "v0.7.0" {
		t.Errorf("report versions = %q → %q", report.InstalledVersion, report.ShippedVersion)
	}
	// Sidecar written with the full upstream copy; user's CONTENT.md untouched.
	sb, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(UpstreamSidecar)))
	if err != nil {
		t.Fatalf("sidecar not written: %v", err)
	}
	if !strings.Contains(string(sb), "## Classification") {
		t.Error("sidecar lacks the upstream Classification section")
	}
	after, _ := os.ReadFile(target)
	if string(after) != stripped {
		t.Error("reconciliation modified the user's CONTENT.md")
	}
}
