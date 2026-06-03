package session

import (
	"net/textproto"
	"testing"

	"github.com/DonovanDiamond/milter"
)

func TestSessionMatchesMilterInterface(t *testing.T) {
	session := NewSession("function execute(commands) { return 'okay'; }")
	// this verifies that all functions specified by the 'milter.Milter'
	// interface are set so that the app will not crash if one of these gets
	// called without being defind.
	_, _ = session.Connect("", "", 0, []byte{}, nil)
	_, _ = session.Helo("", nil)
	_, _ = session.MailFrom("", nil)
	_, _ = session.RcptTo("", nil)
	_, _ = session.Header("", "", nil)
	_, _ = session.Headers(textproto.MIMEHeader{}, nil)
	_, _ = session.BodyChunk([]byte{}, nil)
	_, _ = session.Body(nil)
	_ = session.Abort(nil)
}

func TestNewSession(t *testing.T) {
	script := "function execute(commands) { return 'okay'; }"
	session := NewSession(script)
	if session.script != script {
		t.Errorf("expected script to be set, got '%s'", session.script)
	}
	if session.commands.HELO != "" || session.commands.FROM != "" || len(session.commands.RCPT) != 0 {
		t.Errorf("expected commands to be empty, got %+v", session.commands)
	}
}

func TestSessionConnect(t *testing.T) {
	session := NewSession("function execute(commands) { return 'okay'; }")
	resp, err := session.Connect("client.example.com", "inet", 12345, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	if session.host != "client.example.com" {
		t.Errorf("expected host to be set, got '%s'", session.host)
	}
}

func TestSessionHelo(t *testing.T) {
	session := NewSession("function execute(commands) { return 'okay'; }")
	resp, err := session.Helo("example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	if session.commands.HELO != "example.com" {
		t.Errorf("expected HELO to be set, got '%s'", session.commands.HELO)
	}
}

func TestSessionMailFrom(t *testing.T) {
	session := NewSession("function execute(commands) { return 'okay'; }")
	resp, err := session.MailFrom("joe@example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	if session.commands.FROM != "joe@example.com" {
		t.Errorf("expected FROM to be set, got '%s'", session.commands.FROM)
	}
}

func TestSessionMailFromUppercase(t *testing.T) {
	session := NewSession("function execute(commands) { return 'okay'; }")
	resp, err := session.MailFrom("JOe@exAMPLE.COM", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	if session.commands.FROM != "joe@example.com" {
		t.Errorf("expected FROM to be set to lowercase, got '%s'", session.commands.FROM)
	}
}

func TestSessionMailFromMultiple(t *testing.T) {
	session := NewSession("function execute(commands) { return 'okay'; }")
	resp, err := session.MailFrom("joe@example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	if session.commands.FROM != "joe@example.com" {
		t.Errorf("expected FROM to be set, got '%s'", session.commands.FROM)
	}
	resp, err = session.RcptTo("joe@example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	if len(session.commands.RCPT) != 1 || session.commands.RCPT[0] != "joe@example.com" {
		t.Errorf("expected RCPT to contain 'joe@example.com', got %+v", session.commands.RCPT)
	}
	resp, err = session.MailFrom("mary@example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	if session.commands.FROM != "mary@example.com" {
		t.Errorf("expected FROM to be set, got '%s'", session.commands.FROM)
	}
	if len(session.commands.RCPT) != 0 {
		t.Errorf("expected RCPT to be empty, got %+v", session.commands.RCPT)
	}
}

func TestSessionRcptTo(t *testing.T) {
	session := NewSession("function execute(commands) { return 'okay'; }")
	resp, err := session.RcptTo("joe@example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	if len(session.commands.RCPT) != 1 || session.commands.RCPT[0] != "joe@example.com" {
		t.Errorf("expected RCPT to contain 'joe@example.com', got %+v", session.commands.RCPT)
	}
}

func TestSessionRcptToMultiple(t *testing.T) {
	session := NewSession("function execute(commands) { return 'okay'; }")
	resp, err := session.RcptTo("joe@example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	resp, err = session.RcptTo("mary@example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	if len(session.commands.RCPT) != 2 || session.commands.RCPT[0] != "joe@example.com" || session.commands.RCPT[1] != "mary@example.com" {
		t.Errorf("expected RCPT to contain 'joe@example.com' and 'mary@example.com', got %+v", session.commands.RCPT)
	}
}

func TestSessionRcptToDuplicate(t *testing.T) {
	session := NewSession("function execute(commands) { return 'okay'; }")
	resp, err := session.RcptTo("joe@example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	resp, err = session.RcptTo("joe@example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	if len(session.commands.RCPT) != 2 || session.commands.RCPT[0] != "joe@example.com" || session.commands.RCPT[1] != "joe@example.com" {
		t.Errorf("expected RCPT to contain 'joe@example.com' twice, got %+v", session.commands.RCPT)
	}
}

func TestSessionRcptToUppercase(t *testing.T) {
	session := NewSession("function execute(commands) { return 'okay'; }")
	resp, err := session.RcptTo("JOE@examPLE.cOm", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	if len(session.commands.RCPT) != 1 || session.commands.RCPT[0] != "joe@example.com" {
		t.Errorf("expected RCPT to contain 'joe@example.com', got %+v", session.commands.RCPT)
	}
}

func TestSessionBody(t *testing.T) {
	session := NewSession("function execute(commands) { return 'okay'; }")
	resp, err := session.Body(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
}

func TestFullSessionOkay(t *testing.T) {
	session := NewSession("function execute({ HELO, FROM, RCPT }) { return FROM == 'joe@example.com' ? 'okay' : 'reject' }")
	resp, err := session.MailFrom("joe@example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	resp, err = session.RcptTo("mary@example.net", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	resp, err = session.Body(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
}

func TestFullSessionReject(t *testing.T) {
	session := NewSession("function execute({ HELO, FROM, RCPT }) { return FROM == 'joe@example.com' ? 'okay' : 'reject' }")
	resp, err := session.MailFrom("joe2@example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	resp, err = session.RcptTo("mary@example.net", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("expected response to be RespContinue, got %v", resp)
	}
	resp, err = session.Body(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != milter.RespReject {
		t.Errorf("expected response to be RespReject, got %v", resp)
	}
}
