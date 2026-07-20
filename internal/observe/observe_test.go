package observe_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/meclaw/meclaw/internal/observe"
)

func TestJSONLTracer(t *testing.T) {
	dir := t.TempDir()
	tr, err := observe.NewJSONLTracer(dir)
	if err != nil {
		t.Fatal(err)
	}
	sp := tr.Start("handle", map[string]string{"channel": "stdio"})
	tr.End(sp, nil)
	data, err := os.ReadFile(filepath.Join(dir, "traces.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"name":"handle"`) {
		t.Fatalf("%s", data)
	}
}
