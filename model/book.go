package model

import (
	"database/sql"
	"math"
	"strings"
	"time"

	"github.com/believer/willcodefor-go/utils"
	"github.com/lib/pq"
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
	StartedAt     sql.NullTime    `db:"started_at"`
	FinishedAt    sql.NullTime    `db:"finished_at"`
	WordCount     int             `db:"word_count"`
	PageCount     int             `db:"page_count"`
	CurrentPage   int             `db:"current_page"`
	ReleaseDate   time.Time       `db:"release_date"`
	RecommendedBy sql.NullString  `db:"recommended_by"`
	BookFormat    pq.StringArray  `db:"book_format"`
}

func (b Book) Started() bool {
	return b.StartedAt.Valid
}

func (b Book) Finished() bool {
	return b.FinishedAt.Valid
}

func (b Book) WasRecommended() bool {
	return b.RecommendedBy.Valid
}

func (b Book) ProgressPercent() float64 {
	return (float64(b.CurrentPage) / float64(b.PageCount)) * 100
}

func (b Book) WordsPerPage() float64 {
	return math.Round(float64(b.WordCount) / float64(b.PageCount))
}

func (b Book) WordsRead() float64 {
	return math.Round(b.WordsPerPage() * float64(b.CurrentPage))
}

func (b Book) FormattedWordCount() string {
	return utils.FormatNumber(b.WordCount)
}

func (b Book) DaysElapsed() int {
	if !b.Started() {
		return 0
	}

	elapsed := time.Since(b.StartedAt.Time).Hours() / 24

	if b.Finished() {
		elapsed = b.FinishedAt.Time.Sub(b.StartedAt.Time).Hours() / 24
	}

	if elapsed < 1 {
		elapsed = 1
	}

	return int(elapsed)
}

func (b Book) Pace() int {
	return int(float64(b.WordsRead()) / float64(b.DaysElapsed()))
}

func (b Book) FormattedPace() string {
	return utils.FormatNumber(b.Pace())
}

func (b Book) ExpectedFinish() time.Time {
	wordsLeft := float64(b.WordCount) - b.WordsRead()
	daysLeft := math.Round(float64(wordsLeft) / float64(b.Pace()))

	return time.Now().AddDate(0, 0, int(daysLeft))
}

func (b Book) DaysLeft() int {
	return int(math.Round(b.ExpectedFinish().Sub(time.Now()).Hours() / 24))
}

func (b Book) FormattedBookFormats() string {
	return strings.Join(b.BookFormat, ", ")
}

func (b Book) ReadPercentage() float64 {
	return math.Round((float64(b.CurrentPage) / float64(b.PageCount)) * 100)
}
