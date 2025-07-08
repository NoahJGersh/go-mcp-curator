package main

import (
	"regexp"
)

// Regular expressions
var BookRE, IssueRE, RangeCapRE *regexp.Regexp

type BookIssue struct {
	book  *Book
	issue string
}

// Stores of existing values
var Books = make(map[string]*Book)
var Issues = make(map[BookIssue]*Issue)

func main() {
	// Parse book name, plus vol. or annual
	BookRE = regexp.MustCompile(`(\S+)(\d+|@)`)
	// Parse issue number and story number
	IssueRE = regexp.MustCompile(`(\d+)(?:\/(\d+))?`)
	// Parse page, panel, and balloon
	RangeCapRE = regexp.MustCompile(`(\d+)(?::(\d+)(?::(\d+))?)?`)
}
