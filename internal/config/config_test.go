package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	content := `
antigravity:
  roles:
    - architect
    - senior-engineer
claude:
  roles:
    - junior-engineer
`
	tmpDir, err := os.MkdirTemp("", "castra-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	yamlPath := filepath.Join(tmpDir, "castra.yaml")
	if err := os.WriteFile(yamlPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Parse(yamlPath)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expected := map[string]VendorConfig{
		"antigravity": {Roles: []string{"architect", "senior-engineer"}},
		"claude":      {Roles: []string{"junior-engineer"}},
	}

	if !reflect.DeepEqual(cfg.Vendors, expected) {
		t.Errorf("Expected %+v, got %+v", expected, cfg.Vendors)
	}
}

func TestGenerateTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "castra-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	err = GenerateTemplate(tmpDir)
	if err != nil {
		t.Fatalf("GenerateTemplate failed: %v", err)
	}

	yamlPath := filepath.Join(tmpDir, "castra.yaml")
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		t.Error("castra.yaml was not generated")
	}
}
