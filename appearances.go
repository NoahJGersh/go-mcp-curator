package main

type Book struct {
	Volume   uint8
	IsAnnual bool
	Name     string
	Abbr     string
}

type Issue struct {
	Number uint16
	Book   *Book
}

type Range struct {
	Start RangeCap
	End   RangeCap
	Issue *Issue
}

type RangeCap struct {
	Page    uint8
	Panel   uint8
	Balloon uint8
}

type Appearance struct {
	Character   string
	Alias       string
	RangeMain   Range
	RangeOthers []Range
}
