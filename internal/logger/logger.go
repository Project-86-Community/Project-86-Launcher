package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"p86l/configs"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/spf13/afero"
)

const (
	minKeep    = 10
	latestName = "latest.log"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

type unixWriter struct{ w io.Writer }

func (u unixWriter) Write(p []byte) (int, error) {
	ts := fmt.Sprintf("[%d] ", time.Now().Unix())
	_, err := u.w.Write(append([]byte(ts), p...))
	return len(p), err
}

type localWriter struct{ w io.Writer }

func (l localWriter) Write(p []byte) (int, error) {
	ts := fmt.Sprintf("[%s] ", time.Now().Format("2006-01-02 15:04:05"))
	_, err := l.w.Write(append([]byte(ts), p...))
	return len(p), err
}

func init() {
	flags := log.Lshortfile
	makeStdout := func(prefix string) *log.Logger {
		return log.New(localWriter{os.Stdout}, prefix, flags)
	}
	Debug = makeStdout("[DEBUG] ")
	Info = makeStdout("[INFO]  ")
	Warn = makeStdout("[WARN]  ")
	Error = makeStdout("[ERROR] ")
}

func Init(afs afero.Fs, fake bool) error {
	if fake {
		// No file logging, only stdout
		makeLogger := func(prefix string) *log.Logger {
			return log.New(
				localWriter{os.Stdout},
				prefix,
				log.Lshortfile,
			)
		}
		Debug = makeLogger("[DEBUG] ")
		Info = makeLogger("[INFO]  ")
		Warn = makeLogger("[WARN] ")
		Error = makeLogger("[ERROR] ")
		Info.Println("logger initialised (fake mode)")
		return nil
	}

	bpfs, ok := afs.(*afero.BasePathFs)
	if !ok {
		return fmt.Errorf("logger: requires BasePathFs")
	}
	logsDir, err := bpfs.RealPath("logs")
	if err != nil {
		return fmt.Errorf("logger: resolve logs dir: %w", err)
	}

	latestPath := filepath.Join(logsDir, latestName)
	if _, err := os.Stat(latestPath); err == nil {
		archived := filepath.Join(
			logsDir,
			fmt.Sprintf("%s%d%s", configs.LogPrefix, time.Now().Unix(), configs.LogExt),
		)
		err = os.Rename(latestPath, archived)
		if err != nil {
			return fmt.Errorf("logger: rename latest.log: %w", err)
		}
	}

	if err := prune(logsDir); err != nil {
		return fmt.Errorf("logger: prune: %w", err)
	}

	f, err := os.OpenFile(latestPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("logger: open latest.log: %w", err)
	}

	fileW := unixWriter{f}
	stdoutW := localWriter{os.Stdout}

	makeLogger := func(prefix string) *log.Logger {
		return log.New(
			io.MultiWriter(fileW, stdoutW),
			prefix,
			log.Lshortfile,
		)
	}

	Debug = makeLogger("[DEBUG] ")
	Info = makeLogger("[INFO]  ")
	Warn = makeLogger("[WARN] ")
	Error = makeLogger("[ERROR] ")

	Info.Println("logger initialised")
	return nil
}

func LogStartup(version string) {
	sep := strings.Repeat("-", 60)
	Info.Println(sep)
	Info.Println("Project 86 Launcher starting up")
	Info.Println(sep)
	Info.Printf("Launcher version : %s", version)
	Info.Printf("Go version       : %s", runtime.Version())
	Info.Printf("OS               : %s", runtime.GOOS)
	Info.Printf("Arch             : %s", runtime.GOARCH)
	Info.Println(sep)
}

func prune(logsDir string) error {
	entries, err := os.ReadDir(logsDir)
	if err != nil {
		return err
	}

	var logs []string
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, configs.LogPrefix) && strings.HasSuffix(name, configs.LogExt) {
			logs = append(logs, filepath.Join(logsDir, name))
		}
	}

	sort.Strings(logs)

	for len(logs) > minKeep {
		_ = os.Remove(logs[0])
		logs = logs[1:]
	}

	return nil
}
