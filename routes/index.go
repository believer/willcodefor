package routes

import (
	"log"

	"github.com/believer/willcodefor-go/data"
	"github.com/believer/willcodefor-go/model"
	"github.com/gofiber/fiber/v2"
)

func IndexHandler(c *fiber.Ctx) error {
	posts := []model.Post{}
	err := data.DB.Select(&posts, `
    SELECT
      title,
      til_id,
      slug,
      created_at at time zone 'utc' at time zone 'Europe/Stockholm' as created_at
    FROM post
    WHERE published = true
    ORDER BY id DESC
    LIMIT 5
  `)

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	return c.Render("index", fiber.Map{
		"Path":     "/",
		"Posts":    posts,
		"Projects": data.Projects,
		"Work":     data.Positions,
	})
}

func CommandMenuHandler(c *fiber.Ctx) error {
	search := c.Query("search")
	posts := []model.Post{}

	q := `
    SELECT title, slug
    FROM post
    WHERE 
      CASE
        WHEN $1 <> '"%%"' THEN title ILIKE $1 AND published = true
        ELSE published = true
      END
    ORDER BY id DESC
    LIMIT 5
  `
	err := data.DB.Select(&posts, q, "%"+search+"%")

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	return c.Render("command-menu", fiber.Map{
		"Posts": posts,
	}, "")
}
