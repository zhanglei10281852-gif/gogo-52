package rail

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestReadAuditRejectsMultipleRecordsOnOneLine(t *testing.T) {
	store := NewStore(t.TempDir())
	now := time.Date(2025, 3, 1, 12, 0, 0, 0, time.UTC)
	if _, err := store.AppendAudit("plan", "tester", "demo@1.0.0", map[string]string{"step": "one"}, now); err != nil {
		t.Fatal(err)
	}
	if _, err := store.AppendAudit("apply", "tester", "demo@1.0.0", map[string]string{"step": "two"}, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(store.AuditPath())
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected two audit lines, got %d", len(lines))
	}
	joined := lines[0] + " " + lines[1] + "\n"
	if err := os.WriteFile(store.AuditPath(), []byte(joined), 0o600); err != nil {
		t.Fatal(err)
	}
	verification, err := store.VerifyAudit()
	if err != nil {
		t.Fatal(err)
	}
	if verification.Valid {
		t.Fatal("line-based verification unexpectedly accepted the collapsed audit file")
	}
	if _, err := store.ReadAudit(0); err == nil {
		t.Fatal("ReadAudit accepted an audit file whose line framing verification rejects")
	}
}
