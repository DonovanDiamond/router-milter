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

func TestScriptSHA256(t *testing.T) {
	script := `
	function execute(commands) {
		let hash = sha256('test');
		if (hash != '9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08') {
			throw new Error("hash did not match: " + hash);
		}
		return 'okay';
	}
	`
	result, err := executeScript(script, Commands{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "okay" {
		t.Errorf("expected 'okay', got '%s'", result)
	}
}
