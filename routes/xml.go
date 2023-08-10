package routes

import (
	"database/sql"
	"log"
	"time"

	"github.com/believer/willcodefor-go/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache/v2"
)

type PostWithParsedDate struct {
	Post
	UpdatedAtParsed string
}

func FeedHandler(c *fiber.Ctx, db *sql.DB) error {
	var posts []PostWithParsedDate

	engineXML := mustache.New("./xmls", ".xml")

	if err := engineXML.Load(); err != nil {
		log.Fatal(err)
	}

	q := `
    SELECT
      title,
      slug,
      body,
      updated_at
    FROM post
    WHERE published = true
    ORDER BY created_at DESC
  `

	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	for rows.Next() {
		var post PostWithParsedDate

		if err := rows.Scan(&post.Title, &post.Slug, &post.Body, &post.UpdatedAt); err != nil {
			log.Fatal(err)
			c.JSON("Oh no")
		}

		body := utils.MarkdownToHTML([]byte(post.Body))
		post.Body = body.String()
		post.UpdatedAtParsed = post.UpdatedAt.Format(time.RFC3339)

		posts = append(posts, post)
	}

	c.Type("xml")

	return engineXML.Render(c, "feed", fiber.Map{
		"Metadata": fiber.Map{
			"Title":       "willcodefor.beer",
			"URL":         "https://willcodefor.beer/",
			"Description": "Things I learn while browsing the web",
			"Author": fiber.Map{
				"Name":  "Rickard Natt och Dag",
				"Email": "rickard@willcodefor.dev",
			},
		},
		"Posts":            posts,
		"LatestPostUpdate": posts[0].UpdatedAt.Format(time.RFC3339),
	})
}

func SitemapHandler(c *fiber.Ctx, db *sql.DB) error {
	var posts []PostWithParsedDate

	engineXML := mustache.New("./xmls", ".xml")

	if err := engineXML.Load(); err != nil {
		log.Fatal(err)
	}

	q := `
    SELECT slug, updated_at
    FROM post
    WHERE published = true
    ORDER BY created_at DESC
  `

	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	for rows.Next() {
		var post PostWithParsedDate

		if err := rows.Scan(&post.Slug, &post.UpdatedAt); err != nil {
			log.Fatal(err)
			c.JSON("Oh no")
		}

		post.UpdatedAtParsed = post.UpdatedAt.Format(time.RFC3339)

		posts = append(posts, post)
	}

	c.Type("xml")

	return engineXML.Render(c, "sitemap", fiber.Map{
		"URL":   "https://willcodefor.beer/",
		"Posts": posts,
	})

}
