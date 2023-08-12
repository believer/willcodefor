package routes

import (
	"log"

	"github.com/believer/willcodefor-go/data"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func IndexHandler(c *fiber.Ctx, db *sqlx.DB) error {
	posts := []Post{}
	err := db.Select(&posts, `SELECT title, til_id, slug, created_at FROM post ORDER BY id DESC LIMIT 5`)

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

func CommandMenuHandler(c *fiber.Ctx, db *sqlx.DB) error {
	search := c.Query("search")
	posts := []Post{}

	q := `
    SELECT title, slug
    FROM post
    WHERE 
      CASE
        WHEN $1 <> '"%%"' THEN title LIKE $1 AND published = true
        ELSE published = true
      END
    ORDER BY id DESC
    LIMIT 5
  `
	err := db.Select(&posts, q, "%"+search+"%")

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	return c.Render("command-menu", fiber.Map{
		"Posts": posts,
	}, "")
}
