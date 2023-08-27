package routes

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func timeToQuery(t string) string {
	switch t {
	case "week":
		return time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	case "thirty-days":
		return time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	case "this-year":
		return time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	case "cumulative":
		return "2000-01-01"
	default:
		return time.Now().Format("2006-01-02")
	}
}

func StatsHandler(c *fiber.Ctx, db *sqlx.DB) error {
	var bots int
	var viewsPerDay float64
	var lessThanOnePercent int

	timeQuery := c.Query("time", "today")

	// Views per day
	err := db.Get(&viewsPerDay, `
    SELECT
    ROUND((COUNT(id) / (max(created_at)::DATE - min(created_at)::DATE + 1)::NUMERIC), 2) as "viewsPerDay"
    FROM post_view
    WHERE is_bot = false
  `)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Get(&lessThanOnePercent, `
WITH views_with_percentage AS(
	SELECT
		COUNT(*) AS "count",
		COUNT(*) / SUM(COUNT(*)) OVER() * 100 AS percent_as_number
	FROM post_view
	WHERE is_bot = FALSE
	GROUP BY browser_name, os_name
)
SELECT SUM(v.count)
FROM views_with_percentage AS v
WHERE percent_as_number <= 1
    `)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Get(&bots, `SELECT COUNT(*) FROM post_view WHERE is_bot = true`)

	if err != nil {
		log.Fatal(err)
	}

	return c.Render("stats", fiber.Map{
		"ViewsPerDay":        viewsPerDay,
		"LessThanOnePercent": lessThanOnePercent,
		"Bots":               bots,
		"Time":               timeQuery,
	})

}

func MostViewedHandler(c *fiber.Ctx, db *sqlx.DB) error {
	var posts []PostWithViews

	err := db.Select(&posts, `
SELECT
  COUNT(*) as views,
  p.title,
  p.slug,
  p.created_at,
  p.id,
  p.updated_at,
  p.til_id
FROM post_view AS pv
INNER JOIN post AS p ON p.id = pv.post_id
WHERE pv.is_bot = false
GROUP BY p.id
ORDER BY views DESC
LIMIT 10
`)

	if err != nil {
		log.Fatal(err)
	}

	return c.Render("partials/postList", fiber.Map{
		"Posts":     posts,
		"SortOrder": "views",
	}, "")
}

func BrowsersHandler(c *fiber.Ctx, db *sqlx.DB) error {
	var userAgents []struct {
		Name    string `db:"browser_name"`
		Count   int    `db:"count"`
		Percent string `db:"percent"`
	}

	timeQuery := timeToQuery(c.Query("time", "today"))

	err := db.Select(&userAgents, `
WITH views_with_percentage AS(
	SELECT
		browser_name,
		COUNT(*) AS count,
		COUNT(*) / SUM(COUNT(*)) OVER() * 100 AS percent_as_number
	FROM post_view
	WHERE is_bot = FALSE AND created_at >= $1
	GROUP BY browser_name
	ORDER BY count DESC
)
SELECT
	v.browser_name,
	v.count,
	TO_CHAR(percent_as_number, 'fm99%') as percent
FROM views_with_percentage AS v
WHERE percent_as_number > 1
    `, timeQuery)

	if err != nil {
		log.Fatal(err)
	}

	return c.Render("partials/userAgents", fiber.Map{
		"UserAgents": userAgents,
	}, "")
}

func OSHandler(c *fiber.Ctx, db *sqlx.DB) error {
	var os []struct {
		Name    string `db:"os_name"`
		Count   int    `db:"count"`
		Percent string `db:"percent"`
	}

	timeQuery := timeToQuery(c.Query("time", "today"))

	err := db.Select(&os, `
WITH views_with_percentage AS(
  SELECT
    os_name,
    COUNT(*) AS count,
    COUNT(*) / SUM(COUNT(*)) OVER() * 100 AS percent_as_number
  FROM post_view
  WHERE is_bot = FALSE AND created_at >= $1
  GROUP BY os_name
  ORDER BY count DESC
)
SELECT
  v.os_name,
  v.count,
  TO_CHAR(percent_as_number, 'fm99%') as percent
FROM views_with_percentage AS v
WHERE percent_as_number > 1
    `, timeQuery)

	if err != nil {
		log.Fatal(err)
	}

	return c.Render("partials/userAgents", fiber.Map{
		"UserAgents": os,
	}, "")
}

func TotalViewsHandler(c *fiber.Ctx, db *sqlx.DB) error {
	var count int

	timeQuery := timeToQuery(c.Query("time", "today"))

	err := db.Get(&count, `
SELECT COUNT(*) FROM post_view WHERE is_bot = FALSE AND created_at >= $1
  `, timeQuery)

	if err != nil {
		log.Fatal(err)
	}

	return c.SendString(fmt.Sprint(count))
}

func MostViewedTodayHandler(c *fiber.Ctx, db *sqlx.DB) error {
	var posts []PostWithViews

	err := db.Select(&posts, `
SELECT
  COUNT(*) as views,
  p.title,
  p.slug,
  p.created_at,
  p.id,
  p.updated_at,
  p.til_id
FROM post_view AS pv
INNER JOIN post AS p ON p.id = pv.post_id
WHERE pv.is_bot = false AND pv.created_at >= CURRENT_DATE
GROUP BY p.id
ORDER BY views DESC
`)

	if err != nil {
		log.Fatal(err)
	}

	return c.Render("partials/postList", fiber.Map{
		"Posts":     posts,
		"SortOrder": "views",
	}, "")
}
