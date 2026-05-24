// SPDX-License-Identifier: GPL-3.0-only
// SPDX-FileCopyrightText: 2026 Project 86 Community

package fs

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

const testURL = "https://github.com/Taliayaya/Project-86/releases/download/v0.0.0-alpha/Project86-v0.0.0-alpha.zip"

func TestParseVersion(t *testing.T) {
	version, err := ParseVersion(testURL)
	if err != nil {
		t.Fatal(err)
	}
	if version != "v0.0.0-alpha" {
		t.Fatalf("got %q want %q", version, "v0.0.0-alpha")
	}
}

func TestMemFSHasDirs(t *testing.T) {
	fs := NewMem()
	for _, dir := range []string{"logs", "versions"} {
		info, err := fs.Stat(dir)
		if err != nil {
			t.Fatalf("dir %q missing: %v", dir, err)
		}
		if !info.IsDir() {
			t.Fatalf("%q is not a directory", dir)
		}
	}
}

func TestPathJailPreventsEscape(t *testing.T) {
	root := t.TempDir()
	base := afero.NewOsFs()
	_ = base.MkdirAll(filepath.Join(root, "versions"), 0755)
	jailed := afero.NewBasePathFs(base, root)

	_ = afero.WriteFile(jailed, "../../evil.txt", []byte("escaped"), 0644)

	exists, _ := afero.Exists(base, filepath.Join(filepath.Dir(root), "evil.txt"))
	if exists {
		t.Fatal("path jail failed: file escaped outside root")
	}
}

func TestFakeDownloadAndExtract(t *testing.T) {
	fs := NewMem()
	dl := FakeDownloader{}
	ex := FakeExtractor{}

	zipPath, err := dl.Download(context.Background(), fs, DownloadOptions{
		URL: testURL,
	})
	if err != nil {
		t.Fatalf("fake download: %v", err)
	}

	// zip should be at versions/v0.0.0-alpha/Project86-v0.0.0-alpha.zip
	expected := filepath.Join("versions", "v0.0.0-alpha", "Project86-v0.0.0-alpha.zip")
	if zipPath != expected {
		t.Fatalf("got path %q want %q", zipPath, expected)
	}

	exists, _ := afero.Exists(fs, zipPath)
	if !exists {
		t.Fatalf("zip file not found at %q", zipPath)
	}

	if err := ex.Extract(context.Background(), fs, zipPath, ExtractOptions{}); err != nil {
		t.Fatalf("fake extract: %v", err)
	}

	marker := filepath.Join("versions", "v0.0.0-alpha", "extracted.marker")
	exists, _ = afero.Exists(fs, marker)
	if !exists {
		t.Fatal("expected extracted.marker after extraction")
	}
}

func TestVersionDir(t *testing.T) {
	fs := NewMem()
	dl := FakeDownloader{}

	zipPath, err := dl.Download(context.Background(), fs, DownloadOptions{
		URL: testURL,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Confirm it landed under versions/v0.0.0-alpha/
	dir := filepath.Dir(zipPath)
	if dir != filepath.Join("versions", "v0.0.0-alpha") {
		t.Fatalf("unexpected dir %q", dir)
	}
}
