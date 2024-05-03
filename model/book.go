package model

import (
	"database/sql"
	"time"
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
	Year          time.Time       `db:"year"`
	Days          sql.NullString  `db:"days"`
	Pace          int             `db:"pace"`
}
