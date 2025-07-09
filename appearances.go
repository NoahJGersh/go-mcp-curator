package main

import (
	"strconv"
	"strings"
)

type Book struct {
	Volume   uint8
	IsAnnual bool
	Name     string
	Abbr     string
}

type Issue struct {
	Number      uint16
	PostDecimal string
	Book        *Book
}

type Range struct {
	Flags flags
	Start RangeCap
	End   RangeCap
	Story uint8
	Issue *Issue
}

type RangeCap struct {
	Page    uint8
	Panel   uint8
	Balloon uint8
}

type flags struct {
	isFlashback bool
	isBTS       bool
	isOffPanel  bool
	isVoiceOver bool
}

type Appearance struct {
	Character   string
	Alias       string
	RangeMain   Range
	RangeOthers []Range
}

// Translate a given entry from the MCP listing to more verbose data
func Translate(name, alias, appearanceText string) Appearance {
	// Establish baseline appearance
	appearance := Appearance{
		Character:   name,
		Alias:       alias,
		RangeMain:   Range{},
		RangeOthers: make([]Range, 0),
	}

	// Get all synchronous appearances
	ranges := strings.Split(appearanceText, "~")
	for i, rangeText := range ranges {
		newRange := getRange(strings.Trim(rangeText, " "))
		if i == 0 {
			appearance.RangeMain = newRange
		} else {
			appearance.RangeOthers = append(appearance.RangeOthers, newRange)
		}
	}

	return appearance
}

// Parse an appearance range from the given text
func getRange(rangeText string) Range {
	newRange := Range{}

	fragments := strings.SplitN(rangeText, " ", 3)

	// First fragment: Book
	book := getBook(fragments[0])

	// Second fragment: Issue
	issue, story, _ := getIssue(fragments[1], book)

	// Third fragment: Range Caps
	// TODO

	newRange.Issue = issue
	newRange.Story = story

	return newRange
}

// Get the relevant book from an entry fragment
func getBook(fragment string) *Book {
	// Check if exists
	if book, ok := Books[fragment]; ok {
		return book
	}

	// Build new book from the existing values
	var book Book
	bookValues := BookRE.FindStringSubmatch(fragment)
	if bookValues == nil {
		book = Book{
			Volume:   1,
			IsAnnual: false,
			Name:     getFullName(fragment),
			Abbr:     fragment,
		}
	} else {
		book = Book{
			Name: getFullName(bookValues[1]),
			Abbr: bookValues[1],
		}

		isAnnual := bookValues[2] == "@"

		vol := uint8(0)
		if !isAnnual {
			vol64, _ := strconv.ParseUint(bookValues[2], 10, 8)
			vol = uint8(vol64)
		}

		book.IsAnnual = isAnnual
		book.Volume = vol
	}

	Books[fragment] = &book
	return &book
}

func getIssue(fragment string, book *Book) (issue *Issue, story uint8, flags flags) {
	issueNumber, issueFlags, _ := strings.Cut(fragment, "-")

	// Collate all flags for the appearance in this issue
	flags = getFlags(issueFlags)

	// Get story number, if present
	issueNumber, issueStory, found := strings.Cut(issueNumber, "/")
	story = uint8(1)
	if found {
		story64, _ := strconv.ParseUint(issueStory, 10, 8)
		story = uint8(story64)
	}

	issueIdx := BookIssue{book, issueNumber}

	// Check if issue already tracked
	if iss, ok := Issues[issueIdx]; ok {
		issue = iss
	} else {
		// Parse out post-decimal information
		issueNumber, postDecimal, _ := strings.Cut(issueNumber, ".")
		issueNum64, _ := strconv.ParseUint(issueNumber, 10, 16)
		issue = &Issue{
			Number:      uint16(issueNum64),
			PostDecimal: postDecimal,
			Book:        book,
		}

		// Store new issue
		Issues[issueIdx] = issue
	}

	return
}

// Helper function to get relevant appearance flags from a given string of text
func getFlags(flagString string) flags {
	flags := flags{}
	for flag := range strings.SplitSeq(flagString, "-") {
		switch flag {
		case "FB":
			flags.isFlashback = true
		case "BTS":
			flags.isBTS = true
		case "OP":
			flags.isOffPanel = true
		case "VO":
			flags.isVoiceOver = true
		}
	}

	return flags
}

// TODO: get full name from MCP key page
func getFullName(abbr string) string {
	return abbr
}
