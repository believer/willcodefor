package routes

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/believer/willcodefor-go/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mileusna/useragent"

	_ "github.com/lib/pq"
)

type Post struct {
	Body      string
	CreatedAt time.Time `db:"created_at"`
	Excerpt   string
	ID        int
	Series    string
	Slug      string
	TILID     int `db:"til_id"`
	Title     string
	UpdatedAt time.Time `db:"updated_at"`
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

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func PostsHandler(c *fiber.Ctx, db *sqlx.DB) error {
	posts := []Post{}
	sortOrder := c.Query("sort", "createdAt")
	q := `
    SELECT title, til_id, created_at, updated_at, slug
    FROM post
    WHERE published = true
    ORDER BY created_at DESC
  `

	if sortOrder == "updatedAt" {
		q = `
    SELECT title, til_id, created_at, updated_at, slug
    FROM post
    WHERE published = true
    ORDER BY updated_at DESC
  `
	}

	err := db.Select(&posts, q)

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	return c.Render("posts", fiber.Map{
		"Path":      "/posts",
		"Posts":     posts,
		"SortOrder": sortOrder,
	})
}

func PostsSearchHandler(c *fiber.Ctx, db *sqlx.DB) error {
	posts := []Post{}
	search := c.Query("search")
	q := `
    SELECT title, til_id, created_at, updated_at, slug
    FROM post
    WHERE
      title ILIKE '%' || $1 || '%'
      OR body ILIKE '%' || $1 || '%'
      AND published = TRUE
    ORDER BY created_at DESC
  `

	err := db.Select(&posts, q, search)

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	return c.Render("posts", fiber.Map{
		"Posts":     posts,
		"SortOrder": "createdAt",
		"Search":    search,
	})
}

func PostsViewsHandler(c *fiber.Ctx, db *sqlx.DB) error {
	posts := []PostWithViews{}
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

	err := db.Select(&posts, q)

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	return c.Render("posts", fiber.Map{
		"Path":      "/posts",
		"Posts":     posts,
		"SortOrder": "views",
	})
}

func PostHandler(c *fiber.Ctx, db *sqlx.DB) error {
	slug := c.Params("slug")
	post := Post{}

	stmt, err := db.Preparex(`
 SELECT title, til_id, slug, id, body, created_at, updated_at, COALESCE(series, '') as series, excerpt
 FROM post
 WHERE slug = $1 OR long_slug = $1
  `)

	if err != nil {
		log.Fatal(err)
	}

	if err := stmt.Get(&post, slug); err != nil {
		if err == sql.ErrNoRows {
			return c.Render("404", fiber.Map{
				"Slug": slug,
			})
		}

		log.Fatal(err)
		c.JSON("Oh no")
	}

	body := utils.MarkdownToHTML([]byte(post.Body))
	post.Body = body.String()

	return c.Render("post", fiber.Map{
		"Path": "/posts",
		"Post": post,
		"Metadata": fiber.Map{
			"Excerpt": post.Excerpt,
			"Slug":    post.Slug,
			"Title":   post.Title,
		},
	})
}

func PostNextHandler(c *fiber.Ctx, db *sqlx.DB) error {
	id := c.Params("id")
	nextPost := Post{}
	stmt, err := db.Preparex(`
    SELECT title, slug
    FROM post
    WHERE id > $1 AND published = true
    ORDER BY id ASC
    LIMIT 1
   `)

	if err != nil {
		log.Fatal(err)
	}

	if err := stmt.Get(&nextPost, id); err != nil {
		if err == sql.ErrNoRows {
			return c.SendString("<li></li>")
		}

		log.Fatal(err)
		c.JSON("Oh no")
	}

	return c.Render("partials/postNext", nextPost, "")
}

func PostPreviousHandler(c *fiber.Ctx, db *sqlx.DB) error {
	id := c.Params("id")
	prevPost := Post{}
	stmt, err := db.Preparex(`
    SELECT title, slug
    FROM post
    WHERE id < $1 AND published = true
    ORDER BY id DESC
    LIMIT 1
   `)

	if err != nil {
		log.Fatal(err)
	}

	if err := stmt.Get(&prevPost, id); err != nil {
		if err == sql.ErrNoRows {
			return c.SendString("<li></li>")
		}

		log.Fatal(err)
		c.JSON("Oh no")
	}

	return c.Render("partials/postPrev", prevPost, "")
}

func PostStatsHandler(c *fiber.Ctx, db *sqlx.DB) error {
	var postViews string

	id := c.Params("id")
	env := os.Getenv("APP_ENV")
	userAgent := c.GetReqHeaders()["User-Agent"]

	if userAgent != "" && env == "production" {
		engine := ""
		deviceModel := ""
		deviceVendor := ""
		ua := useragent.Parse(userAgent)

		switch ua.Name {
		case "Firefox":
			engine = "Gecko"
		case "Edge":
		case "Chrome":
			engine = "Blink"
		case "Safari":
			engine = "WebKit"
		}

		switch ua.OS {
		case "macOS":
			deviceModel = "Macintosh"
			deviceVendor = "Apple"
		case "iOS":
			deviceVendor = "Apple"
		}

		_, err := db.Exec(`
		    INSERT INTO post_view (
		      user_agent, post_id, is_bot,
		      browser_name, browser_version,
		      device_type, device_model, device_vendor,
          os_name, os_version,
          engine_version, engine_name
		    )
		    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		  `,
			ua.String,
			id,
			ua.Bot,
			ua.Name,
			ua.Version,
			NewNullString(ua.Device),
			NewNullString(deviceModel),
			NewNullString(deviceVendor),
			ua.OS,
			ua.OSVersion,
			ua.Version,
			NewNullString(engine),
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	q := `SELECT COUNT(*) FROM post_view WHERE post_id = $1`

	if err := db.Get(&postViews, q, id); err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	return c.SendString(postViews)
}

func PostSeriesHandler(c *fiber.Ctx, db *sqlx.DB) error {
	posts := []Post{}
	series := c.Params("series")
	slug := c.Query("slug")
	q := `
    SELECT slug, title
    FROM post
    WHERE series = $1 AND published = true
    ORDER BY id ASC
  `

	err := db.Select(&posts, q, series)

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	seriesNames := map[string]string{
		"applescript": "AppleScript",
		"dataview":    "Dataview",
		"htmx":        "htmx",
		"intl":        "Intl",
		"neovim":      "Neovim",
		"rescript":    "ReScript",
		"tmux":        "tmux",
	}

	return c.Render("partials/series", fiber.Map{
		"Posts":      posts,
		"Slug":       slug,
		"SeriesName": seriesNames[series],
	}, "")
}
