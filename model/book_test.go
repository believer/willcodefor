package model

import (
	"database/sql"
	"testing"
	"time"
)

func TestFinishedBook(t *testing.T) {
	tests := []struct {
		FinishedAt sql.NullTime
		Expected   bool
	}{
		{sql.NullTime{Time: time.Now(), Valid: true}, true},
		{sql.NullTime{Time: time.Time{}, Valid: false}, false},
	}

	for _, test := range tests {
		book := Book{FinishedAt: test.FinishedAt}

		if book.Finished() != test.Expected {
			t.Errorf("Expected book to be finished")
		}
	}
}

func TestProgressPercent(t *testing.T) {
	tests := []struct {
		CurrentPage int
		PageCount   int
		Expected    float64
	}{
		{0, 100, 0},
		{33, 200, 16.5},
		{60, 300, 20},
		{100, 400, 25},
	}

	for _, test := range tests {
		book := Book{
			CurrentPage: test.CurrentPage,
			PageCount:   test.PageCount,
		}

		got := book.ProgressPercent()

		if got != test.Expected {

			t.Errorf("Expected progress percent to be %f, got %f", test.Expected, got)
		}

	}
}

func TestWordsPerPage(t *testing.T) {
	tests := []struct {
		WordCount int
		PageCount int
		Expected  float64
	}{
		{0, 100, 0},
		{85_000, 200, 425},
		{12_345, 300, 41},
		{213_000, 400, 533},
	}

	for _, test := range tests {
		book := Book{
			WordCount: test.WordCount,
			PageCount: test.PageCount,
		}

		got := book.WordsPerPage()

		if got != test.Expected {
			t.Errorf("Expected words per page to be %f, got %f", test.Expected, got)
		}
	}
}

func TestWordsRead(t *testing.T) {
	tests := []struct {
		CurrentPage int
		WordCount   int
		PageCount   int
		Expected    float64
	}{
		{0, 100, 100, 0},
		{33, 20_000, 200, 3300},
		{50, 85_000, 300, 14150},
		{100, 400, 800, 100},
	}

	for _, test := range tests {
		book := Book{
			CurrentPage: test.CurrentPage,
			WordCount:   test.WordCount,
			PageCount:   test.PageCount,
		}

		got := book.WordsRead()

		if got != test.Expected {
			t.Errorf("Expected words read to be %f, got %f", test.Expected, got)
		}
	}
}

func TestFormattedWordCount(t *testing.T) {
	tests := []struct {
		WordCount int
		Expected  string
	}{
		{0, "0"},
		// For the test below, the \u00a0 is a non-breaking space character
		{85_000, "85\u00a0000"},
		{12_345, "12\u00a0345"},
		{213_000, "213\u00a0000"},
	}

	for _, test := range tests {
		book := Book{WordCount: test.WordCount}

		got := book.FormattedWordCount()

		if got != test.Expected {
			t.Errorf("Expected formatted word count to be %s, got %s", test.Expected, got)
		}
	}
}

func TestDaysElapsed(t *testing.T) {
	tests := []struct {
		StartedAt  time.Time
		FinishedAt sql.NullTime
		Expected   int
	}{
		{time.Now(), sql.NullTime{Time: time.Now().AddDate(0, 0, 1), Valid: true}, 1},
		{time.Now(), sql.NullTime{Time: time.Now().AddDate(0, 0, 0), Valid: true}, 1},
		{time.Now(), sql.NullTime{Time: time.Time{}, Valid: false}, 1},
	}

	for _, test := range tests {
		book := Book{
			StartedAt:  test.StartedAt,
			FinishedAt: test.FinishedAt,
		}

		got := book.DaysElapsed()

		if got != test.Expected {
			t.Errorf("Expected days elapsed to be %d, got %d", test.Expected, got)
		}
	}
}

func TestPace(t *testing.T) {
	tests := []struct {
		WordCount   int
		CurrentPage int
		PageCount   int
		StartedAt   time.Time
		FinishedAt  sql.NullTime
		Expected    int
	}{
		{100, 50, 100, time.Now(), sql.NullTime{Time: time.Time{}, Valid: false}, 50},
		{85_000, 50, 200, time.Now(), sql.NullTime{Time: time.Time{}, Valid: false}, 21250},
		{12_345, 50, 300, time.Now(), sql.NullTime{Time: time.Time{}, Valid: false}, 2050},
		{213_000, 50, 400, time.Now(), sql.NullTime{Time: time.Time{}, Valid: false}, 26650},
	}

	for _, test := range tests {
		book := Book{
			WordCount:   test.WordCount,
			CurrentPage: test.CurrentPage,
			PageCount:   test.PageCount,
			StartedAt:   test.StartedAt,
			FinishedAt:  test.FinishedAt,
		}

		got := book.Pace()

		if got != test.Expected {
			t.Errorf("Expected pace to be %d, got %d", test.Expected, got)
		}
	}
}

func TestFormattedPace(t *testing.T) {
	tests := []struct {
		WordCount   int
		CurrentPage int
		PageCount   int
		StartedAt   time.Time
		FinishedAt  sql.NullTime
		Expected    string
	}{
		{100, 50, 100, time.Now(), sql.NullTime{Time: time.Time{}, Valid: false}, "50"},
		{85_000, 50, 200, time.Now(), sql.NullTime{Time: time.Time{}, Valid: false}, "21\u00a0250"},
		{12_345, 50, 300, time.Now(), sql.NullTime{Time: time.Time{}, Valid: false}, "2\u00a0050"},
		{213_000, 50, 400, time.Now(), sql.NullTime{Time: time.Time{}, Valid: false}, "26\u00a0650"},
	}

	for _, test := range tests {
		book := Book{
			WordCount:   test.WordCount,
			CurrentPage: test.CurrentPage,
			PageCount:   test.PageCount,
			StartedAt:   test.StartedAt,
			FinishedAt:  test.FinishedAt,
		}

		got := book.FormattedPace()

		if got != test.Expected {
			t.Errorf("Expected formatted pace to be %s, got %s", test.Expected, got)
		}
	}
}

func TestFormattedBookFormats(t *testing.T) {
	tests := []struct {
		BookFormat []string
		Expected   string
	}{
		{[]string{"Paperback", "Ebook"}, "Paperback, Ebook"},
		{[]string{"Hardcover"}, "Hardcover"},
		{[]string{"Audiobook", "Paperback", "Ebook"}, "Audiobook, Paperback, Ebook"},
	}

	for _, test := range tests {
		book := Book{BookFormat: test.BookFormat}

		got := book.FormattedBookFormats()

		if got != test.Expected {
			t.Errorf("Expected formatted book formats to be %s, got %s", test.Expected, got)
		}
	}
}

func TestReadPercentage(t *testing.T) {
	tests := []struct {
		CurrentPage int
		PageCount   int
		Expected    float64
	}{
		{0, 100, 0},
		{33, 200, 17},
		{50, 300, 17},
		{100, 400, 25},
	}

	for _, test := range tests {
		book := Book{
			CurrentPage: test.CurrentPage,
			PageCount:   test.PageCount,
		}

		got := book.ReadPercentage()

		if got != test.Expected {
			t.Errorf("Expected read percentage to be %f, got %f", test.Expected, got)
		}
	}
}
