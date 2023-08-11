package routes

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/believer/willcodefor-go/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/mileusna/useragent"

	_ "github.com/lib/pq"
)

type Post struct {
	Body      string
	CreatedAt time.Time
	Excerpt   string
	ID        int
	Series    string
	Slug      string
	TILID     int
	Title     string
	UpdatedAt time.Time
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

func PostsHandler(c *fiber.Ctx, db *sql.DB) error {
	var posts []Post

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

	rows, err := db.Query(q)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	for rows.Next() {
		var post Post

		rows.Scan(&post.Title, &post.TILID, &post.CreatedAt, &post.UpdatedAt, &post.Slug)

		posts = append(posts, post)
	}

	return c.Render("posts", fiber.Map{
		"Posts":     posts,
		"SortOrder": sortOrder,
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
		"Posts":     posts,
		"SortOrder": "views",
	})
}

func PostHandler(c *fiber.Ctx, db *sql.DB) error {
	slug := c.Params("slug")
	var post Post

	q := `
	   SELECT title, til_id, slug, id, body, created_at, updated_at, COALESCE(series, ''), excerpt
	   FROM post
	   WHERE slug = $1 OR long_slug = $1
	 `

	if err := db.QueryRow(q, slug).Scan(&post.Title, &post.TILID, &post.Slug, &post.ID, &post.Body, &post.CreatedAt, &post.UpdatedAt, &post.Series, &post.Excerpt); err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	body := utils.MarkdownToHTML([]byte(post.Body))
	post.Body = body.String()

	return c.Render("post", post)
}

func PostNextHandler(c *fiber.Ctx, db *sql.DB) error {
	id := c.Params("id")
	var nextPost Post

	q := `
    SELECT title, slug
    FROM post
    WHERE id > $1 AND published = true
    ORDER BY id ASC
    LIMIT 1
   `

	if err := db.QueryRow(q, id).Scan(&nextPost.Title, &nextPost.Slug); err != nil {
		if err == sql.ErrNoRows {
			return c.SendString("<li></li>")
		}

		log.Fatal(err)
		c.JSON("Oh no")
	}

	return c.Render("partials/postNext", nextPost, "")
}

func PostPreviousHandler(c *fiber.Ctx, db *sql.DB) error {
	var prevPost Post

	id := c.Params("id")
	q := `
    SELECT title, slug
    FROM post
    WHERE id < $1 AND published = true
    ORDER BY id DESC
    LIMIT 1
   `

	if err := db.QueryRow(q, id).Scan(&prevPost.Title, &prevPost.Slug); err != nil {
		if err == sql.ErrNoRows {
			return c.SendString("<li></li>")
		}

		log.Fatal(err)
		c.JSON("Oh no")
	}

	return c.Render("partials/postPrev", prevPost, "")
}

func PostStatsHandler(c *fiber.Ctx, db *sql.DB) error {
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

	if err := db.QueryRow(q, id).Scan(&postViews); err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	return c.SendString(postViews)
}

func PostSeriesHandler(c *fiber.Ctx, db *sql.DB) error {
	var posts []Post

	series := c.Params("series")
	slug := c.Query("slug")
	q := `
    SELECT slug, title
    FROM post
    WHERE series = $1 AND published = true
    ORDER BY id DESC
  `

	rows, err := db.Query(q, series)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
		c.JSON("Oh no")
	}

	for rows.Next() {
		var post Post

		rows.Scan(&post.Slug, &post.Title)

		posts = append(posts, post)
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
