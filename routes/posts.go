package routes

import (
	"database/sql"
	"path"
	"strconv"
	"strings"

	"github.com/believer/willcodefor-go/data"
	"github.com/believer/willcodefor-go/model"
	"github.com/believer/willcodefor-go/utils"
	"github.com/gofiber/fiber/v2"
)

func PostsHandler(c *fiber.Ctx) error {
	var (
		posts     []model.Post
		sortOrder = c.Query("sort", "createdAt")
		q         = "posts-by-created"
	)

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

	umamiPaths, umamiViews, err := data.UmamiPostViews()

	if err != nil {
		return err
	}

	err = data.Dot.Select(data.DB, &posts, "posts-views", umamiPaths, umamiViews, data.UmamiStart)

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
	var (
		nextPost model.Post
		id       = c.Params("id")
	)

	if err := data.Dot.Get(data.DB, &nextPost, "next-post", id); err != nil {
		if err == sql.ErrNoRows {
			return c.SendString("<li></li>")
		}

		return err
	}

	return c.Render("partials/postNext", nextPost, "")
}

func PostPreviousHandler(c *fiber.Ctx) error {
	var (
		prevPost model.Post
		id       = c.Params("id")
	)

	if err := data.Dot.Get(data.DB, &prevPost, "previous-post", id); err != nil {
		if err == sql.ErrNoRows {
			return c.SendString("<li></li>")
		}

		return err
	}

	return c.Render("partials/postPrev", prevPost, "")
}

func PostViewsHandler(c *fiber.Ctx) error {
	var (
		postViews string
		id        = c.Params("id")
	)

	umamiPaths, umamiViews, err := data.UmamiPostViews()

	if err != nil {
		return err
	}

	err = data.Dot.Get(data.DB, &postViews, "post-views", umamiPaths, umamiViews, data.UmamiStart, id)

	if err != nil {
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
