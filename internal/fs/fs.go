package fs

import (
	"fmt"
	"os"
	"p86l/configs"
	"path/filepath"
	"runtime"

	"github.com/spf13/afero"
)

// Dir returns the OS-specific appdata directory.
func Dir() (string, error) {
	var base string
	switch runtime.GOOS {
	case "windows":
		base = os.Getenv("APPDATA")
		if base == "" {
			return "", fmt.Errorf("%%APPDATA%% not set")
		}
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, "Library", "Application Support")
	default:
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			base = xdg
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			base = filepath.Join(home, ".local", "share")
		}
	}
	return filepath.Join(base, configs.Name), nil
}

// New returns a real OS filesystem jailed to the appdata dir.
// All paths are relative to the jail, escapes are prevented.
func New() (afero.Fs, error) {
	root, err := Dir()
	if err != nil {
		return nil, err
	}
	base := afero.NewOsFs()
	for _, dir := range []string{"logs", "versions"} {
		if err := base.MkdirAll(filepath.Join(root, dir), 0755); err != nil {
			return nil, fmt.Errorf("fs: mkdir %s: %w", dir, err)
		}
	}
	return afero.NewBasePathFs(base, root), nil
}

// NewMem returns an in-memory filesystem with the same structure for tests.
func NewMem() afero.Fs {
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("logs", 0755)
	_ = fs.MkdirAll("versions", 0755)
	return fs
}

// RealPath resolves the absolute OS path for a relative jail path.
// Only works on a real (non-memory) jailed FS.
func RealPath(fs afero.Fs, rel string) (string, error) {
	bpfs, ok := fs.(*afero.BasePathFs)
	if !ok {
		return "", fmt.Errorf("fs: RealPath requires a BasePathFs")
	}
	return bpfs.RealPath(rel)
}
