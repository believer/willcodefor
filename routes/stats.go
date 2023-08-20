package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func StatsHandler(c *fiber.Ctx, db *sqlx.DB) error {
	var count int
	var bots int
	var viewsPerDay float64
	var lessThanOnePercent int
	var userAgents []struct {
		BrowserName string `db:"browser_name"`
		OsName      string `db:"os_name"`
		Count       int    `db:"count"`
		Percent     string `db:"percent"`
	}

	// Total views
	err := db.Get(&count, `SELECT COUNT(*) FROM post_view WHERE is_bot = false`)

	if err != nil {
		log.Fatal(err)
	}

	// Views per day
	err = db.Get(&viewsPerDay, `
    SELECT
    ROUND((COUNT(id) / (max(created_at)::DATE - min(created_at)::DATE + 1)::NUMERIC), 2) as "viewsPerDay"
    FROM post_view
    WHERE is_bot = false
  `)

	if err != nil {
		log.Fatal(err)
	}

	// User agents
	// const data = await db
	// .select({
	//   browserName: postView.browserName,
	//   osName: postView.osName,
	//   count: sql<number>`COUNT(*)::int`,
	//   percent: sql<number>`COUNT(*) / SUM(COUNT(*)) OVER()`.as('percent'),
	// })
	// .from(postView)
	// .where(
	//   and(
	//     gt(postView.createdAt, timeToSql(query.time as Time)),
	//     eq(postView.isBot, false)
	//   )
	// )
	// .groupBy(postView.browserName, postView.osName)
	// .orderBy(sql`count DESC`)

	err = db.Select(&userAgents, `
WITH views_with_percentage AS(
	SELECT
		browser_name,
		os_name,
		COUNT(*) AS "count",
		COUNT(*) / SUM(COUNT(*)) OVER() * 100 AS percent_as_number
	FROM post_view
	WHERE is_bot = FALSE
	GROUP BY browser_name, os_name
	ORDER BY COUNT DESC
)
SELECT
	v.browser_name,
	v.os_name,
	v.count,
	TO_CHAR(percent_as_number, 'fm99%') as percent
FROM views_with_percentage AS v
WHERE percent_as_number > 1
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

	err = db.Get(&bots, `
SELECT COUNT(*) FROM post_view WHERE is_bot = true
  `)

	if err != nil {
		log.Fatal(err)
	}

	return c.Render("stats", fiber.Map{
		"Count":              count,
		"ViewsPerDay":        viewsPerDay,
		"UserAgents":         userAgents,
		"LessThanOnePercent": lessThanOnePercent,
		"Bots":               bots,
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

	total := 0

	for _, post := range posts {
		total += post.Views
	}

	return c.Render("partials/postList", fiber.Map{
		"Posts":     posts,
		"SortOrder": "views",
		"Total":     total,
	}, "")
}
