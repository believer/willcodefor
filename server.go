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

	// Define port if it doesn't exist in env
	port := os.Getenv("PORT")

	if port == "" {
		port = "4000"
	}

	// Serve static files
	app.Static("/public", "./public")

	// Start server
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
