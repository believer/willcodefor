package routes

import (
	"time"

	"github.com/believer/willcodefor-go/data"
	"github.com/believer/willcodefor-go/model"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache/v2"
)

type PostWithParsedDate struct {
	model.Post
	UpdatedAtParsed string
}

func FeedHandler(c *fiber.Ctx) error {
	var (
		posts     []model.Post
		engineXML = mustache.New("./xmls", ".xml")
	)

	if err := engineXML.Load(); err != nil {
		return err
	}

	err := data.Dot.Select(data.DB, &posts, "xml-feed")

	if err != nil {
		return err
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
	var (
		posts     []model.Post
		engineXML = mustache.New("./xmls", ".xml")
	)

	if err := engineXML.Load(); err != nil {
		return err
	}

	err := data.Dot.Select(data.DB, &posts, "xml-sitemap")

	if err != nil {
		return err
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
