package scaffold

// Fixture tests for the ac-classify engine and the pre-commit gate's
// delta-scoping — exercised through the real installed scripts, since the
// engine is bash+python and every enforcement point is a thin caller of it.

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func requireTools(t *testing.T, tools ...string) {
	t.Helper()
	for _, tool := range tools {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("%s not available", tool)
		}
	}
}

// installEngine installs the scaffold into a temp dir and returns it.
func installEngine(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if _, err := Install(dir, testVersion); err != nil {
		t.Fatal(err)
	}
	return dir
}

func writePage(t *testing.T, dir, rel, classification, body string) {
	t.Helper()
	target := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	fm := "---\ntitle: Fixture\n"
	if classification != "" {
		fm += "classification: " + classification + "\n"
	}
	fm += "---\n"
	if err := os.WriteFile(target, []byte(fm+body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// runScript runs an installed .agentic-cms script with dir as cwd and returns
// combined stdout. Non-zero exit is fatal unless allowFail is true.
func runScript(t *testing.T, dir string, allowFail bool, script string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(filepath.Join(dir, filepath.FromSlash(script)), args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	code := 0
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("%s %v: %v", script, args, err)
		}
		code = ee.ExitCode()
		if !allowFail {
			t.Fatalf("%s %v exited %d: %s\n%s", script, args, code, out, ee.Stderr)
		}
	}
	return string(out), code
}

// classifyCheck runs ac-classify check on one page and returns its entry.
func classifyCheck(t *testing.T, dir, page string) map[string]any {
	t.Helper()
	out, _ := runScript(t, dir, false, ".agentic-cms/scripts/ac-classify", "check", page)
	var result struct {
		Pages []map[string]any `json:"pages"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("unparseable ac-classify output %q: %v", out, err)
	}
	if len(result.Pages) != 1 {
		t.Fatalf("expected 1 page entry, got %d", len(result.Pages))
	}
	if e, ok := result.Pages[0]["error"]; ok {
		t.Fatalf("ac-classify error for %s: %v", page, e)
	}
	return result.Pages[0]
}

// TestCurrencyFloorRecall pins the floor's recall contract: every enumerated
// currency shape — bare, with no finance context required — implies at least
// C2. Recall is never traded for precision; false positives are resolved by
// the user-only ack, tested below.
func TestCurrencyFloorRecall(t *testing.T) {
	requireTools(t, "python3", "bash")
	dir := installEngine(t)

	cases := []struct{ name, body string }{
		{"symbol-prefixed", "Vendor list price: $100/month per seat.\n"},
		{"symbol-prefixed-euro", "Public plan: €9 per user.\n"},
		{"symbol-suffixed", "Public plan: 9,50€ per user.\n"},
		{"iso-prefixed", "Budget line USD 1,000 approved.\n"},
		{"iso-suffixed", "Contract value 250000 EUR total.\n"},
		{"word-form", "They paid 1,000 dollars for the license.\n"},
		{"cap-table", "Alice Example holds 40%, invested $2,500,000 at seed.\n"},
	}
	for _, c := range cases {
		page := "docs/fixtures/" + c.name + ".md"
		writePage(t, dir, page, "C1", c.body)
		entry := classifyCheck(t, dir, page)
		if entry["floor"] != "C2" {
			t.Errorf("%s: floor = %v, want C2 (body %q)", c.name, entry["floor"], c.body)
		}
		if entry["floor_violation"] != true {
			t.Errorf("%s: floor_violation = %v, want true", c.name, entry["floor_violation"])
		}
	}

	// No-currency control: prose with plain numbers must not trip.
	writePage(t, dir, "docs/fixtures/control.md", "C1", "We evaluated 100 vendors over 12 months.\n")
	entry := classifyCheck(t, dir, "docs/fixtures/control.md")
	if entry["floor"] != nil {
		t.Errorf("control page tripped floor %v on plain numbers", entry["floor"])
	}
}

// TestFloorAckLifecycle pins the ack mechanism: a user ack silences the floor
// violation for the acked body, any body change invalidates it, and a plain
// re-rate withdraws it.
func TestFloorAckLifecycle(t *testing.T) {
	requireTools(t, "python3", "bash")
	dir := installEngine(t)
	page := "docs/pricing/vendors.md"
	writePage(t, dir, page, "C1", "Public vendor pricing: $100/month, $250/month tiers.\n")

	entry := classifyCheck(t, dir, page)
	if entry["floor_violation"] != true {
		t.Fatalf("pre-ack: floor_violation = %v, want true", entry["floor_violation"])
	}

	// User acks the false positive, keeping C1.
	out, _ := runScript(t, dir, false, ".agentic-cms/scripts/ac-page", "classify", page, "C1", "--ack-floor")
	if !strings.Contains(out, "classification-ack") {
		t.Fatalf("ack output lacks classification-ack: %s", out)
	}
	entry = classifyCheck(t, dir, page)
	if entry["floor_violation"] != false {
		t.Errorf("acked: floor_violation = %v, want false", entry["floor_violation"])
	}
	if entry["acked"] != true {
		t.Errorf("acked: acked = %v, want true", entry["acked"])
	}
	if entry["floor"] != "C2" {
		t.Errorf("acked: floor = %v — the hit itself must stay visible", entry["floor"])
	}

	// Any body change invalidates the ack: the violation comes back.
	b, _ := os.ReadFile(filepath.Join(dir, filepath.FromSlash(page)))
	if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(page)),
		append(b, []byte("New tier: $900/month.\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	entry = classifyCheck(t, dir, page)
	if entry["floor_violation"] != true {
		t.Errorf("post-edit: floor_violation = %v, want true (ack must not survive a body change)", entry["floor_violation"])
	}
	if entry["acked"] != false {
		t.Errorf("post-edit: acked = %v, want false", entry["acked"])
	}

	// Re-ack, then a plain re-rate must withdraw the standing ack.
	runScript(t, dir, false, ".agentic-cms/scripts/ac-page", "classify", page, "C1", "--ack-floor")
	runScript(t, dir, false, ".agentic-cms/scripts/ac-page", "classify", page, "C1")
	entry = classifyCheck(t, dir, page)
	if entry["acked"] != false || entry["floor_violation"] != true {
		t.Errorf("post-re-rate: acked = %v, floor_violation = %v — plain classify must withdraw the ack",
			entry["acked"], entry["floor_violation"])
	}
}

// gitIn runs git with args in dir.
func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	base := []string{"-c", "user.email=test@example.invalid", "-c", "user.name=test"}
	cmd := exec.Command("git", append(base, args...)...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

// runPreCommit runs the installed pre-commit gate in dir and returns exit
// code and stderr.
func runPreCommit(t *testing.T, dir string) (int, string) {
	t.Helper()
	cmd := exec.Command(filepath.Join(dir, ".agentic-cms", "hooks", "pre-commit"))
	cmd.Dir = dir
	var stderr strings.Builder
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if err != nil {
		ee, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("pre-commit: %v", err)
		}
		code = ee.ExitCode()
	}
	return code, stderr.String()
}

// TestPreCommitDeltaScope pins the gate's brownfield behavior: pre-existing
// drift elsewhere in the tree is one summarized warning, only the staged
// delta blocks, and bleed blocks whenever the staged index/log itself leaks —
// regardless of which page the leak came from.
func TestPreCommitDeltaScope(t *testing.T) {
	requireTools(t, "python3", "bash", "git")
	dir := installEngine(t)
	gitIn(t, dir, "init", "-q")

	// Pre-existing backlog: a committed floor violation.
	writePage(t, dir, "docs/legacy/dirty.md", "C1", "Legacy cap table: $2,500,000 seed round.\n")
	gitIn(t, dir, "add", "-A")
	gitIn(t, dir, "commit", "-q", "-m", "seed", "--no-verify")

	// Scenario 1: stage one clean unrelated file — must pass with a single
	// backlog summary, not per-file noise.
	writePage(t, dir, "docs/notes/clean.md", "C1", "Meeting notes with no sensitive shapes.\n")
	gitIn(t, dir, "add", "docs/notes/clean.md")
	code, stderr := runPreCommit(t, dir)
	if code != 0 {
		t.Fatalf("clean staged delta blocked by pre-existing drift (exit %d):\n%s", code, stderr)
	}
	if !strings.Contains(stderr, "pre-existing classification issue(s) elsewhere in the tree") {
		t.Errorf("expected one backlog summary warning, got:\n%s", stderr)
	}
	if strings.Contains(stderr, "docs/legacy/dirty.md: rated") {
		t.Errorf("backlog page listed per-file instead of summarized:\n%s", stderr)
	}

	// Scenario 2: edit and stage the dirty page itself — must block. (A plain
	// re-add of an unchanged committed file is index==HEAD: no staged delta,
	// and the gate correctly ignores it.)
	dirty := filepath.Join(dir, "docs", "legacy", "dirty.md")
	db, _ := os.ReadFile(dirty)
	if err := os.WriteFile(dirty, append(db, []byte("Series A: $8,000,000.\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, dir, "add", "docs/legacy/dirty.md")
	code, stderr = runPreCommit(t, dir)
	if code == 0 {
		t.Fatalf("staged floor violation was not blocked:\n%s", stderr)
	}
	if !strings.Contains(stderr, "docs/legacy/dirty.md") || !strings.Contains(stderr, "BLOCKED") {
		t.Errorf("block message missing staged page:\n%s", stderr)
	}
	gitIn(t, dir, "reset", "-q", "docs/legacy/dirty.md")

	// Scenario 3: an acked false positive passes the gate.
	writePage(t, dir, "docs/pricing/public.md", "C1", "Public list price: $100/month.\n")
	runScript(t, dir, false, ".agentic-cms/scripts/ac-page", "classify", "docs/pricing/public.md", "C1", "--ack-floor")
	gitIn(t, dir, "add", "docs/pricing/public.md")
	code, stderr = runPreCommit(t, dir)
	if code != 0 {
		t.Fatalf("acked false positive still blocked (exit %d):\n%s", code, stderr)
	}
	gitIn(t, dir, "commit", "-q", "-m", "acked", "--no-verify")

	// Scenario 4: bleed into a STAGED wiki/index.md blocks, even though the
	// leaking source page (the committed C2 page) is not staged.
	writePage(t, dir, "docs/finance/budget.md", "C2", "FY26 budget details.\n")
	gitIn(t, dir, "add", "docs/finance/budget.md")
	gitIn(t, dir, "commit", "-q", "-m", "budget", "--no-verify")
	index := filepath.Join(dir, "wiki", "index.md")
	b, _ := os.ReadFile(index)
	leak := string(b) + fmt.Sprintf("- [Budget](docs/finance/budget.md) — total is $4,000,000\n")
	if err := os.WriteFile(index, []byte(leak), 0o644); err != nil {
		t.Fatal(err)
	}
	gitIn(t, dir, "add", "wiki/index.md")
	code, stderr = runPreCommit(t, dir)
	if code == 0 {
		t.Fatalf("bleed into staged wiki/index.md was not blocked:\n%s", stderr)
	}
	if !strings.Contains(stderr, "wiki/index.md") || !strings.Contains(stderr, "leaks") {
		t.Errorf("bleed block message malformed:\n%s", stderr)
	}
}
