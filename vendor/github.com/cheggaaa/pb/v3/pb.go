// Package pb provides terminal progress bar functionality.
package pb

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type ProgressBarTemplate int

const (
	Full ProgressBarTemplate = iota
	Default
	Simple
)

func (t ProgressBarTemplate) Start64(total int64) *ProgressBar {
	b := &ProgressBar{
		total:     total,
		startTime: time.Now(),
		writer:    os.Stderr,
	}
	b.render()
	return b
}

type ProgressBar struct {
	total     int64
	current   int64
	startTime time.Time
	mu        sync.Mutex
	writer    io.Writer
	finished  bool
}

func (b *ProgressBar) SetTotal(total int64) *ProgressBar {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.total = total
	return b
}

func (b *ProgressBar) SetCurrent(current int64) *ProgressBar {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current = current
	b.render()
	return b
}

func (b *ProgressBar) Add64(n int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current += n
	b.render()
}

func (b *ProgressBar) NewProxyReader(r io.Reader) io.Reader {
	return &proxyReader{
		Reader: r,
		bar:    b,
	}
}

func (b *ProgressBar) Finish() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.finished {
		return
	}
	b.finished = true
	b.current = b.total
	b.render()
	fmt.Fprintln(b.writer)
}

func (b *ProgressBar) render() {
	width := 30
	var pct float64
	if b.total > 0 {
		pct = float64(b.current) / float64(b.total)
	}
	if pct > 1.0 {
		pct = 1.0
	}
	completed := int(pct * float64(width))
	if completed > width {
		completed = width
	}
	remaining := width - completed

	barStr := strings.Repeat("=", completed)
	if completed < width && completed > 0 {
		barStr = barStr[:len(barStr)-1] + ">"
	}
	emptyStr := strings.Repeat(" ", remaining)

	curHuman := formatBytes(b.current)
	totHuman := formatBytes(b.total)

	fmt.Fprintf(b.writer, "\r[%s%s] %5.1f%% (%s / %s)", barStr, emptyStr, pct*100, curHuman, totHuman)
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

type proxyReader struct {
	io.Reader
	bar *ProgressBar
}

func (pr *proxyReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if n > 0 {
		pr.bar.Add64(int64(n))
	}
	return n, err
}
