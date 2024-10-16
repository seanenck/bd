package shell_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/seanenck/blap/internal/shell"
)

func TestGenerationCompletions(t *testing.T) {
	if err := shell.GenerateCompletions(nil); err != nil {
		t.Errorf("invalid error: %v", err)
	}
	os.Clearenv()
	var buf bytes.Buffer
	t.Setenv("SHELL", "x/a")
	if err := shell.GenerateCompletions(&buf); err == nil || err.Error() != "unable to generate completions for a" {
		t.Errorf("invalid error: %v", err)
	}
	buf = bytes.Buffer{}
	t.Setenv("SHELL", "bash")
	if err := shell.GenerateCompletions(&buf); err != nil {
		t.Errorf("invalid error: %v", err)
	}
	b := buf.String()
	if !strings.Contains(b, "local ") {
		t.Errorf("invalid buffer: %s", b)
	}
}
