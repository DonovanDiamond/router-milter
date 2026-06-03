package session

import (
	"testing"
)

func TestExecuteNoopScript(t *testing.T) {
	script := `
	function execute(commands) {
		return "okay";
	}`
	commands := Commands{}
	result, err := executeScript(script, commands)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "okay" {
		t.Errorf("expected 'okay', got '%s'", result)
	}
}

func TestExecuteRejectScript(t *testing.T) {
	script := `
	function execute({ HELO, FROM, RCPT }) {
		if (FROM === "joe@example.com") {
			return "reject";
		}
		return "okay";
	}`
	commands := Commands{
		HELO: "example.com",
		FROM: "joe@example.com",
		RCPT: []string{""},
	}
	result, err := executeScript(script, commands)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "reject" {
		t.Errorf("expected 'reject', got '%s'", result)
	}
}

func TestExecuteTempFailScript(t *testing.T) {
	script := `
	function execute(commands) {
		throw new Error("temporary failure");
	}`
	commands := Commands{}
	result, err := executeScript(script, commands)
	if err == nil {
		t.Fatalf("expected error, got result: %s", result)
	}
}
