package revisorschemas_test

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	testDebug     bool
	readDebugOnce sync.Once
)

func debug() bool {
	readDebugOnce.Do(func() {
		testDebug = os.Getenv("TEST_DEBUG") == "true"
	})

	return testDebug
}

type testingT interface {
	Helper()
	Fatalf(format string, args ...any)
	Logf(format string, args ...any)
}

func must(t testingT, err error, format string, a ...any) {
	t.Helper()

	if err != nil {
		t.Fatalf("failed: %s: %v", fmt.Sprintf(format, a...), err)
	}

	if debug() {
		t.Logf("success: "+format, a...)
	}
}

// testAgainstGolden compares a result against the contents of the file at the
// goldenPath. Run with regenerate set to true to create or update the file.
func testAgainstGolden[T any](
	t *testing.T,
	regenerate bool,
	got T,
	goldenPath string,
) {
	t.Helper()

	if regenerate {
		data, err := json.MarshalIndent(got, "", "  ")
		must(t, err, "marshal for storage in %q", goldenPath)

		// End all files with a newline
		data = append(data, '\n')

		err = os.WriteFile(goldenPath, data, 0o600)
		must(t, err, "write golden file %q", goldenPath)
	}

	wantData, err := os.ReadFile(goldenPath)
	must(t, err, "read from golden file %q", goldenPath)

	var wantValue T

	err = json.Unmarshal(wantData, &wantValue)
	must(t, err, "unmarshal data from golden file %q", goldenPath)

	equalDiffWithOptions(t, wantValue, got, nil,
		"must match golden file %q", goldenPath)
}

// EqualMessage runs a cmp.Diff with protobuf-specific options.
func equalDiffWithOptions[T any](
	t testingT,
	want T, got T,
	opts cmp.Options,
	format string, a ...any,
) {
	t.Helper()

	diff := cmp.Diff(want, got, opts...)
	if diff != "" {
		msg := fmt.Sprintf(format, a...)
		t.Fatalf("%s: mismatch (-want +got):\n%s", msg, diff)
	}

	if debug() {
		t.Logf("success: "+format, a...)
	}
}
