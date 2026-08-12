// Copyright 2025
// SPDX-License-Identifier: Apache-2.0

package hfdownloader

import (
	"os"
	"path/filepath"
	"strings"
)

// linkOrCopy attempts to create a hard link from src to dst. If that fails
// (cross-filesystem, unsupported OS, or other error) it falls back to copying.
func linkOrCopy(src, dst string) error {
	// Try hard link first (fastest, no data copy).
	if err := os.Link(src, dst); err == nil {
		return nil
	}

	// Fallback: copy the file.
	return copyFile(src, dst)
}

// findSourceFile looks up a file from a local source directory.
// It tries three layout conventions, returning the first match:
//   1. <SourceDir>/<owner>/<model>/<relative-path>
//   2. <SourceDir>/<model>/<relative-path>
//   3. <SourceDir>/<relative-path>
// Returns (sourcePath, true) if found, or ("", false).
func findSourceFile(sourceDir, repoName, relPath string) (string, bool) {
	candidates := []string{
		filepath.Join(sourceDir, repoName, relPath),
	}

	// If repo is "owner/model", also try <SourceDir>/<model>/<relPath>.
	if slashIdx := strings.LastIndex(repoName, "/"); slashIdx >= 0 {
		model := repoName[slashIdx+1:]
		candidates = append(candidates, filepath.Join(sourceDir, model, relPath))
	}

	// Last resort: flat layout <SourceDir>/<relPath>.
	candidates = append(candidates, filepath.Join(sourceDir, relPath))

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}
	}
	return "", false
}
