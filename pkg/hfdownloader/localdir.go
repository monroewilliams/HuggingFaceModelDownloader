// Copyright 2025
// SPDX-License-Identifier: Apache-2.0

package hfdownloader

import (
	"os"
	"path/filepath"
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

// findSourceFile looks up a file from an LM Studio–style SourceDir.
// The layout is: <SourceDir>/<owner>/<model>/<relative-path>, which maps
// directly to SourceDir/repo/relPath.  Returns (sourcePath, true) if found,
// or ("", false).
func findSourceFile(sourceDir, repoName, relPath string) (string, bool) {
	candidate := filepath.Join(sourceDir, repoName, relPath)
	if _, err := os.Stat(candidate); err == nil {
		return candidate, true
	}
	return "", false
}
