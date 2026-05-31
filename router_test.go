package main

import (
	"testing"

	"github.com/DonovanDiamond/milter"
)

func TestRouterMilter_MailFrom(t *testing.T) {
	rejectedFrom := map[string]bool{"rejected@example.com": true}
	r := &RouterMilter{
		rejectedFrom: rejectedFrom,
	}

	// Test rejected address
	resp, err := r.MailFrom("rejected@example.com", nil)
	if err != nil {
		t.Errorf("MailFrom returned error: %v", err)
	}
	if resp != milter.RespReject {
		t.Errorf("Expected RespReject for rejected address, got %v", resp)
	}

	// Test allowed address
	resp, err = r.MailFrom("allowed@example.com", nil)
	if err != nil {
		t.Errorf("MailFrom returned error: %v", err)
	}
	if resp != milter.RespContinue {
		t.Errorf("Expected RespContinue for allowed address, got %v", resp)
	}
}

func TestRouterMilter_RcptTo_Multiple(t *testing.T) {
	rejectedTo := map[string]bool{"rejected@example.com": true}
	r := &RouterMilter{
		rejectedTo: rejectedTo,
	}

	// Test rejected address
	resp, _ := r.RcptTo("rejected@example.com", nil)
	if resp != milter.RespReject {
		t.Errorf("Expected RespReject, got %v", resp)
	}

	// Test rejected address again
	resp, _ = r.RcptTo("rejected@example.com", nil)
	if resp != milter.RespReject {
		t.Errorf("Expected RespReject, got %v", resp)
	}

	// Test allowed address
	resp, _ = r.RcptTo("allowed@example.com", nil)
	if resp != milter.RespContinue {
		t.Errorf("Expected RespContinue, got %v", resp)
	}

	// Test rejected address again
	resp, _ = r.RcptTo("rejected@example.com", nil)
	if resp != milter.RespReject {
		t.Errorf("Expected RespReject, got %v", resp)
	}

	// Verify both were handled (even if one was rejected, we might track them depending on logic)
	// In our case, we only append allowed ones to e.to
	if len(r.to) != 1 || r.to[0] != "allowed@example.com" {
		t.Errorf("Expected 1 allowed recipient, got %v", r.to)
	}
}

func TestRouterMilter_Reset(t *testing.T) {
	r := &RouterMilter{}
	r.MailFrom("sender@example.com", nil)
	r.RcptTo("rcpt@example.com", nil)

	if r.from != "sender@example.com" || len(r.to) != 1 {
		t.Fatalf("Setup failed")
	}

	// Test Abort
	r.Abort(nil)
	if r.from != "" || len(r.to) != 0 {
		t.Errorf("Abort did not reset state: from=%s, to=%v", r.from, r.to)
	}

	// Test MailFrom reset
	r.MailFrom("new@example.com", nil)
	r.RcptTo("new-rcpt@example.com", nil)
	r.MailFrom("another@example.com", nil)
	if r.from != "another@example.com" || len(r.to) != 0 {
		t.Errorf("MailFrom did not reset state: from=%s, to=%v", r.from, r.to)
	}
}
