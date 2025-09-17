package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Size uint64

const (
	Byte     Size = 1
	Kilobyte Size = 1024 * Byte
	Megabyte Size = 1024 * Kilobyte
	Gigabyte Size = 1024 * Megabyte
)

var sizeRegexp = regexp.MustCompile(`([0-9]*)(\.[0-9]*)?([a-z]+)`)

func (s Size) Bytes() uint64 {
	return uint64(s)
}

func (s Size) Kilobytes() float64 {
	return float64(s) / float64(Kilobyte)
}

func (s Size) Megabytes() float64 {
	return float64(s) / float64(Megabyte)
}

func (s Size) Gigabytes() float64 {
	return float64(s) / float64(Gigabyte)
}

func (s Size) String() string {
	switch {
	case s == 0:
		return "0 B"
	case s >= Gigabyte:
		return format(s.Gigabytes(), "GB")
	case s >= Megabyte:
		return format(s.Megabytes(), "MB")
	case s >= Kilobyte:
		return format(s.Kilobytes(), "KB")
	}

	return fmt.Sprintf("%d B", s)
}

func Parse(s string) (Size, error) {
	if s == "" {
		return 0, fmt.Errorf("datasize: empty string")
	}

	n, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return Size(n), nil
	}

	matches := sizeRegexp.FindStringSubmatch(strings.ToLower(s))
	if len(matches) == 0 {
		return 0, fmt.Errorf("datasize: invalid size format: %s", s)
	}

	n, err = strconv.ParseFloat(matches[1]+matches[2], 64)
	if err != nil {
		return 0, fmt.Errorf("datasize: invalid number: %s", s)
	}

	sz, err := suffixSize(matches[3])
	if err != nil {
		return 0, err
	}

	return Size(n) * sz, nil
}

func format(sz float64, suffix string) string {
	if math.Floor(sz) == sz {
		return fmt.Sprintf("%d %s", int64(sz), suffix)
	}
	return fmt.Sprintf("%.2f %s", sz, suffix)
}

func suffixSize(suffix string) (Size, error) {
	switch suffix {
	case "b":
		return Byte, nil
	case "kb":
		return Kilobyte, nil
	case "mb":
		return Megabyte, nil
	case "gb":
		return Gigabyte, nil
	default:
		return 0, fmt.Errorf("datasize: invalid size suffix: %s", suffix)
	}
}
