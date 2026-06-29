package hfdownloader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindSourceFile(t *testing.T) {
	tmp := t.TempDir()

	// Create a layout like: tmp/owner/repo/file.txt
	modelDir := filepath.Join(tmp, "owner", "repo")
	if err := os.MkdirAll(modelDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(modelDir, "file.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		sourceDir string
		repo      string
		path      string
		wantFound bool
	}{
		{
			name:      "found exact path",
			sourceDir: tmp,
			repo:      "owner/repo",
			path:      "file.txt",
			wantFound: true,
		},
		{
			name:      "wrong repo",
			sourceDir: tmp,
			repo:      "other/repo",
			path:      "file.txt",
			wantFound: false,
		},
		{
			name:      "nested path missing",
			sourceDir: tmp,
			repo:      "owner/repo",
			path:      "subdir/config.json",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := findSourceFile(tt.sourceDir, tt.repo, tt.path)
			if found != tt.wantFound {
				t.Errorf("findSourceFile(%q, %q, %q) found=%v, want %v",
					tt.sourceDir, tt.repo, tt.path, found, tt.wantFound)
			}
			if found {
				wantSuffix := filepath.Join(tt.repo, tt.path)
				if !filepath.IsAbs(got) {
					t.Errorf("findSourceFile returned non-absolute path: %s", got)
				}
				if !strings.HasSuffix(got, wantSuffix) {
					t.Errorf("findSourceFile path ends with %q, got: %s", wantSuffix, got)
				}
			}
		})
	}
}

func TestLinkOrCopy(t *testing.T) {
	tmp := t.TempDir()

	// Create source file with known content
	src := filepath.Join(tmp, "src.txt")
	content := []byte("test content for link or copy")
	if err := os.WriteFile(src, content, 0o644); err != nil {
		t.Fatal(err)
	}

	// Test 1: Hard link on same filesystem (this platform supports it)
	dstHard := filepath.Join(tmp, "hardlink.txt")
	if err := linkOrCopy(src, dstHard); err != nil {
		t.Fatalf("linkOrCopy hard link: %v", err)
	}

	got, err := os.ReadFile(dstHard)
	if err != nil {
		t.Fatalf("failed to read link target: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("linkOrCopy content = %q, want %q", string(got), string(content))
	}

	// Test 2: Fallback copy (parent dir must exist — caller ensures this)
	dstCopy := filepath.Join(tmp, "sub", "dir", "copy.txt")
	if err := os.MkdirAll(filepath.Dir(dstCopy), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := linkOrCopy(src, dstCopy); err != nil {
		t.Fatalf("linkOrCopy fallback: %v", err)
	}

	got, err = os.ReadFile(dstCopy)
	if err != nil {
		t.Fatalf("failed to read copy target: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("copy content = %q, want %q", string(got), string(content))
	}

	// Test 3: Non-existent source should fail
	dstFail := filepath.Join(tmp, "fail.txt")
	if err := linkOrCopy(filepath.Join(tmp, "no-such-file"), dstFail); err == nil {
		t.Error("linkOrCopy with non-existent source should fail")
	}
}
