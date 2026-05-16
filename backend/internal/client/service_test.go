package client

import "testing"

func TestNormalizePhone(t *testing.T) {
	t.Parallel()

	got := normalizePhone("(92) 99999-1234")
	if got != "92999991234" {
		t.Fatalf("normalizePhone() = %q, want %q", got, "92999991234")
	}
}

func TestIsValidClientStatus(t *testing.T) {
	t.Parallel()

	valid := []string{"active", "inactive", "blocked"}
	for _, status := range valid {
		if !isValidClientStatus(status) {
			t.Fatalf("expected status %q to be valid", status)
		}
	}

	invalid := []string{"", "paused", "deleted", "ACTIVE"}
	for _, status := range invalid {
		if isValidClientStatus(status) {
			t.Fatalf("expected status %q to be invalid", status)
		}
	}
}
