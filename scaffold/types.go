// Package scaffold: content type overlays. A content type layers a domain on
// top of the base scaffold — its own templates, skills, CONTENT.md/CLAUDE.md
// additions, and (optionally) non-markdown payload — without changing the
// base at all. Types are independent: each lives in its own directory under
// Types and is walked separately; nothing is shared or generated between
// them.
package scaffold

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Types holds every embedded content type, keyed by name under typesRoot.
// The all: prefix mirrors Tree's, so a type's dot-directories (.agentic-cms,
// .claude) are included.
//
//go:embed all:types
var Types embed.FS

const typesRoot = "types"

// typeManifestFile is a type's manifest, relative to its own tree root.
const typeManifestFile = ".agentic-cms/TYPE.md"

// typeManifestPath is where an installed project records its active type.
var typeManifestPath = filepath.FromSlash(".agentic-cms/TYPE.md")

// typeTreeRoot returns the embedded tree root for the named type.
func typeTreeRoot(name string) string {
	return typesRoot + "/" + name + "/tree"
}

// AvailableTypes returns the embedded content type names, sorted.
func AvailableTypes() ([]string, error) {
	entries, err := fs.ReadDir(Types, typesRoot)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

// HasType reports whether name is an embedded content type.
func HasType(name string) bool {
	_, err := fs.Stat(Types, typeTreeRoot(name))
	return err == nil
}

// InstalledType returns the type name recorded in dir's .agentic-cms/TYPE.md,
// or "" if the project is untyped (file absent) or the field can't be read.
// init/update call this to honor an already-installed type when --type is
// omitted, instead of silently falling back to a typeless run.
func InstalledType(dir string) string {
	b, err := os.ReadFile(filepath.Join(dir, typeManifestPath))
	if err != nil {
		return ""
	}
	return frontmatterField(string(b), "name")
}

// frontmatterField extracts one top-level scalar "key: value" line from a
// YAML frontmatter block delimited by "---" lines. This is deliberately not
// a general YAML parser — TYPE.md's manifest is flat key/value plus simple
// "- item" lists that nothing in this codebase needs to read back yet — so a
// few lines of stdlib string handling covers it without taking on a new Go
// dependency, consistent with this project's stdlib-only convention.
func frontmatterField(content, key string) string {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return ""
	}
	prefix := key + ":"
	for _, line := range lines[1:] {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			break
		}
		if val, ok := strings.CutPrefix(trimmed, prefix); ok {
			return strings.Trim(strings.TrimSpace(val), `"`)
		}
	}
	return ""
}
