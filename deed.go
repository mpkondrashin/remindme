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

// 1./(1.+exp(-(t-t0)/T))
func (d *Deed) PickColor() string {
	//log.Printf("PickColor(%s)", d.Name)
	passed := time.Since(d.Last)
	//log.Printf("%s: pased %v", d.Name, passed)
	delta := passed - d.Period
	//log.Println("delta", delta)
	relative := 6. * float64(delta) / float64(d.Period)
	//log.Println("relative", relative)
	f := 1. / (1. + math.Exp(-relative))
	//log.Println("f", f)
	var p Color
	if f < .5 {
		percent := int(f * 200)
		//log.Println("percent", percent)
		p = Pick(Green, Yellow, percent)
	} else {
		percent := int((f - .5) * 200)
		//log.Println("percent", percent)
		p = Pick(Yellow, Red, percent)
	}
	//log.Println("pick", p.red, p.green, p.blue)
	//log.Println("pick string", p.String())
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
	/*if delta > 0 {
		return fmt.Sprintf("overdue: %v", delta)
	} else {
		return fmt.Sprintf("time left: %v", -delta)
	}*/
}
