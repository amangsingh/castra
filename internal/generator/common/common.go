package common

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// EnsureDirs creates the specified directories relative to the base directory.
func EnsureDirs(baseDir string, dirs []string) error {
	for _, dir := range dirs {
		path := filepath.Join(baseDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}
	return nil
}

// WriteTemplateFile reads a file from the template FS and writes it to the destination.
func WriteTemplateFile(templateFS fs.FS, src, dest string) error {
	content, err := fs.ReadFile(templateFS, src)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", src, err)
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("failed to create parent dir for %s: %w", dest, err)
	}

	if err := os.WriteFile(dest, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", dest, err)
	}
	return nil
}

// MapFunc is a function that maps a template path to a destination path.
// If skip is true, the entry is skipped.
type MapFunc func(path string, d fs.DirEntry) (destPath string, skip bool, err error)

// WalkAndMap walks the template FS starting at srcDir and maps files to the destination.
func WalkAndMap(templateFS fs.FS, srcDir, destBaseDir string, mapper MapFunc) error {
	return fs.WalkDir(templateFS, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		destPath, skip, err := mapper(path, d)
		if err != nil {
			return err
		}
		if skip {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			if err := os.MkdirAll(filepath.Join(destBaseDir, destPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
		} else {
			content, err := fs.ReadFile(templateFS, path)
			if err != nil {
				return fmt.Errorf("failed to read template %s: %w", path, err)
			}
			fullDest := filepath.Join(destBaseDir, destPath)
			if err := os.MkdirAll(filepath.Dir(fullDest), 0755); err != nil {
				return fmt.Errorf("failed to create parent dir for %s: %w", fullDest, err)
			}
			if err := os.WriteFile(fullDest, content, 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", fullDest, err)
			}
		}

		return nil
	})
}
