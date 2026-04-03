package fs

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

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
	base.MkdirAll(filepath.Join(root, "versions"), 0755)
	jailed := afero.NewBasePathFs(base, root)

	// Attempt path traversal, which should fail.
	_ = afero.WriteFile(jailed, "../../evil.txt", []byte("escaped"), 0644)

	// Confirm nothing landed outside the jail root.
	exists, _ := afero.Exists(base, "/tmp/evil.txt")
	if exists {
		t.Fatal("path jail failed: file escaped to /tmp/evil.txt")
	}
}

func TestFakeDownload(t *testing.T) {
	fs := NewMem()
	dl := FakeDownloader{}

	err := dl.Download(context.Background(), fs, DownloadOptions{
		URL:  "https://example.com/game.zip",
		Dest: "1.2.3/game.zip",
	})
	if err != nil {
		t.Fatalf("fake download failed: %v", err)
	}

	exists, err := afero.Exists(fs, "versions/1.2.3/game.zip")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("expected versions/1.2.3/game.zip to exist")
	}
}

func TestWriteReadLog(t *testing.T) {
	fs := NewMem()

	if err := afero.WriteFile(fs, "logs/debug.log", []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	data, err := afero.ReadFile(fs, "logs/debug.log")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("got %q want %q", string(data), "hello")
	}
}
