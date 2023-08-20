package main

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"

	"github.com/believer/willcodefor-go/routes"
)

func main() {
	connectionString := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	// Set up Fiber and view engine
	engine := html.New("./views", ".html")

	engine.AddFunc(
		"unescape", func(s string) template.HTML {
			return template.HTML(s)
		},
	)
	engine.AddFunc(
		"add", func(x, y int) int {
			return x + y
		},
	)

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	// Index routes
	// ––––––––––––––––––––––––––––––––––––––––

	app.Get("/", func(c *fiber.Ctx) error {
		return routes.IndexHandler(c, db)
	})

	app.Get("/command-menu", func(c *fiber.Ctx) error {
		return routes.CommandMenuHandler(c, db)
	})

	// Posts routes
	// ––––––––––––––––––––––––––––––––––––––––

	posts := app.Group("/posts")

	posts.Get("/", func(c *fiber.Ctx) error {
		if c.Query("search", "") != "" {
			return routes.PostsSearchHandler(c, db)
		}
		switch c.Query("sort", "createdAt") {
		case "views":
			return routes.PostsViewsHandler(c, db)
		default:
			return routes.PostsHandler(c, db)
		}
	})

	posts.Get("/:slug", func(c *fiber.Ctx) error {
		return routes.PostHandler(c, db)
	})

	posts.Get("/:id/next", func(c *fiber.Ctx) error {
		return routes.PostNextHandler(c, db)
	})

	posts.Get("/:id/previous", func(c *fiber.Ctx) error {
		return routes.PostPreviousHandler(c, db)
	})

	posts.Post("/:id/stats", func(c *fiber.Ctx) error {
		return routes.PostStatsHandler(c, db)
	})

	// Series
	// ––––––––––––––––––––––––––––––––––––––––

	app.Get("/seies/:series", func(c *fiber.Ctx) error {
		return routes.PostSeriesHandler(c, db)
	})

	// XML
	// ––––––––––––––––––––––––––––––––––––––––

	app.Get("/feed.xml", func(c *fiber.Ctx) error {
		return routes.FeedHandler(c, db)
	})

	app.Get("/sitemap.xml", func(c *fiber.Ctx) error {
		return routes.SitemapHandler(c, db)
	})

	app.Get("/iteam", func(c *fiber.Ctx) error {
		return c.Render("iteam", fiber.Map{})
	})

	// Stats

	app.Get("/stats", func(c *fiber.Ctx) error {
		return routes.StatsHandler(c, db)
	})

	app.Get("/stats/most-viewed", func(c *fiber.Ctx) error {
		return routes.MostViewedHandler(c, db)
	})

	app.Get("/stats/most-viewed-today", func(c *fiber.Ctx) error {
		return routes.MostViewedTodayHandler(c, db)
	})

	// Redirects to old page
	// ––––––––––––––––––––––––––––––––––––––––

	app.Get("/admin", func(c *fiber.Ctx) error {
		return c.Redirect("https://willcodefor-htmx.fly.dev/admin", fiber.StatusTemporaryRedirect)
	})

	// Handle short URLs and old posts where images were linked
	// from the root folder.
	app.Get("/:slug", func(c *fiber.Ctx) error {
		filename := fmt.Sprintf("./public/%s", c.Params("slug"))

		if _, err := os.Stat(filename); err == nil {
			return c.SendFile(filename)
		}

		return c.Redirect("/posts/"+c.Params("slug"), fiber.StatusSeeOther)
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "4000"
	}

	// Serve static files
	app.Static("/public", "./public", fiber.Static{
		MaxAge: 86400,
	})

	// Start server
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
