package routes

import (
	"bytes"
	"database/sql"
	"log"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/gofiber/fiber/v2"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/anchor"

	_ "github.com/lib/pq"
)

type Post struct {
	Body      string
	CreatedAt time.Time
	Excerpt   string
	ID        int
	Series    sql.NullString
	Slug      string
	TILID     int
	Title     string
	UpdatedAt time.Time
}

type PostWithDates struct {
	Post
	DateTime string
	Date     string
}

type PostWithViews struct {
	Post
	Views int
}

func PostsHandler(c *fiber.Ctx, db *sql.DB) error {
	var posts []Post

	sortOrder := c.Query("sort", "createdAt")
	q := `
    SELECT title, til_id, created_at, updated_at, slug
    FROM post
    WHERE published = true
    ORDER BY created_at DESC
  `

	if sortOrder == "updatedAt" {
		q = `
    SELECT title, til_id, created_at, updated_at, slug
    FROM post
    WHERE published = true
    ORDER BY updated_at DESC
  `
	}

	rows, err := db.Query(q)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	for rows.Next() {
		var post Post

		rows.Scan(&post.Title, &post.TILID, &post.CreatedAt, &post.UpdatedAt, &post.Slug)

		// Parse dates
		// if sortOrder == "createdAt" {
		// 	dateToParse = createdAt
		// } else if sortOrder == "updatedAt" {
		// 	dateToParse = updatedAt
		// }
		//
		// parsedDate, _ := time.Parse(time.RFC3339, dateToParse)
		// post.DateTime = parsedDate.Format("2006-01-02 15:04")
		// post.Date = parsedDate.Format("2006-01-02")

		posts = append(posts, post)
	}

	return c.Render("posts", fiber.Map{
		"Posts":      posts,
		"IsTimeSort": true,
		"SortOrder":  sortOrder,
	})
}

func PostsViewsHandler(c *fiber.Ctx, db *sql.DB) error {
	var posts []PostWithViews

	q := `
    SELECT
      p.title,
      p.til_id,
      p.slug,
      COUNT(pv.id) AS views
    FROM post AS p
    INNER JOIN post_view AS pv ON p.id = pv.post_id
    WHERE p.published = true AND pv.is_bot = false
    GROUP BY p.id
    ORDER BY views DESC
  `

	rows, err := db.Query(q)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	for rows.Next() {
		var post PostWithViews

		rows.Scan(&post.Title, &post.TILID, &post.Slug, &post.Views)

		posts = append(posts, post)
	}

	return c.Render("posts", fiber.Map{
		"Posts":      posts,
		"IsTimeSort": false,
		"SortOrder":  "views",
	})
}

type customTexter struct{}

func (*customTexter) AnchorText(h *anchor.HeaderInfo) []byte {
	if h.Level == 1 {
		return nil
	}
	return []byte("#")
}

func PostHandler(c *fiber.Ctx, db *sql.DB) error {
	slug := c.Params("slug")
	var post Post

	q := `
	   SELECT title, til_id, slug, id, body, created_at, updated_at, series, excerpt
	   FROM post
	   WHERE slug = $1 OR long_slug = $1
	 `

	if err := db.QueryRow(q, slug).Scan(&post.Title, &post.TILID, &post.Slug, &post.ID, &post.Body, &post.CreatedAt, &post.UpdatedAt, &post.Series, &post.Excerpt); err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	// Markdown rendering
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			&anchor.Extender{
				Attributer: anchor.Attributes{
					"class": "!text-gray-400 dark:!text-gray-500 no-underline",
				},
				Texter: anchor.Text("#"),
			},
			extension.Strikethrough,
			highlighting.NewHighlighting(
				highlighting.WithStyle("base16-snazzy"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
					chromahtml.LinkableLineNumbers(true, "L"),
				),
			),
		),
	)

	if err := md.Convert([]byte(post.Body), &buf); err != nil {
		panic(err)
	}

	post.Body = buf.String()

	return c.Render("post", post)
}
