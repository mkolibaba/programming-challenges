package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const downloadsFolder = "downloads"

type Status int8

const (
	InProgress Status = iota
	Done
)

type Download struct {
	URL        string
	Name       string
	Size       int64
	Downloaded int64
	Status     Status
	Start      time.Time
	Finish     time.Time
	Limit      float64
}

type Option func(*Download)

func WithLimit(limit int64) Option {
	return func(d *Download) {
		d.Limit = float64(limit)
	}
}

func NewDownload(URL string, options ...Option) *Download {
	d := &Download{URL: URL}
	for _, option := range options {
		option(d)
	}
	return d
}

func (d *Download) Run() tea.Msg {
	// Get resource
	resp, err := http.Get(d.URL)
	defer resp.Body.Close()
	if err != nil {
		return ErrMsg(fmt.Errorf("get resource: %w", err))
	}

	// Get file name
	cd := resp.Header.Get("Content-Disposition")
	if cd != "" {
		d.Name = strings.TrimPrefix(cd, "attachment; filename=") // TODO: улучшить алгоритм
	}

	// Get size
	if l := resp.Header.Get("Content-Length"); l != "" { // TODO: что будет, если хедера нет?
		d.Size, err = strconv.ParseInt(l, 10, 64)
		if err != nil {
			return ErrMsg(fmt.Errorf("parse content-length header: %w", err))
		}
	}

	// Create downloads folder
	if err := os.Mkdir(downloadsFolder, 0777); err != nil && !os.IsExist(err) {
		return ErrMsg(fmt.Errorf("create downloads folder: %w", err))
	}

	// Create file
	file, err := os.Create(filepath.Join(downloadsFolder, d.Name))
	if err != nil {
		return ErrMsg(fmt.Errorf("create file: %w", err))
	}
	defer file.Close()

	// Download
	d.Start = time.Now()
	chunkSize := int64(32 * 1024) // 32 K
	for {
		if d.Limit > 0 && d.Speed() > d.Limit {
			continue
		}

		n, err := io.CopyN(file, resp.Body, chunkSize)
		if n > 0 {
			d.Downloaded += n
		}
		if err != nil {
			if err == io.EOF {
				d.Status = Done
				d.Finish = time.Now()
				break
			}
			return ErrMsg(fmt.Errorf("download: %w", err))
		}
	}

	return nil
}

func (d *Download) SizeHumanized() string {
	return humanizeSize(float64(d.Size))
}

func (d *Download) DownloadedHumanized() string {
	return humanizeSize(float64(d.Downloaded))
}

func (d *Download) Speed() float64 {
	return float64(d.Downloaded) / d.Duration().Seconds()
}

func (d *Download) SpeedHumanized() string {
	return fmt.Sprintf("%s/s", humanizeSize(d.Speed()))
}

func (d *Download) Duration() time.Duration {
	switch d.Status {
	case InProgress:
		return time.Now().Sub(d.Start)
	case Done:
		return d.Finish.Sub(d.Start)
	}
	panic("Unknown download status")
}

func (d *Download) DurationHumanized() string {
	duration := d.Duration()
	hours := int64(duration.Hours())
	minutes := int64(duration.Minutes()) % 60
	seconds := int64(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func humanizeSize(size float64) string {
	if size < 1024 {
		return fmt.Sprintf("%.2f B", size)
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%.2f KB", size/1024)
	}
	if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", size/float64(1024*1024))
	}
	return fmt.Sprintf("%.2f GB", size/float64(1024*1024*1024))
}
