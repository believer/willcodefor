package model

import (
	"database/sql"
	"math"
	"strings"
	"time"

	"github.com/lib/pq"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

type BookFormat string

const (
	Physical BookFormat = "physical"
	Kindle   BookFormat = "kindle"
	Audio    BookFormat = "audio"
)

type Book struct {
	ID            int             `db:"id"`
	Author        string          `db:"author"`
	Title         string          `db:"title"`
	Subtitle      sql.NullString  `db:"subtitle"`
	CoverURL      string          `db:"cover_url"`
	Series        sql.NullString  `db:"series"`
	OrderInSeries sql.NullFloat64 `db:"order_in_series"`
	Rating        sql.NullInt32   `db:"rating"`
	StartedAt     time.Time       `db:"started_at"`
	FinishedAt    sql.NullTime    `db:"finished_at"`
	WordCount     int             `db:"word_count"`
	PageCount     int             `db:"page_count"`
	CurrentPage   int             `db:"current_page"`
	ReleaseDate   time.Time       `db:"release_date"`
	RecommendedBy sql.NullString  `db:"recommended_by"`
	BookFormat    pq.StringArray  `db:"book_format"`
}

func (b Book) Finished() bool {
	return b.FinishedAt.Valid
}

func (b Book) ProgressPercent() float64 {
	return (float64(b.CurrentPage) / float64(b.PageCount)) * 100
}

func (b Book) WordsPerPage() int {
	return b.WordCount / b.PageCount
}

func (b Book) WordsRead() int {
	return b.WordsPerPage() * b.CurrentPage
}

func (b Book) FormattedWordCount() string {
	p := message.NewPrinter(language.Swedish)

	return p.Sprintf("%v", number.Decimal(b.WordCount))
}

func (b Book) DaysElapsed() int {
	elapsed := time.Since(b.StartedAt).Hours() / 24

	if b.Finished() {
		elapsed = b.FinishedAt.Time.Sub(b.StartedAt).Hours() / 24
	}

	if elapsed == 0 {
		elapsed = 1
	}

	return int(elapsed)
}

func (b Book) Pace() int {
	return int(float64(b.WordsRead()) / float64(b.DaysElapsed()))
}

func (b Book) FormattedPace() string {
	p := message.NewPrinter(language.Swedish)

	return p.Sprintf("%v", number.Decimal(b.Pace()))
}

func (b Book) ExpectedFinish() time.Time {
	wordsLeft := b.WordCount - b.WordsRead()
	daysLeft := math.Round(float64(wordsLeft) / float64(b.Pace()))

	return time.Now().AddDate(0, 0, int(daysLeft))
}

func (b Book) DaysLeft() int {
	return int(math.Round(b.ExpectedFinish().Sub(time.Now()).Hours() / 24))
}

func (b Book) FormattedBookFormats() string {
	return strings.Join(b.BookFormat, ", ")
}
