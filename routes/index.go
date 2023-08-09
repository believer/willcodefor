package routes

import (
	"database/sql"
	"github.com/believer/willcodefor-go/data"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func IndexHandler(c *fiber.Ctx, db *sql.DB) error {
	var posts []PostWithDates

	rows, err := db.Query("SELECT title, til_id, created_at FROM post ORDER BY id DESC LIMIT 5")
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	for rows.Next() {
		var post PostWithDates
		var createdAt string

		rows.Scan(&post.Title, &post.TILID, &createdAt)

		parsedCreatedAt, _ := time.Parse(time.RFC3339, createdAt)
		post.DateTime = parsedCreatedAt.Format("2006-01-02 15:04")
		post.Date = parsedCreatedAt.Format("2006-01-02")

		posts = append(posts, post)
	}

	return c.Render("index", fiber.Map{
		"Posts":      posts,
		"Projects":   data.Projects,
		"Work":       data.Positions,
		"IsTimeSort": true,
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
