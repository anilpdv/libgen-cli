// Package humanize provides formatting functions.
package humanize

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// Bytes produces a human readable representation of an SI size.
// Bytes(82854982) -> 83 MB
func Bytes(s uint64) string {
	sizes := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(s, 1000, sizes)
}

// IBytes produces a human readable representation of an IEC size.
// IBytes(82854982) -> 79 MiB
func IBytes(s uint64) string {
	sizes := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	return humanateBytes(s, 1024, sizes)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(math.Log(float64(s)) / math.Log(base))
	if int(e) >= len(sizes) {
		e = float64(len(sizes) - 1)
	}
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}
	return fmt.Sprintf(f, val, suffix)
}

// ParseBytes parses a string representation of bytes back to uint64.
func ParseBytes(s string) (uint64, error) {
	lastDigit := 0
	hasPoint := false
	for _, r := range s {
		if !(unicode.IsDigit(r) || r == '.' || r == ',') {
			break
		}
		if r == '.' || r == ',' {
			hasPoint = true
		}
		lastDigit++
	}
	numStr := strings.ReplaceAll(s[:lastDigit], ",", "")
	unit := strings.TrimSpace(s[lastDigit:])
	if !hasPoint {
		n, err := strconv.ParseUint(numStr, 10, 64)
		if err != nil {
			return 0, err
		}
		mult := unitMultiplier(unit)
		return n * mult, nil
	}
	f, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, err
	}
	mult := float64(unitMultiplier(unit))
	return uint64(f * mult), nil
}

func unitMultiplier(unit string) uint64 {
	switch strings.ToLower(unit) {
	case "k", "kb", "kib":
		return 1024
	case "m", "mb", "mib":
		return 1024 * 1024
	case "g", "gb", "gib":
		return 1024 * 1024 * 1024
	case "t", "tb", "tib":
		return 1024 * 1024 * 1024 * 1024
	default:
		return 1
	}
}
