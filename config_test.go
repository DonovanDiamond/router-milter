package main

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	content := `
protocol: tcp
address: :1234
reject_from:
  - joe@example.com
reject_to:
  - alice@example.com
reject_to_regex:
  - .*@regex\.com
`
	tmpfile, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	expected := &Config{
		Protocol:   "tcp",
		Address:    ":1234",
		RejectFrom: []string{"joe@example.com"},
		RejectTo:   []string{"alice@example.com"},
		RejectToRegex: []string{".*@regex\\.com"},
	}

	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("Expected %+v, got %+v", expected, cfg)
	}
}
