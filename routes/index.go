package routes

import (
	"github.com/believer/willcodefor-go/data"
	"github.com/believer/willcodefor-go/model"
	"github.com/gofiber/fiber/v2"
)

func IndexHandler(c *fiber.Ctx) error {
	var posts []model.Post

	err := data.Dot.Select(data.DB, &posts, "five-latest-posts")

	if err != nil {
		return err
	}

	return c.Render("index", fiber.Map{
		"Path":     "/",
		"Posts":    posts,
		"Projects": data.Projects,
		"Work":     data.Positions,
	})
}

func CommandMenuHandler(c *fiber.Ctx) error {
	var posts []model.Post

	search := c.Query("search")

	err := data.Dot.Select(data.DB, &posts, "command-menu-search", "%"+search+"%")

	if err != nil {
		return err
	}

	return c.Render("command-menu", fiber.Map{
		"Posts": posts,
	}, "")
}
