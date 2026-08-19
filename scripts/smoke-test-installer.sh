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

# Nesting sandboxes inside this repo means a plain (non-git-initialized)
# sandbox would otherwise be discovered by git as "inside" this repo's own
# working tree (git walks up looking for a .git) — InstallGitHook would then
# correctly-but-unwantedly wire its hook into THIS repo's real
# .git/hooks/pre-commit. Capping discovery at SANDBOX_ROOT stops that: a
# sandbox with no .git of its own is treated as outside any repository, while
# a sandbox that explicitly runs `git init` (the 3b scenarios below) is still
# found normally, since GIT_CEILING_DIRECTORIES only stops the upward walk,
# never hides a repo at or below the ceiling itself.
export GIT_CEILING_DIRECTORIES="$SANDBOX_ROOT"
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
# .agentic-cms/VERSION is stamped with the installing binary's version.
BIN_VERSION="$("$BINARY" version | awk '{print $2}')"
sed -i "s/{{VERSION}}/$BIN_VERSION/g" "$EXPECTED/.agentic-cms/VERSION"

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

classify_out="$(.agentic-cms/scripts/ac-classify check docs/smoke/check.md)"
echo "    $classify_out"
assert_contains "$classify_out" '"clean": true' "ac-classify check reports clean on a freshly created C1 page"

python3 -c "
text = open('docs/smoke/check.md').read()
text = text.replace('## Content', '## Content\n\napi_key: sk-abcdef1234567890abcdef')
open('docs/smoke/check.md', 'w').write(text)
"
floor_out="$(.agentic-cms/scripts/ac-classify check docs/smoke/check.md)"
echo "    $floor_out"
assert_contains "$floor_out" '"floor_violation": true' "ac-classify check detects a credential-shaped floor violation"
assert_contains "$floor_out" '"implied_level": "C3"' "ac-classify check implies C3 for a credential-shaped string"

reclass_out="$(.agentic-cms/scripts/ac-page classify docs/smoke/check.md C3)"
echo "    $reclass_out"
assert_contains "$reclass_out" '"classification": "C3"' "ac-page classify raises the rating"
clean_after_out="$(.agentic-cms/scripts/ac-classify check docs/smoke/check.md)"
echo "    $clean_after_out"
assert_contains "$clean_after_out" '"clean": true' "ac-classify check reports clean after re-rating to the implied floor"

# Floor ack: a user-acked false positive stops blocking; the ack is bound to
# the body hash, so it must be honored only while the body is unchanged.
.agentic-cms/scripts/ac-page new doc docs/smoke/pricing.md --title "Public Pricing" --topic smoke >/dev/null
python3 -c "
text = open('docs/smoke/pricing.md').read()
text = text.replace('## Content', '## Content\n\nPublic list price: \$100/month.')
open('docs/smoke/pricing.md', 'w').write(text)
"
floor2_out="$(.agentic-cms/scripts/ac-classify check docs/smoke/pricing.md)"
assert_contains "$floor2_out" '"floor_violation": true' "bare currency figure trips the C2 floor (recall kept)"
ack_out="$(.agentic-cms/scripts/ac-page classify docs/smoke/pricing.md C1 --ack-floor)"
echo "    $ack_out"
assert_contains "$ack_out" '"classification-ack"' "ac-page classify --ack-floor stamps the ack"
acked_out="$(.agentic-cms/scripts/ac-classify check docs/smoke/pricing.md)"
echo "    $acked_out"
assert_contains "$acked_out" '"clean": true' "acked false positive reports clean"
assert_contains "$acked_out" '"acked": true' "acked page reports acked: true"

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

# --- 2b. CONTENT.md schema reconciliation ---
echo
echo "-- CONTENT.md reconciliation --"
RECON="$SANDBOX/reconcile"
mkdir -p "$RECON"
"$BINARY" init "$RECON" >/dev/null
sed -i 's/^## Classification$/## My Local Rules/' "$RECON/CONTENT.md"
out="$("$BINARY" init "$RECON" 2>&1)"
printf '%s\n' "$out" | sed 's/^/    /'
assert_contains "$out" "CONTENT.md reconciliation" "re-init over a customized CONTENT.md emits a reconciliation report"
assert_contains "$out" "Classification" "reconciliation report names the missing section"
if [[ -f "$RECON/.agentic-cms/CONTENT.upstream.md" ]]; then
    pass "upstream sidecar written"
else
    fail "upstream sidecar written — .agentic-cms/CONTENT.upstream.md missing"
fi
assert_contains "$(cat "$RECON/CONTENT.md")" "## My Local Rules" "user's customized CONTENT.md left untouched"
assert_no_line_matches "$(cat "$RECON/CONTENT.md")" "## Classification" "reconciliation did not edit the user's CONTENT.md"

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

# --- 3b. Git pre-commit hook: fresh install, existing-hook append, core.hooksPath report ---
echo
echo "-- git pre-commit hook install --"
if command -v git >/dev/null 2>&1; then
    GIT1="$SANDBOX/git-fresh"
    mkdir -p "$GIT1"
    git -C "$GIT1" init -q
    out="$("$BINARY" init "$GIT1" 2>&1)"
    printf '%s\n' "$out" | sed 's/^/    /'
    assert_contains "$out" "  created  git hook: pre-commit (installed)" "fresh git repo: pre-commit hook installed"
    if [[ -x "$GIT1/.git/hooks/pre-commit" ]]; then
        pass "installed pre-commit hook is executable"
    else
        fail "installed pre-commit hook is executable — not executable or missing"
    fi
    out2="$("$BINARY" init "$GIT1" 2>&1)"
    assert_contains "$out2" "  skipped  git hook: pre-commit (block already present)" "re-init skips, does not re-merge the git hook"

    GIT2="$SANDBOX/git-existing-hook"
    mkdir -p "$GIT2/.git/hooks"
    git -C "$GIT2" init -q
    printf '#!/usr/bin/env bash\necho existing\n' > "$GIT2/.git/hooks/pre-commit"
    chmod +x "$GIT2/.git/hooks/pre-commit"
    out="$("$BINARY" init "$GIT2" 2>&1)"
    printf '%s\n' "$out" | sed 's/^/    /'
    assert_contains "$out" "  merged   git hook: pre-commit (block appended)" "existing hook: managed block appended"
    assert_contains "$(cat "$GIT2/.git/hooks/pre-commit")" "echo existing" "existing hook content preserved after append"

    GIT3="$SANDBOX/git-hookspath"
    mkdir -p "$GIT3/custom-hooks"
    git -C "$GIT3" init -q
    git -C "$GIT3" config core.hooksPath custom-hooks
    out="$("$BINARY" init "$GIT3" 2>&1)"
    printf '%s\n' "$out" | sed 's/^/    /'
    assert_contains "$out" "core.hooksPath=custom-hooks" "core.hooksPath set: reported, not written"
    if [[ ! -e "$GIT3/.git/hooks/pre-commit" ]] && [[ -z "$(ls -A "$GIT3/custom-hooks" 2>/dev/null)" ]]; then
        pass "core.hooksPath set: nothing written to either location"
    else
        fail "core.hooksPath set: something was written despite the report-only contract"
    fi
else
    echo "    git not on PATH — skipping git pre-commit hook scenarios"
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
