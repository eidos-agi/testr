package testr

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	ModelRelPath   = ".testr/product-test-model.json"
	AttemptsRelDir = ".testr/test-attempts"
	IgnoreEntry    = ".testr/"
	ShiprModelRel  = ".shipr/product-release-model.json"
	// Version is the Go CLI / config schema generation version.
	Version = "0.2.0"
)

// TestModel is the durable per-product test config AI agents read.
// testr does not run tests; it stores how-to-prove + attempt memory.
type TestModel map[string]any

func exists(root string, parts ...string) bool {
	_, err := os.Stat(filepath.Join(append([]string{root}, parts...)...))
	return err == nil
}

func ensureIgnored(root string) {
	if !exists(root, ".git") && !exists(root, ".gitignore") {
		return
	}
	gi := filepath.Join(root, ".gitignore")
	b, _ := os.ReadFile(gi)
	text := string(b)
	for _, line := range strings.Split(text, "\n") {
		s := strings.TrimSpace(line)
		if s == ".testr" || s == ".testr/" || s == ".testr/*" {
			return
		}
	}
	if text != "" && !strings.HasSuffix(text, "\n") {
		text += "\n"
	}
	_ = os.WriteFile(gi, []byte(text+IgnoreEntry+"\n"), 0o644)
}

func uniqueKeepOrder(ss []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, s := range ss {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func stringSlice(v any) []string {
	switch t := v.(type) {
	case []string:
		return t
	case []any:
		var out []string
		for _, x := range t {
			if s, ok := x.(string); ok && s != "" {
				out = append(out, s)
			}
		}
		return out
	default:
		return nil
	}
}

// loadShiprLink reads sibling shipr model for cross-reference.
func loadShiprLink(root string) map[string]any {
	out := map[string]any{
		"operator":   "shipr",
		"model_path": ShiprModelRel,
		"loaded":     false,
		"note":       "shipr should use these test_commands as proof_commands when shipping",
	}
	path := filepath.Join(root, ShiprModelRel)
	b, err := os.ReadFile(path)
	if err != nil {
		return out
	}
	var sm map[string]any
	if json.Unmarshal(b, &sm) != nil {
		return out
	}
	out["loaded"] = true
	out["abs_path"] = path
	if pid, ok := sm["product_id"].(string); ok {
		out["product_id"] = pid
	}
	if pc := stringSlice(sm["proof_commands"]); len(pc) > 0 {
		out["shipr_proof_commands"] = pc
	}
	return out
}

// DetectTestModel builds (does not write) the product test config.
// Cross-links .shipr/product-release-model.json when present.
func DetectTestModel(project, description string) TestModel {
	root, _ := filepath.Abs(project)
	product := filepath.Base(root)
	var suites, commands []string

	// Prefer go when go.mod is present (canonical CLI language for shipr/testr).
	if exists(root, "go.mod") {
		suites = append(suites, "go")
		commands = append(commands, "go test ./...", "go build ./...")
	}
	if exists(root, "pyproject.toml") || exists(root, "pytest.ini") || exists(root, "tests") {
		// Only add pytest if there is an actual tests dir or pyproject without being go-only
		if exists(root, "tests") || exists(root, "pytest.ini") || (exists(root, "pyproject.toml") && !exists(root, "go.mod")) {
			suites = append(suites, "pytest")
			commands = append(commands, "python -m pytest -q")
		} else if exists(root, "pyproject.toml") && exists(root, "src") {
			// legacy python package still present alongside go
			suites = append(suites, "pytest")
			commands = append(commands, "python -m pytest -q")
		}
	}
	if exists(root, "package.json") {
		suites = append(suites, "node")
		commands = append(commands, "npm test")
	}
	if exists(root, "docs", "emf") {
		suites = append(suites, "emf")
		commands = append(commands, "python3 -m emf.validate docs/emf/")
	}
	if exists(root, "Makefile") {
		if b, err := os.ReadFile(filepath.Join(root, "Makefile")); err == nil {
			if regexp.MustCompile(`(?m)^test:`).Match(b) {
				suites = append(suites, "make-test")
				commands = append(commands, "make test")
			}
		}
	}
	if len(commands) == 0 {
		commands = append(commands, "define product-specific test command before claiming green")
	}
	sort.Strings(suites)
	return TestModel{
		"schema_version": 1,
		"role":           "ai_config_and_memory",
		"purpose":        "Tell AI agents how this product is proven. Store repeatable test config and attempt ledgers. Does not execute tests.",
		"product_id":     product,
		"project_root":   root,
		"description":    description,
		"test_suites":    suites,
		"test_commands":  uniqueKeepOrder(commands),
		"evidence_paths": []string{"test output / junit / coverage when configured"},
		"methods_source": "test-forge (retired) → testr",
		"related_shipr":  loadShiprLink(root),
		"related_operators": []string{"shipr"},
		"learning_questions": []string{
			"What failed that the ship path assumed green?",
			"Which test should become a shipr proof_command?",
			"What is flaky vs broken?",
		},
		"memory_paths": map[string]string{
			"model":        ModelRelPath,
			"attempts_dir": AttemptsRelDir,
		},
		"updated_at": time.Now().UTC().Format(time.RFC3339Nano),
	}
}

func writeJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

func WriteTestModel(project string, model TestModel) (string, error) {
	root, _ := filepath.Abs(project)
	ensureIgnored(root)
	path := filepath.Join(root, ModelRelPath)
	return path, writeJSON(path, model)
}

func LoadTestModel(project string) (TestModel, error) {
	root, _ := filepath.Abs(project)
	b, err := os.ReadFile(filepath.Join(root, ModelRelPath))
	if err != nil {
		return nil, err
	}
	var m TestModel
	return m, json.Unmarshal(b, &m)
}

func slug(text string) string {
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s := strings.Trim(re.ReplaceAllString(strings.ToLower(text), "-"), "-")
	if s == "" {
		return "test"
	}
	if len(s) > 60 {
		s = s[:60]
	}
	return s
}

// RecordAttempt appends a test-attempt ledger entry. It does not execute tests.
func RecordAttempt(project, goal, status, notes string, proofs []string) (string, map[string]any, error) {
	root, _ := filepath.Abs(project)
	ensureIgnored(root)
	model, err := LoadTestModel(root)
	if err != nil {
		model = DetectTestModel(root, "")
	}
	if proofs == nil {
		proofs = []string{}
	}
	ts := time.Now().UTC().Format("20060102T150405Z")
	path := filepath.Join(root, AttemptsRelDir, ts+"-"+slug(goal)+".json")
	attempt := map[string]any{
		"schema_version":         1,
		"product_id":             model["product_id"],
		"goal":                   goal,
		"status":                 status,
		"notes":                  notes,
		"proofs":                 proofs,
		"test_commands_snapshot": model["test_commands"],
		"created_at":             time.Now().UTC().Format(time.RFC3339Nano),
	}
	return path, attempt, writeJSON(path, attempt)
}

func Frontier(project string) map[string]any {
	root, _ := filepath.Abs(project)
	model, err := LoadTestModel(root)
	if err != nil {
		model = DetectTestModel(root, "")
	}
	attemptsDir := filepath.Join(root, AttemptsRelDir)
	var files []string
	if entries, e := os.ReadDir(attemptsDir); e == nil {
		for _, ent := range entries {
			if !ent.IsDir() && strings.HasSuffix(ent.Name(), ".json") {
				files = append(files, filepath.Join(attemptsDir, ent.Name()))
			}
		}
		sort.Strings(files)
	}
	var latest map[string]any
	if len(files) > 0 {
		b, _ := os.ReadFile(files[len(files)-1])
		_ = json.Unmarshal(b, &latest)
	}
	status := "model_missing"
	if err == nil {
		status = "model_ready"
	}
	next := []string{
		"AI: run test_commands from the model (testr does not execute them)",
		"record result with `testr attempt`",
		"keep shipr proof_commands aligned with test_commands",
	}
	if err != nil {
		next = append([]string{"run `testr model --write` to materialize test config for AI"}, next...)
	}
	out := map[string]any{
		"product_id":     model["product_id"],
		"role":           "ai_config_and_memory",
		"model_path":     filepath.Join(root, ModelRelPath),
		"test_commands":  model["test_commands"],
		"latest_status":  nil,
		"latest_attempt": nil,
		"latest":         latest,
		"status":         status,
		"next_actions":   next,
		"related_shipr":  loadShiprLink(root),
	}
	if latest != nil {
		out["latest_status"] = latest["status"]
		if len(files) > 0 {
			out["latest_attempt"] = files[len(files)-1]
		}
	}
	return out
}
