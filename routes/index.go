package routes

import (
	"github.com/believer/willcodefor-go/data"
	"github.com/believer/willcodefor-go/model"
	"github.com/gofiber/fiber/v2"
)

func IndexHandler(c *fiber.Ctx) error {
	var posts []model.Post
	var books []model.Book

	err := data.Dot.Select(data.DB, &posts, "five-latest-posts")

	if err != nil {
		return err
	}

	err = data.Dot.Select(data.DB, &books, "currently-reading")

	if err != nil {
		return err
	}

	return c.Render("index", fiber.Map{
		"Path":     "/",
		"Posts":    posts,
		"Projects": data.Projects,
		"Work":     data.Positions,
		"Books":    books,
		"HasBooks": len(books) > 0,
	})
}

func CommandMenuHandler(c *fiber.Ctx) error {
	var posts []model.Post

	search := c.Query("search")

	err := data.Dot.Select(data.DB, &posts, "command-menu-search", "%"+search+"%")

	if err != nil {
		return err
	}

	return c.Render("command-menu", posts, "")
}
