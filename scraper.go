package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

// Endpoints
const baseURL string = `https://www.chronologyproject.com/`
const keyPath string = `key.php`

// Key page data

type keyData struct {
	books map[string]bookData
}

type bookData struct {
	title string
	date  string
}

var KeyData keyData

// Entry page data
type characterData struct {
	names       []string
	appearances []string
}

// Map page IDs to useable data
var EntryData map[string]characterData

// Reverse mapping for character name lookup
var CharacterNames map[string]string

// Colly collectors
var k *colly.Collector // Key
var c *colly.Collector // Character entries

func InitScraper() {
	// Establish blank data stores
	KeyData = keyData{
		books: make(map[string]bookData),
	}
	EntryData = make(map[string]characterData)
	CharacterNames = make(map[string]string)

	// Initialize collectors
	k = colly.NewCollector(
		colly.AllowedDomains(baseURL),
	)

	c = k.Clone()

	// Key data events
	k.OnHTML("table", parseKeyData)
}

func GetKeyData() {
	k.Visit(baseURL + keyPath)
}

func parseKeyData(e *colly.HTMLElement) {
	if caption := e.ChildText("caption"); caption != "TITLE KEY by KEY" {
		return
	}

	// Rows
	// 1:  Column headers
	// 2:  @ (annual)
	// 3:  ' (annual by year)
	// 4+: Actual titles
	books := e.DOM.Filter("tr:nth-child(n+4)")
	for _, book := range books.EachIter() {
		// Code | Title | Date
		code := book.Children().First()
		title := code.Next()
		date := title.Next()

		// Add to data store
		KeyData.books[code.Text()] = bookData{title.Text(), date.Text()}
	}
}

func GetEntryData(page, target string) {
	selector := fmt.Sprintf("div#chrons > p#%s", target)
	c.OnHTML(selector, parseEntryData)
	c.Visit(baseURL + page + ".php")
	c.OnHTMLDetach(selector)
}

func parseEntryData(e *colly.HTMLElement) {

}

func GetAllCharacters() {
	c.OnHTML("div#chrons", func(e *colly.HTMLElement) {
		characters := e.DOM.Find("p")
		for _, character := range characters.EachIter() {
			id, ok := character.Attr("id")
			if !ok {
				continue
			}

			names := strings.Split(character.Find("span.char").Text(), "/")
			appearances := strings.Split(character.Find("span.chron").Text(), "<br>")
			for i, appearance := range appearances {
				appearances[i] = strings.Trim(appearance, `\n`)
			}

			EntryData[id] = characterData{names, appearances}
		}
	})
}
