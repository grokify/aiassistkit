package requirements

import (
	"runtime"
	"testing"
)

func TestCheckerIsInstalled(t *testing.T) {
	checker := NewChecker()

	// "go" should be installed in most dev environments
	if !checker.IsInstalled("go") {
		t.Skip("go not installed, skipping test")
	}

	// Non-existent tool should not be installed
	if checker.IsInstalled("this-tool-definitely-does-not-exist-xyz123") {
		t.Error("non-existent tool reported as installed")
	}
}

func TestCheckerCheck(t *testing.T) {
	checker := NewChecker()

	// Check a mix of installed and missing tools
	result := checker.Check([]string{"go", "this-tool-does-not-exist-abc"})

	if len(result.Satisfied) == 0 {
		t.Skip("go not installed, skipping test")
	}

	// go should be satisfied
	found := false
	for _, name := range result.Satisfied {
		if name == "go" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'go' in satisfied list")
	}

	// Unknown tool should be in unknown list
	found = false
	for _, name := range result.Unknown {
		if name == "this-tool-does-not-exist-abc" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected unknown tool in unknown list")
	}
}

func TestRegistryGet(t *testing.T) {
	reg := DefaultRegistry

	// Known tool
	req := reg.Get("go")
	if req == nil {
		t.Fatal("expected to find 'go' in registry")
	}
	if req.Name != "go" {
		t.Errorf("Name = %q, want %q", req.Name, "go")
	}
	if req.Purpose == "" {
		t.Error("Purpose should not be empty")
	}

	// Unknown tool
	if reg.Get("unknown-tool-xyz") != nil {
		t.Error("expected nil for unknown tool")
	}
}

func TestRegistryMerge(t *testing.T) {
	base := Registry{
		"tool1": {Name: "tool1", Purpose: "base tool"},
		"tool2": {Name: "tool2", Purpose: "base tool 2"},
	}

	overlay := Registry{
		"tool2": {Name: "tool2", Purpose: "overlay tool 2"},
		"tool3": {Name: "tool3", Purpose: "overlay tool 3"},
	}

	merged := base.Merge(overlay)

	if len(merged) != 3 {
		t.Errorf("merged length = %d, want 3", len(merged))
	}

	// tool1 from base
	if merged["tool1"].Purpose != "base tool" {
		t.Errorf("tool1 = %q, want %q", merged["tool1"].Purpose, "base tool")
	}

	// tool2 overridden by overlay
	if merged["tool2"].Purpose != "overlay tool 2" {
		t.Errorf("tool2 = %q, want %q", merged["tool2"].Purpose, "overlay tool 2")
	}

	// tool3 from overlay
	if merged["tool3"].Purpose != "overlay tool 3" {
		t.Errorf("tool3 = %q, want %q", merged["tool3"].Purpose, "overlay tool 3")
	}
}

func TestFindAvailableMethods(t *testing.T) {
	checker := NewChecker()

	// Get go requirement
	req := checker.Registry.Get("golangci-lint")
	if req == nil {
		t.Fatal("expected golangci-lint in registry")
	}

	methods := checker.findAvailableMethods(*req)

	// Should have at least one method available (go install if go is installed)
	if checker.IsInstalled("go") && len(methods) == 0 {
		t.Error("expected at least one install method for golangci-lint")
	}

	// Verify platform filtering works
	for _, m := range methods {
		if len(m.Platforms) > 0 {
			found := false
			for _, p := range m.Platforms {
				if p == runtime.GOOS {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("method %s has platform restriction but was included for %s", m.Name, runtime.GOOS)
			}
		}
	}
}

func TestCheckResultAllSatisfied(t *testing.T) {
	// All satisfied
	r1 := CheckResult{Satisfied: []string{"go", "git"}}
	if !r1.AllSatisfied() {
		t.Error("expected AllSatisfied to be true")
	}

	// Has missing
	r2 := CheckResult{
		Satisfied: []string{"go"},
		Missing:   []MissingRequirement{{Requirement: Requirement{Name: "foo"}}},
	}
	if r2.AllSatisfied() {
		t.Error("expected AllSatisfied to be false with missing")
	}

	// Has unknown
	r3 := CheckResult{
		Satisfied: []string{"go"},
		Unknown:   []string{"bar"},
	}
	if r3.AllSatisfied() {
		t.Error("expected AllSatisfied to be false with unknown")
	}
}

func TestGetInstallCommand(t *testing.T) {
	checker := NewChecker()

	// Known tool with available method
	if checker.IsInstalled("go") {
		cmd := checker.GetInstallCommand("golangci-lint")
		if cmd == "" {
			t.Error("expected install command for golangci-lint")
		}
	}

	// Unknown tool
	cmd := checker.GetInstallCommand("unknown-tool-xyz")
	if cmd != "" {
		t.Error("expected empty command for unknown tool")
	}
}
