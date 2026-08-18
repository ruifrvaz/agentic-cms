#!/usr/bin/env bash
# End-to-end smoke test for the agentic-cms installer binary.
#
# Usage: scripts/smoke-test-installer.sh [path-to-binary]
#
# Runs the real compiled binary against mktemp sandboxes and verifies:
#   - a greenfield install matches scaffold/tree/ (after {{DATE}} substitution)
#   - a second run is idempotent (no re-creates, no overwrites)
#   - a brownfield CLAUDE.md gets the managed block merged, not duplicated
#   - init against an invalid target directory fails with a non-zero exit
#   - --version / --help behave
#
# Set KEEP_SMOKE_DIR=1 to keep the sandbox directory for inspection instead of
# removing it on exit.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
TREE_DIR="$REPO_ROOT/scaffold/tree"

BINARY="${1:-$REPO_ROOT/installer/agentic-cms}"
if [[ ! -x "$BINARY" ]]; then
    echo "error: binary not found or not executable: $BINARY" >&2
    echo "hint: run 'make build' first" >&2
    exit 1
fi
BINARY="$(cd "$(dirname "$BINARY")" && pwd)/$(basename "$BINARY")"

# Sandboxes live inside the repo (gitignored) rather than system /tmp so
# KEEP_SMOKE_DIR=1 leaves an inspectable, project-local artifact.
SANDBOX_ROOT="$REPO_ROOT/.smoke-test"
mkdir -p "$SANDBOX_ROOT"
SANDBOX="$(mktemp -d "$SANDBOX_ROOT/run.XXXXXX")"
cleanup() {
    if [[ -n "${KEEP_SMOKE_DIR:-}" ]]; then
        echo "KEEP_SMOKE_DIR set — leaving sandbox at $SANDBOX"
        return
    fi
    case "$SANDBOX" in
        "$SANDBOX_ROOT"/run.*) rm -rf "$SANDBOX" ;;
        *) echo "warning: refusing to remove unexpected path: $SANDBOX" >&2 ;;
    esac
}
trap cleanup EXIT

FAILURES=0

fail() {
    echo "[FAIL] $1" >&2
    FAILURES=$((FAILURES + 1))
}

pass() {
    echo "[OK] $1"
}

assert_contains() {
    local haystack="$1" needle="$2" msg="$3"
    if [[ "$haystack" == *"$needle"* ]]; then pass "$msg"; else fail "$msg — did not find: $needle"; fi
}

# Counts lines starting with $prefix (Result.Print() indents each entry with
# two spaces, so this avoids false positives from the "Done: N created, M
# merged, K skipped." summary line, which has no leading indent).
assert_no_line_matches() {
    local text="$1" prefix="$2" msg="$3"
    local count
    count="$(printf '%s\n' "$text" | grep -c "^$prefix" || true)"
    if [[ "$count" -eq 0 ]]; then pass "$msg"; else fail "$msg — found $count matching line(s)"; fi
}

assert_tree_matches() {
    local expected="$1" actual="$2" msg="$3"
    local diff_out
    if diff_out="$(diff -rq "$expected" "$actual" 2>&1)"; then
        pass "$msg"
    else
        fail "$msg"
        printf '%s\n' "$diff_out" | sed 's/^/    /' >&2
    fi
}

TODAY="$(date +%Y-%m-%d)"

echo "== agentic-cms installer smoke test =="
echo "binary:  $BINARY"
echo "sandbox: $SANDBOX"
echo

# Expected greenfield tree: scaffold/tree/ with {{DATE}} resolved to today,
# matching the substitution Install() performs on every embedded file — except
# .agentic-cms/templates/ (keeps {{DATE}} live for ac-page new to fill in at
# page-creation time) and .agentic-cms/scripts/ (source code referencing the
# literal placeholder string, not a page carrying it).
EXPECTED="$SANDBOX/expected"
cp -r "$TREE_DIR" "$EXPECTED"
find "$EXPECTED" -type f \
    -not -path "$EXPECTED/.agentic-cms/templates/*" \
    -not -path "$EXPECTED/.agentic-cms/scripts/*" \
    -print0 | xargs -0 sed -i "s/{{DATE}}/$TODAY/g"

# --- 1. Greenfield install ---
echo "-- greenfield install --"
GREEN="$SANDBOX/greenfield"
mkdir -p "$GREEN"
out="$("$BINARY" init "$GREEN" 2>&1)"
printf '%s\n' "$out" | sed 's/^/    /'
assert_no_line_matches "$out" "  skipped  " "greenfield install reports no skips"
assert_tree_matches "$EXPECTED" "$GREEN" "greenfield tree matches scaffold/tree/ (post {{DATE}} substitution)"
assert_contains "$(cat "$GREEN/.agentic-cms/templates/doc.md")" "{{TITLE}}" "doc.md template keeps non-date placeholders"
assert_contains "$(cat "$GREEN/.agentic-cms/templates/doc.md")" "{{DATE}}" "doc.md template keeps the {{DATE}} placeholder live"
if [[ -x "$GREEN/.agentic-cms/scripts/ac-index" ]]; then
    pass "ac-index installed executable"
else
    fail "ac-index installed executable — not executable or missing"
fi

# --- 1b. Toolkit functional check ---
# Exercises the installed ac-* scripts against the real filesystem output of
# Install(), not just scaffold/tree/ in isolation — this is what would have
# caught the installer's {{DATE}} substitution corrupting ac-page's own source
# (the dict key "{{DATE}}" inside .agentic-cms/scripts/ac-page).
echo
echo "-- toolkit functional check --"
pushd "$GREEN" >/dev/null

page_out="$(.agentic-cms/scripts/ac-page new doc docs/smoke/check.md --title "Smoke Check" --topic smoke)"
echo "    $page_out"
if [[ "$page_out" == *'"ok": true'* ]]; then pass "ac-page new returns ok"; else fail "ac-page new returns ok — got: $page_out"; fi
page_body="$(cat docs/smoke/check.md)"
assert_contains "$page_body" "created: $TODAY" "ac-page new fills {{DATE}} with today's date"
if [[ "$page_body" == *'{{DATE}}'* ]]; then fail "ac-page new left a raw {{DATE}} placeholder unfilled"; else pass "ac-page new left no unfilled {{DATE}} placeholder"; fi

idx_out="$(.agentic-cms/scripts/ac-index add topics docs/smoke/check.md "smoke test page")"
echo "    $idx_out"
if [[ "$idx_out" == *'"ok": true'* ]]; then pass "ac-index add returns ok"; else fail "ac-index add returns ok — got: $idx_out"; fi

log_out="$(.agentic-cms/scripts/ac-log append new-item "smoke/check")"
echo "    $log_out"
if [[ "$log_out" == *'"ok": true'* ]]; then pass "ac-log append returns ok"; else fail "ac-log append returns ok — got: $log_out"; fi

check_out="$(.agentic-cms/scripts/ac-index check)"
echo "    $check_out"
assert_contains "$check_out" '"clean": true' "ac-index check reports clean after add"

links_out="$(.agentic-cms/scripts/ac-links check)"
echo "    $links_out"
assert_contains "$links_out" '"clean": true' "ac-links check reports clean"

popd >/dev/null

# --- 2. Idempotent re-run ---
echo
echo "-- idempotent re-run --"
IDEM="$SANDBOX/idempotent"
mkdir -p "$IDEM"
"$BINARY" init "$IDEM" >/dev/null
echo "user content" > "$IDEM/wiki/index.md"
out="$("$BINARY" init "$IDEM" 2>&1)"
printf '%s\n' "$out" | sed 's/^/    /'
assert_no_line_matches "$out" "  created  " "second install creates no files"
assert_contains "$(cat "$IDEM/wiki/index.md")" "user content" "mutated file survives re-init"

# --- 3. Brownfield CLAUDE.md merge ---
echo
echo "-- brownfield CLAUDE.md merge --"
BROWN="$SANDBOX/brownfield"
mkdir -p "$BROWN"
printf '# My project\n\nExisting instructions.\n' > "$BROWN/CLAUDE.md"
out="$("$BINARY" init "$BROWN" 2>&1)"
printf '%s\n' "$out" | sed 's/^/    /'
assert_contains "$out" "  merged   CLAUDE.md (block appended)" "first brownfield run merges the block"
claude_content="$(cat "$BROWN/CLAUDE.md")"
assert_contains "$claude_content" "# My project" "original CLAUDE.md content preserved"
assert_contains "$claude_content" "<!-- agentic-cms:begin -->" "managed block begin marker present"
assert_contains "$claude_content" "<!-- agentic-cms:end -->" "managed block end marker present"

out2="$("$BINARY" init "$BROWN" 2>&1)"
assert_contains "$out2" "  skipped  CLAUDE.md (block already present)" "second brownfield run skips, does not re-merge"
marker_count="$(grep -o '<!-- agentic-cms:begin -->' "$BROWN/CLAUDE.md" | wc -l)"
if [[ "$marker_count" -eq 1 ]]; then
    pass "managed block not duplicated on re-init"
else
    fail "managed block duplicated on re-init (found $marker_count begin markers)"
fi

# --- 4. Invalid target directory ---
echo
echo "-- invalid target directory --"
set +e
"$BINARY" init "$SANDBOX/does-not-exist" >/dev/null 2>"$SANDBOX/stderr.txt"
status=$?
set -e
if [[ "$status" -ne 0 ]]; then pass "init against a missing directory exits non-zero"; else fail "init against a missing directory should fail"; fi
assert_contains "$(cat "$SANDBOX/stderr.txt")" "agentic-cms:" "error message is reported on stderr"

# --- 5. --version / --help sanity ---
echo
echo "-- version / help --"
version_out="$("$BINARY" --version)"
assert_contains "$version_out" "agentic-cms" "--version reports the binary name"
help_out="$("$BINARY" --help)"
assert_contains "$help_out" "Usage:" "--help prints usage"

echo
if [[ "$FAILURES" -eq 0 ]]; then
    echo "== all checks passed =="
    exit 0
else
    echo "== $FAILURES check(s) failed ==" >&2
    exit 1
fi
