package routes

import (
	"cmp"
	"fmt"
	"log"
	"slices"
	"strconv"
	"time"

	"github.com/believer/willcodefor-go/data"
	"github.com/believer/willcodefor-go/model"
	"github.com/gofiber/fiber/v2"
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

func StatsHandler(c *fiber.Ctx) error {
	var bots int
	var totalViews float64
	var lessThanOnePercent int

	timeQuery := c.Query("time", "today")
	timeQueryString := timeToQuery(c.Query("time", "today"))

	err := data.DB.Get(&totalViews, `
SELECT COUNT(*) FROM post_view WHERE is_bot = FALSE AND created_at >= $1
  `, timeQueryString)

	if err != nil {
		log.Fatal(err)
	}

	err = data.DB.Get(&lessThanOnePercent, `
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

	err = data.DB.Get(&bots, `SELECT COUNT(*) FROM post_view WHERE is_bot = true`)

	if err != nil {
		log.Fatal(err)
	}

	return c.Render("stats", fiber.Map{
		"LessThanOnePercent": lessThanOnePercent,
		"Bots":               bots,
		"Time":               timeQuery,
		"Path":               "/stats",
	})
}

func StatsPostHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	var totalViews int
	var biggestDay struct {
		Date  time.Time `db:"date"`
		Count int       `db:"count"`
	}

	err := data.DB.Get(&totalViews, `SELECT COUNT(*) FROM post_view WHERE post_id = $1 AND is_bot = FALSE`, id)

	if err != nil {
		log.Fatal(err)
	}

	err = data.DB.Get(&biggestDay, `
SELECT
  DATE(created_at) AS DATE,
  COUNT(*) AS COUNT
FROM post_view
WHERE post_id = $1 AND is_bot = FALSE
GROUP BY DATE(created_at)
ORDER BY COUNT DESC LIMIT 1`, id)

	if err != nil {
		log.Fatal(err)
	}

	return c.Render("statsPost", fiber.Map{
		"ID":         id,
		"TotalViews": totalViews,
		"BiggestDay": biggestDay,
		"Path":       "/stats",
	})
}

func StatsPostViewsHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var views []CountData

	err := data.DB.Select(&views, `WITH days AS (
	SELECT
		GENERATE_SERIES(
      (SELECT created_at FROM post WHERE id = $1),
			CURRENT_DATE,
			'1 day'::INTERVAL
		)::DATE AS DAY
)
SELECT
  days.day as date,
	TO_CHAR(days.day, 'Mon DD YY') AS label,
	COUNT(pv.id)::INT AS count
FROM
	days
	LEFT JOIN post_view AS pv ON DATE_TRUNC('day', created_at) = days.day
	AND pv.is_bot = FALSE
	AND post_id = $1
GROUP BY 1
ORDER BY 1 ASC`, id)

	if err != nil {
		log.Fatal(err)
	}

	p, err := constructLineGraphFromData(views)

	if err != nil {
		log.Fatal(err)
	}

	return c.Render("partials/graphLine", fiber.Map{
		"D":     p.D,
		"YGrid": p.YGrid,
	}, "")
}

func ViewsPerDay(c *fiber.Ctx) error {
	var err error
	var totalViews float64

	timeQuery := c.Query("time", "today")
	timeQuerySQL := timeToQuery(timeQuery)

	if timeQuery == "week" {
		err = data.DB.Get(&totalViews, `
SELECT COUNT(*) 
FROM post_view 
WHERE is_bot = FALSE AND date_trunc('week', created_at) = date_trunc('week', now());
`)
	} else {
		err = data.DB.Get(&totalViews, `
SELECT COUNT(*) FROM post_view WHERE is_bot = FALSE AND created_at >= $1
  `, timeQuerySQL)
	}

	if err != nil {
		return err
	}

	viewsPerDay := totalViews

	now := time.Now()
	daysThisYear := now.YearDay()
	firstViewDate := time.Date(2022, 6, 8, 17, 41, 0, 0, time.UTC)
	daysSinceFirstView := now.Sub(firstViewDate).Hours() / 24

	switch timeQuery {
	case "week":
		viewsPerDay = totalViews / 7
	case "thirty-days":
		viewsPerDay = totalViews / 30
	case "this-year":
		viewsPerDay = totalViews / float64(daysThisYear)
	case "cumulative":
		viewsPerDay = totalViews / float64(daysSinceFirstView)
	}

	return c.SendString(fmt.Sprintf("%.2f", viewsPerDay))
}

func MostViewedHandler(c *fiber.Ctx) error {
	var posts []model.Post

	err := data.DB.Select(&posts, `
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
		"Path":      "stats",
	}, "")
}

func BrowsersHandler(c *fiber.Ctx) error {
	var userAgents []struct {
		Name    string `db:"browser_name"`
		Count   int    `db:"count"`
		Percent string `db:"percent"`
	}

	timeQuery := timeToQuery(c.Query("time", "today"))

	err := data.DB.Select(&userAgents, `
  SELECT
		browser_name,
		COUNT(*) AS count,
		TO_CHAR(COUNT(*) / SUM(COUNT(*)) OVER() * 100, 'fm99%') as percent
	FROM post_view
	WHERE is_bot = FALSE AND created_at >= $1
	GROUP BY browser_name
	ORDER BY count DESC
	LIMIT 5
    `, timeQuery)

	if err != nil {
		log.Fatal(err)
	}

	return c.Render("partials/userAgents", fiber.Map{
		"UserAgents": userAgents,
	}, "")
}

func OSHandler(c *fiber.Ctx) error {
	var os []struct {
		Name    string `db:"os_name"`
		Count   int    `db:"count"`
		Percent string `db:"percent"`
	}

	timeQuery := timeToQuery(c.Query("time", "today"))

	err := data.DB.Select(&os, `
  SELECT
    os_name,
    COUNT(*) AS count,
    TO_CHAR(COUNT(*) / SUM(COUNT(*)) OVER() * 100, 'fm99%') as percent
  FROM post_view
  WHERE is_bot = FALSE AND created_at >= $1
  GROUP BY os_name
  ORDER BY count DESC
  LIMIT 5
    `, timeQuery)

	if err != nil {
		log.Fatal(err)
	}

	return c.Render("partials/userAgents", fiber.Map{
		"UserAgents": os,
	}, "")
}

func TotalViewsHandler(c *fiber.Ctx) error {
	var err error
	var count int

	timeQuery := c.Query("time", "today")
	timeQuerySQL := timeToQuery(timeQuery)

	if timeQuery == "week" {
		err = data.DB.Get(&count, `
SELECT COUNT(*) 
FROM post_view 
WHERE is_bot = FALSE AND date_trunc('week', created_at) = date_trunc('week', now());
`)
	} else {
		err = data.DB.Get(&count, `
SELECT COUNT(*) FROM post_view WHERE is_bot = FALSE AND created_at >= $1
  `, timeQuerySQL)
	}

	if err != nil {
		log.Fatal(err)
	}

	return c.SendString(fmt.Sprint(count))
}

func MostViewedTodayHandler(c *fiber.Ctx) error {
	var posts []model.Post

	err := data.DB.Select(&posts, `
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
		"Path":      "stats",
	}, "")
}

type CountData struct {
	Date  string `db:"date"`
	Label string `db:"label"`
	Count int    `db:"count"`
}

func ChartHandler(c *fiber.Ctx) error {
	var views []CountData
	var err error

	time := c.Query("time", "today")

	if time == "today" {
		err = data.DB.Select(&views, `WITH days AS (
  SELECT generate_series(CURRENT_DATE, CURRENT_DATE + '1 day'::INTERVAL, '1 hour') AS hour
)

SELECT
	days.hour as date,
  to_char(days.hour, 'HH24:MI') as label,
  count(pv.id)::int as count
FROM days
LEFT JOIN post_view AS pv ON DATE_TRUNC('hour', created_at at time zone 'utc' at time zone 'Europe/Stockholm') = days.hour AND pv.is_bot = false
LEFT JOIN post AS p ON p.id = pv.post_id
GROUP BY 1
ORDER BY 1 ASC`,
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	if time == "week" {
		err = data.DB.Select(&views, `WITH days AS (
  SELECT generate_series(date_trunc('week', current_date), date_trunc('week', current_date) + '6 days'::INTERVAL, '1 day')::DATE as day
)

SELECT
	days.day as date,
  to_char(days.day, 'Mon DD') as label,
  count(pv.id)::int as count
FROM days
LEFT JOIN post_view AS pv ON DATE_TRUNC('day', created_at) = days.day AND pv.is_bot = false
GROUP BY 1
ORDER BY 1 ASC`,
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	if time == "thirty-days" {
		err = data.DB.Select(&views, `WITH days AS (
          SELECT generate_series(CURRENT_DATE - '30 days'::INTERVAL, CURRENT_DATE, '1 day')::DATE AS day
        )

        SELECT
        	days.day as date,
          to_char(days.day, 'Mon DD') as label,
          count(pv.id)::int as count
        FROM days
        LEFT JOIN post_view AS pv ON DATE_TRUNC('day', created_at) = days.day AND pv.is_bot = false
        GROUP BY 1
        ORDER BY 1 ASC`,
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	if time == "this-year" {
		err = data.DB.Select(&views, `WITH months AS (
	SELECT (DATE_TRUNC('year', NOW()) + (INTERVAL '1' MONTH * GENERATE_SERIES(0,11)))::DATE AS MONTH
)

SELECT
  months.month as date,
	to_char(months.month, 'Mon') as label,
  COUNT(pv.id)::int as count
FROM
	months
	LEFT JOIN post_view AS pv ON DATE_TRUNC('month', created_at) = months.month AND pv.is_bot = false
GROUP BY 1
ORDER BY 1 ASC`,
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	if time == "cumulative" {
		err = data.DB.Select(&views, `WITH data AS (
  SELECT
    date_trunc('month', created_at) as month,
    count(1)::int
  FROM post_view WHERE is_bot = false GROUP BY 1
)

select
  month::DATE as date,
	to_char(month, 'Mon YY') as label,
  sum(count) over (order by month asc rows between unbounded preceding and current row)::int as count
from data`,
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	data, err := constructGraphFromData(views)

	if err != nil {
		log.Fatal(err)
	}

	return c.Render("partials/graph", fiber.Map{
		"Bars":     data,
		"Animated": true,
	}, "")
}

func PostsStatsHandler(c *fiber.Ctx) error {
	var posts []CountData

	err := data.DB.Select(&posts, `
WITH months AS (
	SELECT GENERATE_SERIES('2020-01-01', CURRENT_DATE, '1 month') AS MONTH
)
SELECT
	months.month AS date,
	TO_CHAR(months.month, 'Mon YY') AS label,
	COUNT(p.id) AS count
FROM
	months
	LEFT JOIN post AS p ON DATE_TRUNC('month', p.created_at) = months.month
WHERE p.published = true
GROUP BY 1
ORDER BY 1
  `)

	if err != nil {
		log.Fatal(err)
	}

	data, err := constructGraphFromData(posts)

	if err != nil {
		log.Fatal(err)
	}

	return c.Render("partials/graph", fiber.Map{
		"Bars":     data,
		"Animated": true,
	}, "")
}

type Bar struct {
	Label     string
	Value     int
	BarHeight int
	BarWidth  int
	BarX      int
	BarY      int
	LabelX    float64
	LabelY    float64
	ValueX    float64
	ValueY    int
}

func constructGraphFromData(data []CountData) ([]Bar, error) {
	var graphData []Bar

	graphHeight := 200
	graphWidth := 900
	maxCount := calculateMaxCount(data)

	// The data is used for a bar chart, so we need to convert the data
	for i, row := range data {
		var (
			elementsInGraph = graphWidth / len(data)
			// Calcualte the bar Height
			// Subtract 40 from the graph height to make room for the labels
			barHeight = int(float64(row.Count)/float64(maxCount)*float64(graphHeight-40)) - 6
			barWidth  = int(elementsInGraph) - 5

			// Space the bars evenly across the graph
			// Plus one px for border of first bar
			barX = elementsInGraph*i + 1
			barY = graphHeight - barHeight - 26
		)

		if barWidth <= 0 {
			barWidth = elementsInGraph
			barX = barX + 20
		}

		// Position centered on the bar. Subtract 3.4 which is half the width of the text.
		charWidth := 8.67 // Uses tabular nums so all characters are the same width
		numberOfCharsInCount := len(strconv.Itoa(row.Count))
		numberOfCharsInLabel := len(row.Label)

		halfWidthOfCount := charWidth * float64(numberOfCharsInCount) / 2
		halfWidthOfLabel := charWidth * float64(numberOfCharsInLabel) / 2

		valueX := float64(barX+(barWidth/2)) - halfWidthOfCount
		labelX := float64(barX+(barWidth/2)) - halfWidthOfLabel

		// If it's the first bar, we want to position the label at the start of the graph
		if i == 0 {
			labelX = float64(barX)
		}

		// If it's the last bar, we want to position the label at the end of the graph
		if i == len(data)-1 {
			labelX = float64(barX+barWidth) - charWidth*float64(numberOfCharsInLabel)
		}

		// Subtract 8 to put some space between the text and the bar
		valueY := barY - 8
		// 16,5 is the height of the text
		labelY := float64(barY) + float64(barHeight) + 20

		// Add the data to the graphData slice
		graphData = append(graphData, Bar{
			Label:     row.Label,
			Value:     row.Count,
			BarHeight: barHeight,
			BarWidth:  barWidth,
			BarX:      barX,
			BarY:      barY,
			ValueX:    valueX,
			ValueY:    valueY,
			LabelX:    labelX,
			LabelY:    labelY,
		})
	}

	return graphData, nil
}

type GridLine struct {
	Y1    int
	Y2    int
	Label int
}

type LineGraph struct {
	D     string
	YGrid []GridLine
}

func constructLineGraphFromData(data []CountData) (LineGraph, error) {
	graphHeight := 200
	graphWidth := 900
	maxCount := calculateMaxCount(data)
	var yGrid []GridLine

	// Start the path at the bottom left corner
	path := "M 0 " + strconv.Itoa(graphHeight)

	for i, row := range data {
		// Calculate the x and y values for the line
		x := float64(graphWidth) / float64(len(data)) * float64(i)
		y := float64(graphHeight) - float64(row.Count)/float64(maxCount)*float64(graphHeight)

		// Add point to the path
		path += " L " + strconv.FormatFloat(x, 'f', 3, 64) + " " + strconv.FormatFloat(y, 'f', 3, 64)
	}

	spacing := (graphHeight - 20) / 3

	for i := range 3 {
		ii := i + 1

		yGrid = append(yGrid, GridLine{
			Y1:    graphHeight - spacing*ii,
			Y2:    graphHeight - spacing*ii,
			Label: maxCount / 3 * ii,
		})
	}

	return LineGraph{
		D:     path,
		YGrid: yGrid,
	}, nil
}

func calculateMaxCount(data []CountData) int {
	c := slices.MaxFunc(data, func(a, b CountData) int {
		return cmp.Compare(a.Count, b.Count)
	})

	return c.Count
}
