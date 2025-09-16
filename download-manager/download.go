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
)

const downloadsFolder = "downloads"

type Status int8

const (
	InProcess Status = iota
	Done
)

type Download struct {
	URL        string
	Name       string
	Size       int64
	Downloaded int64
	Status     Status
}

func NewDownload(URL string) *Download {
	return &Download{URL: URL}
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
	chunkSize := int64(128 * 1024) // 128 KB
	for {
		n, err := io.CopyN(file, resp.Body, chunkSize)
		if n > 0 {
			d.Downloaded += n
		}
		if err != nil {
			if err == io.EOF {
				d.Status = Done
				break
			}
			return ErrMsg(fmt.Errorf("download: %w", err))
		}
	}

	return nil
}

func (d *Download) SizeHumanized() string {
	return humanizeSize(d.Size)
}

func (d *Download) DownloadedHumanized() string {
	return humanizeSize(d.Downloaded)
}

func humanizeSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%.2f K", float64(size)/float64(1024))
	}
	if size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f M", float64(size)/float64(1024*1024))
	}
	return fmt.Sprintf("%.2f G", float64(size)/float64(1024*1024*1024))
}
