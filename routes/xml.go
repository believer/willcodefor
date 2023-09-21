package routes

import (
	"log"
	"time"

	"github.com/believer/willcodefor-go/data"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache/v2"
)

type PostWithParsedDate struct {
	Post
	UpdatedAtParsed string
}

func FeedHandler(c *fiber.Ctx) error {
	posts := []Post{}
	engineXML := mustache.New("./xmls", ".xml")

	if err := engineXML.Load(); err != nil {
		log.Fatal(err)
	}

	q := `
    SELECT title, slug, body, updated_at at time zone 'utc' at time zone 'Europe/Stockholm' as updated_at
    FROM post
    WHERE published = true
    ORDER BY created_at DESC
  `

	err := data.DB.Select(&posts, q)

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
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

func SitemapHandler(c *fiber.Ctx) error {
	posts := []Post{}
	engineXML := mustache.New("./xmls", ".xml")

	if err := engineXML.Load(); err != nil {
		log.Fatal(err)
	}

	err := data.DB.Select(&posts, "SELECT slug, updated_at FROM post WHERE published = true ORDER BY created_at DESC")

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	var updatedPosts []PostWithParsedDate

	for _, post := range posts {
		var parsedPost PostWithParsedDate

		parsedPost.Title = post.Title
		parsedPost.UpdatedAtParsed = post.UpdatedAt.Format(time.RFC3339)

		updatedPosts = append(updatedPosts, parsedPost)
	}

	c.Type("xml")

	return engineXML.Render(c, "sitemap", fiber.Map{
		"URL":   "https://willcodefor.beer/",
		"Posts": updatedPosts,
	})

}
