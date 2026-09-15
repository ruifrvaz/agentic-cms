package scaffold

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// unmangleEmbeddedPath reverses the two source-tree-only renames a type's
// non-markdown payload needs so go:embed picks it up without also polluting
// this repo's own go build/vet/test ./... — an installed project always sees
// the real names:
//
//   - "go.mod" is authored as "go.mod.embedded". A file literally named
//     go.mod makes go:embed treat that subtree as a separate module and
//     silently exclude it, regardless of the all: prefix.
//   - "exercises" is authored as "_exercises". go build/vet/test ./... skips
//     any path with a leading "_"/"." segment, so this repo's own test run
//     never compiles or executes a coding exercise's deliberately-unfinished
//     (panics until fixed) sample code. go:embed's all: prefix still reaches
//     past that same leading-underscore exclusion to embed it.
func unmangleEmbeddedPath(rel string) string {
	sep := string(filepath.Separator)
	switch {
	case rel == "_exercises":
		rel = "exercises"
	case strings.HasPrefix(rel, "_exercises"+sep):
		rel = "exercises" + sep + strings.TrimPrefix(rel, "_exercises"+sep)
	}
	if filepath.Base(rel) == "go.mod.embedded" {
		rel = filepath.Join(filepath.Dir(rel), "go.mod")
	}
	return rel
}

// InstallType overlays the named content type's tree onto dir, on top of a
// prior base Install. Every file a type ships is framework-owned — always
// overwritten, exactly like the base scaffold's own skills/templates — since
// a type has no user-authored content of its own: docs/ and wiki/ ship empty,
// and the type's skills create every page per engagement, not the installer.
// Appends to res so the caller reports overlay files alongside the base
// tree's. version stamps TYPE.md's base: field, mirroring how Install stamps
// .agentic-cms/VERSION.
func InstallType(dir, name, version string, res *Result) error {
	if !HasType(name) {
		return fmt.Errorf("unknown content type %q", name)
	}
	root := typeTreeRoot(name)
	date := time.Now().Format("2006-01-02")

	return fs.WalkDir(Types, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		rel = unmangleEmbeddedPath(rel)
		target := filepath.Join(dir, filepath.FromSlash(rel))

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		data, err := Types.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(data)
		// Templates keep {{DATE}} live for ac-page (page-creation time, not
		// install time); exercises/ is non-markdown payload, same reasoning
		// as the base tree's .agentic-cms/scripts/ exclusion. TYPE.md's own
		// prose names {{DATE}} as one of ac-page's placeholders (documentation
		// for a human/agent reader, not a live placeholder) — substituting it
		// there corrupts that sentence into today's literal date.
		if !strings.HasPrefix(rel, filepath.FromSlash(".agentic-cms/templates/")) &&
			!strings.HasPrefix(rel, filepath.FromSlash("exercises/")) &&
			rel != filepath.FromSlash(typeManifestFile) {
			content = strings.ReplaceAll(content, "{{DATE}}", date)
		}
		if rel == filepath.FromSlash(typeManifestFile) {
			content = strings.ReplaceAll(content, "{{VERSION}}", version)
		}

		_, statErr := os.Stat(target)
		exists := statErr == nil
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
			return err
		}
		if exists {
			res.Updated = append(res.Updated, rel)
		} else {
			res.Created = append(res.Created, rel)
		}
		return nil
	})
}

// typeMarkerBegin is the per-type CLAUDE.md managed block's begin marker —
// the type's own append-once marker, alongside the base <!--
// agentic-cms:begin/end -->. Only the begin marker is checked for presence
// (mirroring installClaudeMD's own base-block check), so no matching
// typeMarkerEnd is needed.
func typeMarkerBegin(name string) string { return "<!-- agentic-cms:type:" + name + ":begin -->" }

// InstallTypeClaudeMD appends the type's CLAUDE.md managed block if not
// already present, mirroring installClaudeMD's append-once semantics for the
// base block. Assumes CLAUDE.md already exists (Install always creates or
// merges it before this runs).
func InstallTypeClaudeMD(dir, name string, res *Result) error {
	block, err := Types.ReadFile(typesRoot + "/" + name + "/claude.fragment.md")
	if err != nil {
		return err
	}
	target := filepath.Join(dir, "CLAUDE.md")
	existing, err := os.ReadFile(target)
	if err != nil {
		return err
	}
	if strings.Contains(string(existing), typeMarkerBegin(name)) {
		res.Skipped = append(res.Skipped, "CLAUDE.md (type block already present)")
		return nil
	}
	out := strings.TrimRight(string(existing), "\n") + "\n\n" + string(block)
	if err := os.WriteFile(target, []byte(out), 0o644); err != nil {
		return err
	}
	res.Merged = append(res.Merged, "CLAUDE.md (type block appended)")
	return nil
}

// ComposeTypeContentMD rewrites dir's CONTENT.md — freshly written by Install
// moments earlier — with the type's fragment applied on top. Only call this
// right after Install created CONTENT.md for the first time under an active
// type; an already-existing CONTENT.md is never rewritten here (see
// ReconcileContentMD for the re-init path, which reports drift instead of
// editing the user's file).
func ComposeTypeContentMD(dir, name string) error {
	frag, err := Types.ReadFile(typesRoot + "/" + name + "/content.fragment.md")
	if err != nil {
		return err
	}
	target := filepath.Join(dir, "CONTENT.md")
	base, err := os.ReadFile(target)
	if err != nil {
		return err
	}
	composed, err := applyTypeFragment(string(base), string(frag))
	if err != nil {
		return fmt.Errorf("composing CONTENT.md for type %q: %w", name, err)
	}
	return os.WriteFile(target, []byte(composed), 0o644)
}

// typeAnchorRe matches one fragment directive line: an exact base-file line
// to locate, and whether the fragment content that follows (up to the next
// directive or EOF) is inserted after it, before it, or replaces it.
var typeAnchorRe = regexp.MustCompile(`(?m)^<!-- ac-type-anchor: (after|before|replace) "(.*)" -->$`)

type typeAnchorEdit struct {
	mode    string // after | before | replace
	anchor  string // exact line to locate in the base text
	content string // lines to insert or substitute
}

// parseTypeFragment splits a fragment file into its ordered anchor edits.
func parseTypeFragment(frag string) []typeAnchorEdit {
	matches := typeAnchorRe.FindAllStringSubmatchIndex(frag, -1)
	edits := make([]typeAnchorEdit, 0, len(matches))
	for i, m := range matches {
		mode := frag[m[2]:m[3]]
		anchor := frag[m[4]:m[5]]
		start := m[1]
		end := len(frag)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}
		content := frag[start:end]
		content = strings.TrimPrefix(content, "\n")
		content = strings.TrimSuffix(content, "\n")
		edits = append(edits, typeAnchorEdit{mode: mode, anchor: anchor, content: content})
	}
	return edits
}

// applyTypeFragment applies an anchor-based fragment (see parseTypeFragment)
// to base text, returning the composed result. Every anchor line must be
// found exactly once in base — a missing anchor means the fragment and the
// base schema have drifted apart, which is an authoring error to fix, not a
// case to silently skip.
func applyTypeFragment(base, frag string) (string, error) {
	lines := strings.Split(base, "\n")
	for _, edit := range parseTypeFragment(frag) {
		idx := -1
		for i, line := range lines {
			if line == edit.anchor {
				idx = i
				break
			}
		}
		if idx == -1 {
			return "", fmt.Errorf("anchor not found: %q", edit.anchor)
		}
		contentLines := strings.Split(edit.content, "\n")
		switch edit.mode {
		case "after":
			lines = spliceLines(lines, idx+1, 0, contentLines)
		case "before":
			lines = spliceLines(lines, idx, 0, contentLines)
		case "replace":
			lines = spliceLines(lines, idx, 1, contentLines)
		default:
			return "", fmt.Errorf("unknown anchor mode %q", edit.mode)
		}
	}
	return strings.Join(lines, "\n"), nil
}

// spliceLines removes `remove` lines at index `at` and inserts `insert`
// there, returning a new slice.
func spliceLines(lines []string, at, remove int, insert []string) []string {
	out := make([]string, 0, len(lines)-remove+len(insert))
	out = append(out, lines[:at]...)
	out = append(out, insert...)
	out = append(out, lines[at+remove:]...)
	return out
}
