package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// githubRelease holds the fields we care about from the GitHub releases API.
type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// runUpdate self-updates the binary to the latest GitHub release, then
// re-initializes the project directory if it looks like an installed project.
func runUpdate(localVersion string) {
	projectDir := resolveDefaultProjectDir(".")

	release, err := fetchLatestRelease()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching latest release: %v\n", err)
		os.Exit(1)
	}

	remoteVersion := strings.TrimPrefix(release.TagName, "v")
	localVersionTrimmed := strings.TrimPrefix(localVersion, "v")

	cmp, err := compareVersions(localVersionTrimmed, remoteVersion)
	if err != nil {
		// A "dev" build (or one without an injected version) can't be compared
		// as semver — treat it as always eligible to update rather than failing.
		cmp = -1
	}
	if cmp == 0 {
		fmt.Printf("Already up to date (%s)\n", localVersion)
		checkAndReInit(projectDir, localVersion)
		return
	}
	if cmp > 0 {
		fmt.Printf("Local version (%s) is newer than latest release (%s). Nothing to do.\n", localVersion, release.TagName)
		checkAndReInit(projectDir, localVersion)
		return
	}

	// Build the asset name for the current platform and architecture.
	assetName := fmt.Sprintf("agentic-cms_%s_%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		assetName += ".exe"
	}
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		fmt.Fprintf(os.Stderr, "No asset named %q found in release %s\n", assetName, release.TagName)
		os.Exit(1)
	}

	tmpFile := filepath.Join(os.TempDir(), "agentic-cms.new")
	if err := downloadBinary(downloadURL, tmpFile); err != nil {
		_ = os.Remove(tmpFile)
		fmt.Fprintf(os.Stderr, "Error downloading binary: %v\n", err)
		os.Exit(1)
	}

	if err := os.Chmod(tmpFile, 0o755); err != nil {
		_ = os.Remove(tmpFile)
		fmt.Fprintf(os.Stderr, "Error setting executable bit: %v\n", err)
		os.Exit(1)
	}

	// Resolve current binary path using os.Executable (cross-platform).
	currentBin, err := os.Executable()
	if err != nil {
		_ = os.Remove(tmpFile)
		fmt.Fprintf(os.Stderr, "Error detecting binary path: %v\n", err)
		os.Exit(1)
	}
	currentBin, err = filepath.EvalSymlinks(currentBin)
	if err != nil {
		_ = os.Remove(tmpFile)
		fmt.Fprintf(os.Stderr, "Error resolving binary path: %v\n", err)
		os.Exit(1)
	}

	if err := replaceBinary(tmpFile, currentBin); err != nil {
		_ = os.Remove(tmpFile)
		fmt.Fprintf(os.Stderr, "Error replacing binary: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated from %s to %s\n", localVersion, release.TagName)

	// Re-exec the freshly-downloaded binary to perform the reinit, rather than
	// reinitializing in this (now stale) process. go:embed content is compiled
	// into the binary at build time, so this still-running process only has
	// the OLD version's embedded scaffold in memory even though the file on
	// disk was just replaced — reinitializing here would silently skip
	// anything new in the release just downloaded.
	reinitWithBinary(currentBin, projectDir)
}

// resolveDefaultProjectDir finds the project root for `update` when invoked
// without an explicit target, so it works from any subdirectory of an
// installed project rather than only the current directory. It walks up to
// the nearest ancestor containing .agentic-cms/; a directory with no such
// ancestor (a fresh, not-yet-initialized project) falls back to startDir.
func resolveDefaultProjectDir(startDir string) string {
	absStart, err := filepath.Abs(startDir)
	if err != nil {
		return startDir
	}
	if root, ok := findAncestorWithEntry(absStart, ".agentic-cms"); ok {
		return root
	}
	return absStart
}

// findAncestorWithEntry walks up from startDir looking for a directory
// containing entryName, returning the first match or false if none is found
// before reaching the filesystem root.
func findAncestorWithEntry(startDir, entryName string) (string, bool) {
	dir := startDir
	for {
		if _, err := os.Stat(filepath.Join(dir, entryName)); err == nil {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

// fetchLatestRelease queries the GitHub API and returns release metadata.
func fetchLatestRelease() (*githubRelease, error) {
	const apiURL = "https://api.github.com/repos/ruifrvaz/agentic-cms/releases/latest"

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "agentic-cms/"+version)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("GitHub API rate limit reached. Try again in a few minutes, or download manually from https://github.com/ruifrvaz/agentic-cms/releases")
	}
	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("GitHub API request forbidden (status 403). Check your network or try again later")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("parsing GitHub API response: %w", err)
	}
	if release.TagName == "" {
		return nil, fmt.Errorf("no tag_name in GitHub API response")
	}
	return &release, nil
}

// compareVersions compares two semver strings (without leading "v").
// Returns -1 if a < b, 0 if equal, 1 if a > b, and a non-nil error if
// either version string cannot be parsed.
func compareVersions(a, b string) (int, error) {
	parse := func(v string) ([]int, error) {
		parts := strings.Split(v, ".")
		nums := make([]int, len(parts))
		for i, p := range parts {
			n, err := strconv.Atoi(p)
			if err != nil {
				return nil, fmt.Errorf("invalid version segment %q in %q", p, v)
			}
			nums[i] = n
		}
		return nums, nil
	}

	av, aerr := parse(a)
	bv, berr := parse(b)
	if aerr != nil {
		return 0, aerr
	}
	if berr != nil {
		return 0, berr
	}

	// Pad shorter slice
	for len(av) < len(bv) {
		av = append(av, 0)
	}
	for len(bv) < len(av) {
		bv = append(bv, 0)
	}

	for i := range av {
		if av[i] < bv[i] {
			return -1, nil
		}
		if av[i] > bv[i] {
			return 1, nil
		}
	}
	return 0, nil
}

// downloadBinary downloads url into destPath.
func downloadBinary(url, destPath string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "agentic-cms/"+version)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}
	return nil
}

// replaceBinary atomically replaces currentPath with tmpPath.
// Falls back to a same-filesystem copy + rename when /tmp is on a different
// filesystem than the binary (e.g., tmpfs vs ext4).
func replaceBinary(tmpPath, currentPath string) error {
	// Fast path: rename is atomic if src and dst are on the same filesystem.
	if err := os.Rename(tmpPath, currentPath); err == nil {
		return nil
	}

	// Fallback: write to a temp file in the same directory, then rename.
	sameDir := filepath.Dir(currentPath)
	tmpSameFS := filepath.Join(sameDir, fmt.Sprintf(".agentic-cms-%d.new", os.Getpid()))

	if err := copyFile(tmpPath, tmpSameFS); err != nil {
		_ = os.Remove(tmpSameFS)
		return fmt.Errorf("copy to same filesystem: %w", err)
	}
	if err := os.Chmod(tmpSameFS, 0o755); err != nil {
		_ = os.Remove(tmpSameFS)
		return err
	}
	if err := os.Rename(tmpSameFS, currentPath); err != nil {
		_ = os.Remove(tmpSameFS)
		return fmt.Errorf("rename on same filesystem: %w", err)
	}
	_ = os.Remove(tmpPath)
	return nil
}

// copyFile copies src to dst, creating dst if necessary.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("copying file: %w", err)
	}
	return nil
}

// checkAndReInit checks whether dir contains a .agentic-cms/ directory. If so
// it re-runs the init command to deploy updated skills, agents, and
// templates. Safe to run in-process only when this process's own embedded
// content is known to match what's on disk (i.e., no binary replacement just
// happened in this run) — see reinitWithBinary for the post-replacement case.
func checkAndReInit(dir, version string) {
	scaffoldPath := filepath.Join(dir, ".agentic-cms")
	if _, err := os.Stat(scaffoldPath); err != nil {
		// .agentic-cms/ not present — skip auto-init
		fmt.Println("Run `agentic-cms init` to install or update the project scaffolding")
		return
	}

	fmt.Println("Detected .agentic-cms/ — re-initializing project scaffolding...")
	if err := runInit(dir, version); err != nil {
		fmt.Fprintf(os.Stderr, "agentic-cms: %v\n", err)
		os.Exit(1)
	}
}

// reinitWithBinary checks whether dir contains a .agentic-cms/ directory. If
// so it re-invokes binaryPath (the just-downloaded binary on disk) as a
// subprocess to perform the reinit, instead of calling runInit() in this
// process. This process's own go:embed content was compiled in at build time
// and does not reflect the binary that now sits on disk after a self-update —
// only a fresh process image loaded from that file has the new release's
// skills, agents, and templates.
func reinitWithBinary(binaryPath, dir string) {
	scaffoldPath := filepath.Join(dir, ".agentic-cms")
	if _, err := os.Stat(scaffoldPath); err != nil {
		// .agentic-cms/ not present — skip auto-init
		fmt.Println("Run `agentic-cms init` to install the project scaffolding")
		return
	}

	fmt.Println("Detected .agentic-cms/ — re-initializing project scaffolding with the updated binary...")
	cmd := exec.Command(binaryPath, "init", dir)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error re-initializing project scaffolding: %v\n", err)
		os.Exit(1)
	}
}
