package revisorschemas_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/ttab/newsdoc"
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
	expectSchemas := 7

	if len(files) != expectSchemas {
		t.Fatalf("expected there to be %d schema files, got %d",
			expectSchemas, len(files))
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

func TestValidDocuments(t *testing.T) {
	schemaFS := revisorschemas.Files()

	files, err := schemaFS.ReadDir(".")
	if err != nil {
		t.Fatalf("failed to read directory")
	}

	var constraints []revisor.ConstraintSet

	for _, f := range files {
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

		constraints = append(constraints, cs)
	}

	validator, err := revisor.NewValidator(constraints...)
	if err != nil {
		t.Fatalf("failed to create validator with schemas: %v", err)
	}

	dataDir := os.DirFS(filepath.Join("testdata", "valid-docs"))

	validDocs, err := fs.Glob(dataDir, "*.json")
	if err != nil {
		t.Fatalf("failed to glob for valid documents: %v", err)
	}

	for _, name := range validDocs {
		t.Run(name, func(t *testing.T) {
			data, err := fs.ReadFile(dataDir, name)
			if err != nil {
				t.Fatalf("failed to read document file: %v", err)
			}

			dec := json.NewDecoder(bytes.NewReader(data))

			dec.DisallowUnknownFields()

			var doc newsdoc.Document

			err = dec.Decode(&doc)
			if err != nil {
				t.Fatalf("failed to unmarshal document file: %v", err)
			}

			validationErrors, err := validator.ValidateDocument(
				context.Background(), &doc)
			if err != nil {
				t.Fatalf("validation failed: %v", err)
			}

			for _, e := range validationErrors {
				t.Error(e.String())
			}
		})
	}
}
