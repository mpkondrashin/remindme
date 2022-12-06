package main

import "fmt"

var Red = Color{255, 0, 0}
var Green = Color{0, 255, 0}
var Yellow = Color{255, 255, 0}

type Color struct {
	red   int
	green int
	blue  int
}

func Pick(a, b Color, percent int) Color {
	return Color{
		red:   (a.red*(100-percent) + b.red*percent) / 100,
		green: (a.green*(100-percent) + b.green*percent) / 100,
		blue:  (a.blue*(100-percent) + b.blue*percent) / 100,
	}
}

func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x", c.red, c.green, c.blue)
}
