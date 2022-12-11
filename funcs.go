package main

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

var GlobalFuncMap = template.FuncMap{
	"duration":           Duration,
	"duration_number":    DurationNumber,
	"duration_dimension": DurationDimension,
}

func Duration(d time.Duration) string {
	var sb strings.Builder
	sign := ""
	if d < 0 {
		sign = "-"
		d = -d
	}
	second := time.Second
	minute := time.Minute
	hour := time.Hour
	day := 24 * hour
	week := 7 * day

	weeks := d / week
	if weeks > 0 {
		fmt.Fprintf(&sb, " %d weeks", weeks)
		d /= week
	}

	days := d / day
	if days > 0 {
		fmt.Fprintf(&sb, " %d days", days)
		d /= day
	}

	hours := d / hour
	if hours > 0 {
		fmt.Fprintf(&sb, " %dh", hours)
		d /= hour
	}

	minutes := d / minute
	if minutes > 0 {
		fmt.Fprintf(&sb, " %dm", minutes)
		d /= minute
	}

	seconds := d / second
	if seconds > 0 {
		fmt.Fprintf(&sb, " %ds", seconds)
		d /= second
	}

	return sign + sb.String()[1:]
}

func DurationNumberAndDimension(d time.Duration) (int, string) {
	second := time.Second
	minute := time.Minute
	hour := time.Hour
	day := 24 * hour
	week := 7 * day

	weeks := int(d / week)
	if weeks > 0 {
		return weeks, "w"
	}

	days := int(d / day)
	if days > 0 {
		return days, "d"
	}

	hours := int(d / hour)
	if hours > 0 {
		return hours, "h"
	}

	minutes := int(d / minute)
	if minutes > 0 {
		return minutes, "m"
	}

	seconds := int(d / second)
	return seconds, "s"
}

func DurationNumber(d time.Duration) int {
	num, _ := DurationNumberAndDimension(d)
	return num
}
func DurationDimension(d time.Duration) string {
	_, dim := DurationNumberAndDimension(d)
	return dim
}

func PeriodStrToDuration(s string) time.Duration {
	return map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
		"d": time.Hour * 24,
		"w": time.Hour * 24 * 7,
	}[s]
}
