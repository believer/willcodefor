package routes

import (
	"database/sql"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/believer/willcodefor-go/data"
	"github.com/believer/willcodefor-go/model"
	"github.com/believer/willcodefor-go/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/mileusna/useragent"
)

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func PostsHandler(c *fiber.Ctx) error {
	var (
		posts     []model.Post
		sortOrder = c.Query("sort", "createdAt")
	)

	q := "posts-by-created"

	if sortOrder == "updatedAt" {
		q = "posts-by-updated"
	}

	err := data.Dot.Select(data.DB, &posts, q)

	if err != nil {
		return err
	}

	return c.Render("posts", fiber.Map{
		"Path":      "/posts",
		"Posts":     posts,
		"SortOrder": sortOrder,
	})
}

func PostsSearchHandler(c *fiber.Ctx) error {
	var (
		posts  []model.Post
		search = c.Query("search")
	)

	err := data.Dot.Select(data.DB, &posts, "post-search", search)

	if err != nil {
		return err
	}

	return c.Render("posts", fiber.Map{
		"SortOrder": "createdAt",
		"Posts":     posts,
		"Search":    search,
	})
}

func PostsViewsHandler(c *fiber.Ctx) error {
	var posts []model.Post

	err := data.Dot.Select(data.DB, &posts, "posts-views")

	if err != nil {
		return err
	}

	return c.Render("posts", fiber.Map{
		"Path":      "/posts",
		"SortOrder": "views",
		"Posts":     posts,
	})
}

func PostHandler(c *fiber.Ctx) error {
	var post model.Post

	// Self healing slug
	var (
		slug     = c.Params("slug")
		parts    = strings.Split(path.Base(slug), "-")
		lastPart = parts[len(parts)-1]
	)

	tilId, err := strconv.Atoi(lastPart)

	if err != nil {
		tilId = 0
	}

	if err := data.Dot.Get(data.DB, &post, "post-by-slug", slug, tilId); err != nil {
		if err == sql.ErrNoRows {
			return c.Render("404", fiber.Map{
				"Slug": slug,
			})
		}

		return err
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

func PostNextHandler(c *fiber.Ctx) error {
	var nextPost model.Post

	id := c.Params("id")

	if err := data.Dot.Get(data.DB, &nextPost, "next-post", id); err != nil {
		if err == sql.ErrNoRows {
			return c.SendString("<li></li>")
		}

		return err
	}

	return c.Render("partials/postNext", nextPost, "")
}

func PostPreviousHandler(c *fiber.Ctx) error {
	var prevPost model.Post

	id := c.Params("id")

	if err := data.Dot.Get(data.DB, &prevPost, "previous-post", id); err != nil {
		if err == sql.ErrNoRows {
			return c.SendString("<li></li>")
		}

		return err
	}

	return c.Render("partials/postPrev", prevPost, "")
}

func PostStatsHandler(c *fiber.Ctx) error {
	var postViews string

	id := c.Params("id")
	env := os.Getenv("APP_ENV")
	userAgent := c.GetReqHeaders()["User-Agent"][0]

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

		_, err := data.Dot.Exec(data.DB, "insert-view",
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
			return err
		}
	}

	if err := data.Dot.Get(data.DB, &postViews, "post-views", id); err != nil {
		return err
	}

	return c.SendString(postViews)
}

func PostSeriesHandler(c *fiber.Ctx) error {
	var (
		posts  []model.Post
		series = c.Params("series")
		slug   = c.Query("slug")
	)

	err := data.Dot.Select(data.DB, &posts, "post-series", series)

	if err != nil {
		return err
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
