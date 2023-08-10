package routes

import (
	"database/sql"
	"github.com/believer/willcodefor-go/data"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"log"
)

type NavigationItem struct {
	Title    string
	URL      string
	IsActive bool
}

func IndexHandler(c *fiber.Ctx, db *sql.DB) error {
	var posts []Post

	rows, err := db.Query("SELECT title, til_id, slug, created_at FROM post ORDER BY id DESC LIMIT 5")
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	for rows.Next() {
		var post Post

		rows.Scan(&post.Title, &post.TILID, &post.Slug, &post.CreatedAt)

		posts = append(posts, post)
	}

	return c.Render("index", fiber.Map{
		"Posts":    posts,
		"Projects": data.Projects,
		"Work":     data.Positions,
	})
}

func CommandMenuHandler(c *fiber.Ctx, db *sql.DB) error {
	search := c.Query("search")
	var posts []Post

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
	rows, err := db.Query(q, "%"+search+"%")
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	for rows.Next() {
		var post Post

		rows.Scan(&post.Title, &post.Slug)

		posts = append(posts, post)
	}

	return c.Render("command-menu", fiber.Map{
		"Posts": posts,
	}, "")
}
