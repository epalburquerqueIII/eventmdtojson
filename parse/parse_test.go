package parse

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"../download"
)

func TestParse(t *testing.T) {
	const (
		dir  string = "/tmp/events"
		repo string = "https://github.com/nosajio/writing"
	)

	// Ensure events will be available in specified location
	download.RepoToDisk(repo, dir)

	t.Run("Files(<dir>)", func(t *testing.T) {
		events, err := Files(dir)
		if err != nil {
			t.Errorf("Files(%s) failed with an error: %s", dir, err.Error())
		}
		if events == nil || len(events) == 0 {
			t.Errorf("Files(%s) returned an empty result. Should be slice of Parsed types", dir)
		}

		// Test an individual event for evidence of successful parsing
		firstEvent := events[0]
		if firstEvent.Title == "" || len(firstEvent.Title) == 0 {
			t.Errorf("Files(%s) doesn't parse the event Title", dir)
		}
		if strings.Contains(firstEvent.BodyHTML, "<p>") == false {
			t.Errorf("Files(%s) doesn't parse HTML in BodyHTML", dir)
		}

		// Test for specific props like "date" "slug" etc
		if firstEvent.Slug == "" {
			t.Errorf("Files(%s) first item has an empty slug", dir)
		}
		if firstEvent.Date.IsZero() {
			t.Errorf("Files(%s) first item has an empty date", dir)
		}

		// Test the result order (should be chronological new -> old)
		var prevDate time.Time
		for i := range events {
			p := events[i]
			if prevDate.IsZero() {
				prevDate = p.Date
				continue
			}
			if prevDate.Before(p.Date) {
				t.Errorf("Files(%s) doesn't sort events chronologically", dir)
			}
			prevDate = p.Date
		}

		// Test custom markdown tags for images
		imgTagPattern := regexp.MustCompile(`(?im)\%img\[.*\]\(.*\)`)
		// The parser will misparse by assuming img tags are links with %img before them
		imgTagBadParsingPattern := regexp.MustCompile(`(?im)\%img<a`)
		for i := range events {
			p := events[i]
			if imgTagPattern.Match([]byte(p.BodyPlain)) && imgTagBadParsingPattern.Match([]byte(p.BodyHTML)) {
				t.Errorf("Files(%s) doesn't parse custom %%img[]() tags into HTML", dir)
			}
		}
	})
}
