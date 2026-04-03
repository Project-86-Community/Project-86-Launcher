package fs

import (
	"context"
	"fmt"
	"path/filepath"

	grab "github.com/cavaliergopher/grab/v3"
	"github.com/spf13/afero"
)

// DownloadOptions configures a file download.
type DownloadOptions struct {
	// URL to download from.
	URL string
	// Dest is the relative path inside versions/, e.g. "1.2.3/game.zip"
	Dest string
	// Progress is called periodically with bytes done and total. Can be nil.
	Progress func(done, total int64)
}

// Downloader abstracts downloading so tests can swap in a fake.
type Downloader interface {
	Download(ctx context.Context, fs afero.Fs, opts DownloadOptions) error
}

// GrabDownloader is the real implementation using cavaliergopher/grab.
type GrabDownloader struct{}

func (GrabDownloader) Download(ctx context.Context, fs afero.Fs, opts DownloadOptions) error {
	destRel := filepath.Join("versions", opts.Dest)

	// Ensure destination directory exists inside the jail.
	if err := fs.MkdirAll(filepath.Dir(destRel), 0755); err != nil {
		return fmt.Errorf("download: mkdir: %w", err)
	}

	// grab needs a real OS path, not a relative path.
	destAbs, err := RealPath(fs, destRel)
	if err != nil {
		return fmt.Errorf("download: resolve path: %w", err)
	}

	req, err := grab.NewRequest(destAbs, opts.URL)
	if err != nil {
		return fmt.Errorf("download: build request: %w", err)
	}
	req = req.WithContext(ctx)

	// grab auto-resumes if a partial file exists and server supports Range.
	client := grab.NewClient()
	resp := client.Do(req)

	// Stream progress updates if caller wants them.
	if opts.Progress != nil {
		for !resp.IsComplete() {
			opts.Progress(resp.BytesComplete(), resp.Size())
			// small sleep is fine here, grab updates these atomically.
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-resp.Done:
			}
		}
	} else {
		<-resp.Done
	}

	if err := resp.Err(); err != nil {
		return fmt.Errorf("download: %w", err)
	}
	return nil
}

// FakeDownloader writes a stub file, use in tests to avoid real HTTP.
type FakeDownloader struct{}

func (FakeDownloader) Download(_ context.Context, fs afero.Fs, opts DownloadOptions) error {
	dest := filepath.Join("versions", opts.Dest)
	if err := fs.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	return afero.WriteFile(fs, dest, []byte("fake content"), 0644)
}
