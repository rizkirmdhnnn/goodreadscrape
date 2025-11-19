package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadURLsFromFile(t *testing.T) {
	// Helper function to create a temp file with content
	createTempFile := func(t *testing.T, content string) string {
		t.Helper()
		tmpfile, err := os.CreateTemp("", "test_urls_*.txt")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := tmpfile.WriteString(content); err != nil {
			t.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}
		return tmpfile.Name()
	}

	tests := []struct {
		name        string
		fileContent string
		want        []string
		wantErr     bool
	}{
		{
			name: "Valid file with URLs",
			fileContent: `http://example.com/1
http://example.com/2
http://example.com/3`,
			want: []string{
				"http://example.com/1",
				"http://example.com/2",
				"http://example.com/3",
			},
			wantErr: false,
		},
		{
			name: "File with comments and empty lines",
			fileContent: `
# This is a comment
http://example.com/1

http://example.com/2
# Another comment
`,
			want: []string{
				"http://example.com/1",
				"http://example.com/2",
			},
			wantErr: false,
		},
		{
			name:        "Empty file",
			fileContent: "",
			want:        nil,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile := createTempFile(t, tt.fileContent)
			defer os.Remove(tmpfile)

			got, err := LoadURLsFromFile(tmpfile)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadURLsFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadURLsFromFile() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("File not found", func(t *testing.T) {
		_, err := LoadURLsFromFile("non_existent_file.txt")
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
	})
}
