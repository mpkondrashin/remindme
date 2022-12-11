package main

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

var GlobalFuncMap = template.FuncMap{
	"duration": Duration,
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
