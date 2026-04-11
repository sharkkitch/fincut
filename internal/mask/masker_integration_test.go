package mask_test

import (
	"strings"
	"testing"

	"github.com/user/fincut/internal/mask"
)

func TestMasker_LargeInput_AllSensitiveLinesRedacted(t *testing.T) {
	m, err := mask.New(mask.Options{
		Patterns: []string{`api_key=[A-Za-z0-9]+`},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := make([]string, 500)
	for i := range lines {
		if i%2 == 0 {
			lines[i] = "request api_key=supersecret123 received"
		} else {
			lines[i] = "heartbeat ok"
		}
	}

	out := m.Apply(lines)
	for i, line := range out {
		if i%2 == 0 && strings.Contains(line, "supersecret123") {
			t.Errorf("line %d: api_key not redacted: %q", i, line)
		}
		if i%2 != 0 && line != "heartbeat ok" {
			t.Errorf("line %d: unexpected modification: %q", i, line)
		}
	}

	count := m.CountRedacted(lines)
	if count != 250 {
		t.Errorf("expected 250 redacted lines, got %d", count)
	}
}

func TestMasker_ChainedPatterns_AllApplied(t *testing.T) {
	m, err := mask.New(mask.Options{
		Patterns: []string{
			`password=\S+`,
			`email=[^\s]+@[^\s]+`,
		},
		Replacement: "<hidden>",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := "user email=bob@example.com password=hunter2"
	out := m.Apply([]string{line})
	if strings.Contains(out[0], "bob@example.com") {
		t.Errorf("email not redacted: %q", out[0])
	}
	if strings.Contains(out[0], "hunter2") {
		t.Errorf("password not redacted: %q", out[0])
	}
	if !strings.Contains(out[0], "<hidden>") {
		t.Errorf("expected <hidden> in output: %q", out[0])
	}
}
