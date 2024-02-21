package revisorschemas_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ttab/revisor"
	"github.com/ttab/revisorschemas"
)

func TestFS_ValidSchemas(t *testing.T) {
	schemaFS := revisorschemas.Files()

	files, err := schemaFS.ReadDir(".")
	if err != nil {
		t.Fatalf("failed to read directory")
	}

	// Sanity check so that we don't sail through on any issues that cause
	// embedding to fail.
	if len(files) != 4 {
		t.Fatal("expected there to be 4 schema files")
	}

	for _, f := range files {
		t.Run(f.Name(), func(t *testing.T) {
			data, err := schemaFS.ReadFile(f.Name())
			if err != nil {
				t.Fatalf("failed to read schema file: %v", err)
			}

			dec := json.NewDecoder(bytes.NewReader(data))

			dec.DisallowUnknownFields()

			var cs revisor.ConstraintSet

			err = dec.Decode(&cs)
			if err != nil {
				t.Fatalf("failed to unmarshal schema file: %v", err)
			}
		})
	}
}
