package model

import (
	"time"

	"github.com/believer/willcodefor-go/utils"
)

type Post struct {
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
	Excerpt   string    `db:"excerpt"`
	ID        int       `db:"id"`
	Series    string    `db:"series"`
	Slug      string    `db:"slug"`
	TILID     int       `db:"til_id"`
	Title     string    `db:"title"`
	UpdatedAt time.Time `db:"updated_at"`
	Views     int       `db:"views"`
}

func (p Post) BodyAsHTML() string {
	body := utils.MarkdownToHTML([]byte(p.Body))
	return body.String()
}

func (p Post) BodyAsXML() string {
	body := utils.MarkdownToXML([]byte(p.Body))
	return body.String()
}

func (p Post) UpdatedAtAsISO() string {
	return p.UpdatedAt.Format(time.RFC3339)
}
