package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	_ "github.com/lib/pq"

	"github.com/believer/willcodefor-go/routes"
)

func main() {
	connectionString := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connectionString)

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

	app.Get("/", func(c *fiber.Ctx) error {
		return routes.IndexHandler(c, db)
	})

	app.Get("/command-menu", func(c *fiber.Ctx) error {
		return routes.CommandMenuHandler(c, db)
	})

	app.Get("/posts", func(c *fiber.Ctx) error {
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

	app.Get("/posts/:slug", func(c *fiber.Ctx) error {
		return routes.PostHandler(c, db)
	})

	app.Get("/posts/:id/next", func(c *fiber.Ctx) error {
		return routes.PostNextHandler(c, db)
	})

	app.Get("/posts/:id/previous", func(c *fiber.Ctx) error {
		return routes.PostPreviousHandler(c, db)
	})

	app.Post("/posts/:id/stats", func(c *fiber.Ctx) error {
		return routes.PostStatsHandler(c, db)
	})

	app.Get("/series/:series", func(c *fiber.Ctx) error {
		return routes.PostSeriesHandler(c, db)
	})

	app.Get("/feed.xml", func(c *fiber.Ctx) error {
		return routes.FeedHandler(c, db)
	})

	app.Get("/sitemap.xml", func(c *fiber.Ctx) error {
		return routes.SitemapHandler(c, db)
	})

	app.Get("/iteam", func(c *fiber.Ctx) error {
		return c.Render("iteam", fiber.Map{})
	})

	// Redirects to old page
	app.Get("/stats", func(c *fiber.Ctx) error {
		return c.Redirect("https://willcodefor-htmx.fly.dev/stats", fiber.StatusTemporaryRedirect)
	})

	app.Get("/admin", func(c *fiber.Ctx) error {
		return c.Redirect("https://willcodefor-htmx.fly.dev/admin", fiber.StatusTemporaryRedirect)
	})

	app.Get("/:slug", func(c *fiber.Ctx) error {
		_, err := os.Stat(fmt.Sprintf("./public/%s", c.Params("slug")))

		if err == nil {
			return c.SendFile(fmt.Sprintf("./public/%s", c.Params("slug")))
		}

		return c.Redirect("/posts/"+c.Params("slug"), fiber.StatusSeeOther)
	})

	// Define port if it doesn't exist in env
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
