package revisorschemas_test

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/ttab/newsdoc"
	"github.com/ttab/revisor"
	"github.com/ttab/revisorschemas"
)

func TestValidateSampleDocuments(t *testing.T) {
	sFS := revisorschemas.Files()

	core := decodeConstraintSetsFS(t, sFS,
		"core.json", "core-planning.json",
		"tt.json", "tt-planning.json",
	)

	validator, err := revisor.NewValidator(core...)
	must(t, err, "create validator")

	samplesDir := os.DirFS("testdata/samples")

	sampleDocs, err := fs.Glob(samplesDir, "*.json")
	must(t, err, "glob for sample documents")

	for _, name := range sampleDocs {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			data, err := fs.ReadFile(samplesDir, name)
			must(t, err, "read sample document")

			var doc newsdoc.Document

			err = decodeBytes(t, data, &doc)
			must(t, err, "decode document")

			res := validator.ValidateDocument(ctx, &doc)
			for _, e := range res {
				t.Error(e.String())
			}
		})
	}
}

func decodeConstraintSetsFS(
	t *testing.T, sFS embed.FS, names ...string,
) []revisor.ConstraintSet {
	t.Helper()

	var constraints []revisor.ConstraintSet

	for _, n := range names {
		data, err := sFS.ReadFile(n)
		must(t, err, "load constraints from %q", n)

		var c revisor.ConstraintSet

		err = decodeBytes(t, data, &c)
		must(t, err, "load constraints from %q", n)

		constraints = append(constraints, c)
	}

	return constraints
}

func decodeBytes(t *testing.T, data []byte, o interface{}) error {
	t.Helper()

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	return dec.Decode(o)
}

func must(t *testing.T, err error, msg string, a ...any) {
	t.Helper()

	if err != nil {
		t.Fatalf("%s: %v", fmt.Sprintf(msg, a...), err)
	}
}
