package main

import (
	"math"
	"time"

	"github.com/google/uuid"
)

type Deed struct {
	ID     string
	Name   string
	Period time.Duration
	Last   time.Time
}

func NewDeed(name string, period time.Duration) *Deed {
	return &Deed{
		ID:     uuid.New().String(),
		Name:   name,
		Period: period,
		Last:   time.Now(),
	}
}

func (d *Deed) Update() {
	d.Last = time.Now()
}

// PockColor — return color using following formula: 1./(1.+exp(-(t-t0)/T)).
func (d *Deed) PickColor() string {
	passed := time.Since(d.Last)
	delta := passed - d.Period
	relative := 3. * float64(delta) / float64(d.Period)
	f := 1. / (1. + math.Exp(-relative))
	var p Color
	if f < .5 {
		percent := int(f * 200)
		p = Pick(Green, Yellow, percent)
	} else {
		percent := int((f - .5) * 200)
		p = Pick(Yellow, Red, percent)
	}
	return p.String()
}

type DeedModel struct {
	ID      string
	Name    string
	Period  time.Duration
	Color   string
	Overdue time.Duration
}

type DeedsModel []*DeedModel

func (d *Deed) GetModel() *DeedModel {
	return &DeedModel{
		ID:      d.ID,
		Name:    d.Name,
		Period:  d.Period,
		Color:   d.PickColor(),
		Overdue: d.Overdue(),
	}
}

func (d *Deed) Overdue() time.Duration {
	passed := time.Since(d.Last)
	delta := passed - d.Period
	switch {
	case d.Period >= 7*24*time.Hour:
		delta = delta.Round(24 * time.Hour)
	case d.Period >= 24*time.Hour:
		delta = delta.Round(time.Hour)
	case d.Period >= time.Hour:
		delta = delta.Round(time.Minute)
	default:
		delta = delta.Round(time.Second)
	}
	return delta
}
