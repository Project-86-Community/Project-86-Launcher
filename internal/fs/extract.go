package fs

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/saracen/fastzip"
	"github.com/spf13/afero"
)

type ExtractOptions struct {
	OnFile func(extracted, total int)
}

type Extractor interface {
	Extract(ctx context.Context, fs afero.Fs, zipPath string, opts ExtractOptions) error
}

type FastExtractor struct{}

func (FastExtractor) Extract(ctx context.Context, afs afero.Fs, zipRelPath string, opts ExtractOptions) error {
	zipAbs, err := RealPath(afs, zipRelPath)
	if err != nil {
		return fmt.Errorf("extract: resolve zip path: %w", err)
	}

	destRel := filepath.Dir(zipRelPath)
	destAbs, err := RealPath(afs, destRel)
	if err != nil {
		return fmt.Errorf("extract: resolve dest path: %w", err)
	}

	e, err := fastzip.NewExtractor(zipAbs, destAbs)
	if err != nil {
		return fmt.Errorf("extract: create extractor: %w", err)
	}

	total := len(e.Files())

	var extractErr error
	if opts.OnFile != nil {
		done := make(chan error, 1)
		go func() { done <- e.Extract(ctx) }()

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

	loop:
		for {
			select {
			case extractErr = <-done:
				_, entries := e.Written()
				opts.OnFile(int(entries), total)
				break loop
			case <-ticker.C:
				_, entries := e.Written()
				opts.OnFile(int(entries), total)
			case <-ctx.Done():
				_ = e.Close()
				return ctx.Err()
			}
		}
	} else {
		extractErr = e.Extract(ctx)
	}

	_ = e.Close()

	if extractErr != nil {
		return fmt.Errorf("extract: %w", extractErr)
	}

	if err := afs.Remove(zipRelPath); err != nil {
		return fmt.Errorf("extract: cleanup zip: %w", err)
	}

	return nil
}

type FakeExtractor struct{}

func (FakeExtractor) Extract(ctx context.Context, afs afero.Fs, zipRelPath string, opts ExtractOptions) error {
	const total = 80 // fake 80 files

	if opts.OnFile != nil {
		for i := range total {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(50 * time.Millisecond):
				opts.OnFile(i+1, total)
			}
		}
	}

	marker := filepath.Join(filepath.Dir(zipRelPath), "extracted.marker")
	if err := afero.WriteFile(afs, marker, []byte("ok"), 0644); err != nil {
		return err
	}
	_ = afs.Remove(zipRelPath)
	return nil
}
