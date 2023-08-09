package routes

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

type Post struct {
	Title string
	TILID string
	Slug  string
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
	var posts []PostWithDates

	sortOrder := c.Query("sort", "createdAt")
	q := `
    SELECT title, til_id, created_at, updated_at, slug
    FROM post
    WHERE published = true
    ORDER BY created_at DESC
  `

	switch sortOrder {
	case "updatedAt":
		q = `
    SELECT title, til_id, created_at, updated_at, slug
    FROM post
    WHERE published = true
    ORDER BY updated_at DESC
  `
	case "views":
		q = `
    SELECT
      p.title,
      p.til_id,
      p.created_at,
      p.updated_at,
      p.slug,
      COUNT(pv.id) AS views
    FROM post AS p
    INNER JOIN post_view AS pv ON p.id = pv.post_id
    WHERE p.published = true AND pv.is_bot = false
    GROUP BY p.id
    ORDER BY views DESC
  `
	}

	rows, err := db.Query(q)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	for rows.Next() {
		var post PostWithDates
		var createdAt string
		var updatedAt string
		var dateToParse string

		rows.Scan(&post.Title, &post.TILID, &createdAt, &updatedAt, &post.Slug)

		// Parse dates
		if sortOrder == "createdAt" {
			dateToParse = createdAt
		} else if sortOrder == "updatedAt" {
			dateToParse = updatedAt
		}

		parsedDate, _ := time.Parse(time.RFC3339, dateToParse)
		post.DateTime = parsedDate.Format("2006-01-02 15:04")
		post.Date = parsedDate.Format("2006-01-02")

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
