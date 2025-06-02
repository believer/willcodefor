package main

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"

	"github.com/believer/willcodefor-go/data"
	"github.com/believer/willcodefor-go/routes"
	"github.com/believer/willcodefor-go/utils"
)

func main() {
	godotenv.Load()
	err := data.InitDB()

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
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
	engine.AddFunc(
		"formatNumber", func(x int) string {
			return utils.FormatNumber(x)
		},
	)

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	// Index routes
	// ––––––––––––––––––––––––––––––––––––––––

	app.Get("/", routes.IndexHandler)
	app.Get("/books", routes.BooksHandler)
	app.Get("/command-menu", routes.CommandMenuHandler)

	// Posts routes
	// ––––––––––––––––––––––––––––––––––––––––

	posts := app.Group("/posts")

	posts.Get("/", func(c *fiber.Ctx) error {
		if c.Query("search", "") != "" {
			return routes.PostsSearchHandler(c)
		}

		switch c.Query("sort", "createdAt") {
		case "views":
			return routes.PostsViewsHandler(c)
		default:
			return routes.PostsHandler(c)
		}
	})

	posts.Get("/:slug", routes.PostHandler)
	app.Get("/til/:slug", routes.PostHandler)

	posts.Get("/:id/next", routes.PostNextHandler)
	posts.Get("/:id/previous", routes.PostPreviousHandler)
	posts.Post("/:id/stats", routes.PostStatsHandler)

	// Series
	// ––––––––––––––––––––––––––––––––––––––––
	app.Get("/series/:series", routes.PostSeriesHandler)

	// XML
	// ––––––––––––––––––––––––––––––––––––––––

	app.Get("/feed.xml", routes.FeedHandler)
	app.Get("/sitemap.xml", routes.SitemapHandler)

	// Other
	// ––––––––––––––––––––––––––––––––––––––––
	app.Get("/iteam", func(c *fiber.Ctx) error {
		return c.Render("iteam", fiber.Map{})
	})

	// Stats
	// ––––––––––––––––––––––––––––––––––––––––
	stats := app.Group("/stats")

	stats.Get("/", routes.StatsHandler)
	stats.Get("/total-views", routes.TotalViewsHandler)
	stats.Get("/views-per-day", routes.ViewsPerDay)
	stats.Get("/browsers", routes.BrowsersHandler)
	stats.Get("/os", routes.OSHandler)
	stats.Get("/most-viewed", routes.MostViewedHandler)
	stats.Get("/most-viewed-today", routes.MostViewedTodayHandler)
	stats.Get("/chart", routes.ChartHandler)
	stats.Get("/posts", routes.PostsStatsHandler)
	stats.Get("/:id", routes.StatsPostHandler)
	stats.Get("/:id/views", routes.StatsPostViewsHandler)

	// Redirects to old page
	// ––––––––––––––––––––––––––––––––––––––––

	// Handle short URLs and old posts where images were linked
	// from the root folder.
	app.Get("/:slug", func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		filename := fmt.Sprintf("./public/%s", slug)

		if _, err := os.Stat(filename); err == nil {
			return c.SendFile(filename)
		}

		return c.Redirect("/posts/"+slug, fiber.StatusSeeOther)
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	// Serve static files
	app.Static("/public", "./public", fiber.Static{
		MaxAge: 86400,
	})

	// Start server
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
