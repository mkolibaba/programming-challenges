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

const (
	downloadsFolder = "downloads"
	chunkSize       = int64(32 * 1000)
)

type Status int8

const (
	Fetching Status = iota
	InProgress
	Done
)

type Download struct {
	URL        string
	Name       string
	Size       Size
	Downloaded Size
	Status     Status
	Start      time.Time
	Finish     time.Time
	Limit      Size
}

type Option func(*Download)

func WithLimit(limit Size) Option {
	return func(d *Download) {
		d.Limit = limit
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
		sz, err := strconv.ParseInt(l, 10, 64)
		if err != nil {
			return ErrMsg(fmt.Errorf("parse content-length header: %w", err))
		}
		d.Size = Size(sz)
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
	d.Status = InProgress
	d.Start = time.Now()
	for {
		if d.Limit > 0 && d.Speed() > d.Limit {
			continue
		}

		n, err := io.CopyN(file, resp.Body, chunkSize)
		if n > 0 {
			d.Downloaded += Size(n)
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

	return DoneMsg{}
}

func (d *Download) Speed() Size {
	return Size(float64(d.Downloaded) / d.Duration().Seconds())
}

func (d *Download) SpeedHumanized() string {
	return fmt.Sprintf("%s/s", d.Speed())
}

func (d *Download) Duration() time.Duration {
	switch d.Status {
	case Fetching:
		return 0
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
