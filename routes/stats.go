package routes

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
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
		"Path":               "/stats",
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

func OSHandler(c *fiber.Ctx, db *sqlx.DB) error {
	var os []struct {
		Name    string `db:"os_name"`
		Count   int    `db:"count"`
		Percent string `db:"percent"`
	}

	timeQuery := timeToQuery(c.Query("time", "today"))

	err := db.Select(&os, `
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

type CountData struct {
	Date  string `db:"date"`
	Label string `db:"label"`
	Count int    `db:"count"`
}

func heatMapBase(data []CountData) *charts.HeatMap {
	weekDays := [...]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	hm := charts.NewHeatMap()
	hm.SetGlobalOptions(
		charts.WithLegendOpts(opts.Legend{
			Show: false,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type:      "category",
			SplitArea: &opts.SplitArea{Show: true},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type:      "category",
			Data:      weekDays,
			SplitArea: &opts.SplitArea{Show: true},
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Min:        0,
			Max:        10,
			InRange: &opts.VisualMapInRange{
				Color: []string{"#50a3ba", "#eac736", "#d94e5d"},
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    true,
			Trigger: "item",
		}),
	)
	dayHrs := [...]string{
		"12a", "1a", "2a", "3a", "4a", "5a", "6a", "7a", "8a", "9a", "10a", "11a",
		"12p", "1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p", "10p", "11p",
	}

	var heatmapData []opts.HeatMapData

	for _, v := range data {
		x, _ := strconv.Atoi(v.Label)
		y, _ := strconv.Atoi(v.Date)

		if v.Count == 0 {
			heatmapData = append(heatmapData, opts.HeatMapData{
				Value: []interface{}{x, y, nil},
			})
		} else {
			heatmapData = append(heatmapData, opts.HeatMapData{
				Value: []interface{}{x, y, v.Count},
			})
		}
	}

	hm.SetXAxis(dayHrs).AddSeries("heatmap", heatmapData)

	return hm
}

func ChartHandler(c *fiber.Ctx, db *sqlx.DB) error {
	var views []CountData
	var weekData []CountData

	err := db.Select(&views, `WITH days AS (
  SELECT generate_series(CURRENT_DATE, CURRENT_DATE + '1 day'::INTERVAL, '1 hour') AS hour
)

SELECT
	days.hour as date,
  to_char(days.hour, 'HH24:MI') as label,
  count(pv.id)::int as count
FROM days
LEFT JOIN post_view AS pv ON DATE_TRUNC('hour', created_at) = days.hour
LEFT JOIN post AS p ON p.id = pv.post_id
GROUP BY 1
ORDER BY 1 ASC`,
	)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Select(&weekData, `WITH days AS (
    SELECT generate_series(date_trunc('week', current_date), date_trunc('week', current_date) + '6 days'::INTERVAL, '1 hour') as hour
)

SELECT
	extract(isodow FROM days.hour) - 1 as date,
  to_char(days.hour, 'HH24')::int as label,
  count(pv.id)::int as count
FROM days
LEFT JOIN post_view AS pv ON DATE_TRUNC('hour', created_at) = days.hour
LEFT JOIN post AS p ON p.id = pv.post_id
GROUP BY 1, days.hour
ORDER BY 1,2 ASC`,
	)

	if err != nil {
		log.Fatal(err)
	}

	bar := charts.NewBar()

	bar.SetGlobalOptions(
		charts.WithLegendOpts(opts.Legend{
			Show: false,
		}),
		charts.WithGridOpts(opts.Grid{
			Left:   "5%",
			Right:  "2%",
			Bottom: "5%",
			Top:    "5%",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			SplitLine: &opts.SplitLine{
				Show: true,
				LineStyle: &opts.LineStyle{
					Type:  "dashed",
					Color: "#333",
				},
			},
		}),
	)

	var xAxis []string
	var yAxis []opts.BarData

	for _, v := range views {
		xAxis = append(xAxis, v.Label)
		yAxis = append(yAxis, opts.BarData{
			Value: v.Count,
			ItemStyle: &opts.ItemStyle{
				Color: "#65bcff",
			},
		})
	}

	// Put data into instance
	bar.SetXAxis(xAxis).AddSeries("Data", yAxis).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "top",
			}),
		)

	page := components.NewPage()
	page.AddCharts(
		heatMapBase(weekData),
		bar,
	)

	return page.Render(c)
}
