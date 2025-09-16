package main

import (
	"fmt"
	"math"
)

type Size uint64

const (
	Byte     Size = 1
	Kilobyte Size = 1024 * Byte
	Megabyte Size = 1024 * Kilobyte
	Gigabyte Size = 1024 * Megabyte
)

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

func format(sz float64, suffix string) string {
	if math.Floor(sz) == sz {
		return fmt.Sprintf("%d %s", int64(sz), suffix)
	}
	return fmt.Sprintf("%.2f %s", sz, suffix)
}
