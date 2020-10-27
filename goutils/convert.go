package goutils

import (
	"strconv"
)

func String2Float(s string) float64 {
	if value, err := strconv.ParseFloat(s, 64); err == nil {
		return value
	}
	return 0.0
}

func String2Int64(s string) int64 {
	if value, err := strconv.ParseInt(s, 10, 0); err == nil {
		return value
	}
	return 0
}

func String2Int(s string) int {
	if value, err := strconv.Atoi(s); err == nil {
		return value
	}
	return 0
}