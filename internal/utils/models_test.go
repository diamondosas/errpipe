package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// expectedModelsPath returns the path where models.json should be written.
func expectedModelsPath(t *testing.T) string {
	t.Helper()
	configDir, err := os.UserConfigDir()
	if err != nil {
		t.Fatalf("cannot determine UserConfigDir: %v", err)
	}
	return filepath.Join(configDir, "errpipe", "models.json")
}

// TestFetchModels_WritesToAppData verifies that FetchModels:
//  1. Reaches the remote API and retrieves at least one non-empty model string.
//  2. Writes models.json to the expected AppData location.
//  3. The written file is valid JSON and contains the same models.
func TestFetchModels_WritesToAppData(t *testing.T) {
	wantPath := expectedModelsPath(t)

	// Remove any leftover cache so the test is deterministic.
	_ = os.Remove(wantPath)

	// Run synchronously (init() already called it in a goroutine; we call it
	// directly here so the test does not race against the background fetch).
	FetchModels()

	// ── 1. File must exist ────────────────────────────────────────────────────
	info, err := os.Stat(wantPath)
	if err != nil {
		t.Fatalf(
			"FAIL: models.json was not written to AppData\n"+
				"  expected path : %s\n"+
				"  stat error    : %v\n"+
				"  hint          : check FetchModels() in internal/utils/models.go",
			wantPath, err,
		)
	}
	if info.Size() == 0 {
		t.Fatalf(
			"FAIL: models.json exists but is empty\n"+
				"  path : %s\n"+
				"  hint : the HTTP response may have been empty or JSON marshal failed (models.go:69-71)",
			wantPath,
		)
	}
	fmt.Printf("PASS: models.json written to: %s (%d bytes)\n", wantPath, info.Size())

	// ── 2. File must contain valid JSON ───────────────────────────────────────
	raw, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("FAIL: cannot read models.json at %s: %v", wantPath, err)
	}

	var models []string
	if err := json.Unmarshal(raw, &models); err != nil {
		t.Fatalf(
			"FAIL: models.json is not valid JSON\n"+
				"  path    : %s\n"+
				"  content : %s\n"+
				"  error   : %v\n"+
				"  hint    : check json.Marshal call in FetchModels() (models.go:69)",
			wantPath, string(raw), err,
		)
	}

	// ── 3. Must have at least one non-empty model ─────────────────────────────
	nonEmpty := 0
	for _, m := range models {
		if m != "" {
			nonEmpty++
		}
	}
	if nonEmpty == 0 {
		t.Fatalf(
			"FAIL: models slice has no non-empty entries\n"+
				"  path    : %s\n"+
				"  models  : %v\n"+
				"  hint    : API may have returned empty strings or the Response struct\n"+
				"            fields (One/Two/Three) may not match the JSON keys (models.go:14-18)",
			wantPath, models,
		)
	}
	fmt.Printf("PASS: models fetched successfully: %v\n", models)
}

// TestGetModels_ReadsFromAppData verifies that GetModels reads whatever is in
// the AppData cache file and returns it faithfully.
func TestGetModels_ReadsFromAppData(t *testing.T) {
	wantPath := expectedModelsPath(t)

	// Write a known payload so we are not relying on network availability.
	want := []string{"model-alpha", "model-beta", "model-gamma"}
	data, _ := json.Marshal(want)

	if err := os.MkdirAll(filepath.Dir(wantPath), 0755); err != nil {
		t.Fatalf("cannot create directory for cache: %v", err)
	}
	if err := os.WriteFile(wantPath, data, 0644); err != nil {
		t.Fatalf("cannot write test cache at %s: %v", wantPath, err)
	}

	// Point the package-level variable to the file we just created.
	ModelsPath = wantPath

	got := GetModels()

	if len(got) == 0 {
		t.Fatalf(
			"FAIL: GetModels() returned nil/empty\n"+
				"  cache path : %s\n"+
				"  hint       : check os.ReadFile / json.Unmarshal in GetModels() (models.go:76-80)",
			wantPath,
		)
	}

	for i, g := range got {
		if i >= len(want) || g != want[i] {
			t.Errorf(
				"FAIL: GetModels()[%d] = %q, want %q\n"+
					"  cache path : %s\n"+
					"  hint       : value mismatch — check deserialization in GetModels() (models.go:78)",
				i, g, want[i], wantPath,
			)
		}
	}

	fmt.Printf("PASS: GetModels() returned: %v\n", got)
}

// TestFetchModels_UpdatesPackageVar verifies that after FetchModels() the
// package-level ModelsPath variable is non-empty and points at the expected
// location, so that GetModels() can actually find the file.
func TestFetchModels_UpdatesPackageVar(t *testing.T) {
	// Reset so the background init goroutine doesn't interfere.
	ModelsPath = ""

	// Give the init goroutine a short head-start, then call synchronously.
	time.Sleep(50 * time.Millisecond)
	FetchModels()

	if ModelsPath == "" {
		t.Fatalf(
			"FAIL: ModelsPath is still empty after FetchModels()\n"+
				"  hint: models.go line 65 uses ':=' (short variable declaration) which creates\n"+
				"        a LOCAL variable that shadows the package-level ModelsPath —\n"+
				"        change it to '=' (assignment) so the package var is actually set.\n"+
				"  file: internal/utils/models.go:65",
		)
	}

	wantSuffix := filepath.Join("errpipe", "models.json")
	if !filepath.IsAbs(ModelsPath) {
		t.Errorf(
			"FAIL: ModelsPath is not an absolute path\n"+
				"  got  : %s\n"+
				"  hint : check os.UserConfigDir() usage in FetchModels() (models.go:60-65)",
			ModelsPath,
		)
	}
	if base := filepath.Base(ModelsPath); base != "models.json" {
		t.Errorf(
			"FAIL: ModelsPath does not end in 'models.json'\n"+
				"  got  : %s\n"+
				"  want suffix: %s",
			ModelsPath, wantSuffix,
		)
	}

	fmt.Printf("PASS: ModelsPath correctly set to: %s\n", ModelsPath)
}
