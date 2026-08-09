package testr

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDetectTestModel_LinksShipr(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "module example.com/p\n\ngo 1.22\n")
	sm := map[string]any{
		"schema_version": 1,
		"product_id":     "p",
		"proof_commands": []string{"go test ./...", "go build ./..."},
	}
	b, _ := json.MarshalIndent(sm, "", "  ")
	writeFile(t, filepath.Join(dir, ShiprModelRel), string(b)+"\n")

	m := DetectTestModel(dir, "test product")
	if m["role"] != "ai_config_and_memory" {
		t.Fatalf("role: %v", m["role"])
	}
	rel, ok := m["related_shipr"].(map[string]any)
	if !ok || rel["loaded"] != true {
		t.Fatalf("related_shipr: %v", m["related_shipr"])
	}
	if rel["product_id"] != "p" {
		t.Fatalf("shipr product_id: %v", rel["product_id"])
	}
	cmds := m["test_commands"].([]string)
	if len(cmds) == 0 || cmds[0] != "go test ./..." {
		t.Fatalf("test_commands: %v", m["test_commands"])
	}
}

func TestWriteLoadAttemptFrontier(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "module example.com/p\n\ngo 1.22\n")
	writeFile(t, filepath.Join(dir, ".gitignore"), "")
	m := DetectTestModel(dir, "")
	path, err := WriteTestModel(dir, m)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadTestModel(dir)
	if err != nil {
		t.Fatal(err)
	}
	if loaded["role"] != "ai_config_and_memory" {
		t.Fatalf("loaded role: %v", loaded["role"])
	}
	ap, attempt, err := RecordAttempt(dir, "full suite", "passed", "", []string{"go test ./..."})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(ap); err != nil {
		t.Fatal(err)
	}
	if attempt["status"] != "passed" {
		t.Fatalf("attempt: %v", attempt)
	}
	fr := Frontier(dir)
	if fr["status"] != "model_ready" {
		t.Fatalf("frontier: %v", fr)
	}
	if fr["role"] != "ai_config_and_memory" {
		t.Fatalf("frontier role: %v", fr["role"])
	}
}

func TestWriteDoesNotGitignoreAndCreatesSibling(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "module example.com/p\n\ngo 1.22\n")
	writeFile(t, filepath.Join(dir, ".gitignore"), "*.log\n.testr/\n.shipr/\n")
	m := DetectTestModel(dir, "")
	path, err := WriteTestModel(dir, m)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, ShiprModelRel)); err != nil {
		t.Fatalf("expected sibling shipr model: %v", err)
	}
	gi, _ := os.ReadFile(filepath.Join(dir, ".gitignore"))
	text := string(gi)
	if strings.Contains(text, ".shipr/") || strings.Contains(text, ".testr/") {
		t.Fatalf("gitignore still ignores configs:\n%s", text)
	}
}

func TestEnsureProductConfigsCreatesBoth(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "module example.com/p\n\ngo 1.22\n")
	if err := EnsureProductConfigs(dir); err != nil {
		t.Fatal(err)
	}
	if !exists(dir, ModelRelPath) || !exists(dir, ShiprModelRel) {
		t.Fatal("expected both models")
	}
}

func TestFrontierUsesLatestAttemptByTime(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "go.mod"), "module example.com/p\n\ngo 1.22\n")
	m := DetectTestModel(dir, "")
	if _, err := WriteTestModel(dir, m); err != nil {
		t.Fatal(err)
	}
	for _, st := range []string{"planned", "running", "passed"} {
		if _, _, err := RecordAttempt(dir, "stress-"+st, st, "", []string{"go test ./..."}); err != nil {
			t.Fatal(err)
		}
	}
	fr := Frontier(dir)
	if fr["latest_status"] != "passed" {
		t.Fatalf("want latest_status=passed got %v (latest=%v)", fr["latest_status"], fr["latest_attempt"])
	}
}
