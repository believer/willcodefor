package routes

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/believer/willcodefor-go/data"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
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
	var posts []Post

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
	var posts []Post

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
	}, "")
}

type CountData struct {
	Date  string `db:"date"`
	Label string `db:"label"`
	Count int    `db:"count"`
}

var ToolTipFormatter = `
function (info) {
  var [,,value] =info.value;
	return '<div class="tooltip-title">' + value + '</div>';
}
`

func HeatMapHandler(c *fiber.Ctx) error {
	var weekData []CountData

	err := data.DB.Select(&weekData, `WITH days AS (
    SELECT generate_series(date_trunc('week', current_date), date_trunc('week', current_date) + '6 days'::INTERVAL, '1 hour') as hour
)

SELECT
	extract(isodow FROM days.hour) - 1 as date,
  to_char(days.hour, 'HH24')::int as label,
  count(pv.id)::int as count
FROM days
LEFT JOIN post_view AS pv ON DATE_TRUNC('hour', created_at at time zone 'utc' at time zone 'Europe/Stockholm') = days.hour AND pv.is_bot = false
LEFT JOIN post AS p ON p.id = pv.post_id
GROUP BY 1, days.hour
ORDER BY 1,2 ASC`,
	)

	if err != nil {
		log.Fatal(err)
	}

	weekDays := [...]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	hm := charts.NewHeatMap()

	hm.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "250px",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: false,
		}),
		charts.WithGridOpts(opts.Grid{
			Left:   "8%",
			Right:  "2%",
			Bottom: "20%",
			Top:    "5%",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "category",
			SplitArea: &opts.SplitArea{
				Show: true,
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type: "category",
			Data: weekDays,
			SplitArea: &opts.SplitArea{
				Show: true,
			},
			AxisLine: &opts.AxisLine{
				LineStyle: &opts.LineStyle{
					Color: "#a3a3a3",
				},
			},
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Min:        0,
			Max:        10,
			InRange: &opts.VisualMapInRange{
				Color: []string{
					"#f0f9ff",
					"#e0f2fe",
					"#bae6fd",
					"#7dd3fc",
					"#38bdf8",
					"#0ea5e9",
					"#0284c7",
					"#0369a1",
					"#075985",
					"#0c4a6e",
					"#082f49",
				},
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "item",
			Formatter: opts.FuncOpts(ToolTipFormatter),
		}),
	)

	var heatmapData []opts.HeatMapData

	for _, v := range weekData {
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

	hm.AddSeries("Views for hour", heatmapData)

	hm.PageTitle = "Rickard Natt och Dag"

	return hm.Render(c)
}

func barAxisValues(views []CountData) ([]string, []opts.BarData) {
	var xAxis []string
	var yAxis []opts.BarData

	for _, v := range views {
		xAxis = append(xAxis, v.Label)

		if v.Count == 0 {
			yAxis = append(yAxis, opts.BarData{
				Value: nil,
				Label: &opts.Label{
					Show:  true,
					Color: "#a3a3a3",
				},
				ItemStyle: &opts.ItemStyle{
					Color: "#0ea5e9",
				},
			})
		} else {
			yAxis = append(yAxis, opts.BarData{
				Value: v.Count,
				Label: &opts.Label{
					Show:  true,
					Color: "#a3a3a3",
				},
				ItemStyle: &opts.ItemStyle{
					Color: "#0ea5e9",
				},
			})
		}
	}

	return xAxis, yAxis
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

	bar := charts.NewBar()

	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "250px",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: false,
		}),
		charts.WithGridOpts(opts.Grid{
			Left:   "8%",
			Right:  "2%",
			Bottom: "8%",
			Top:    "8%",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "item",
			Formatter: "{b}: {c}",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			AxisLabel: &opts.AxisLabel{
				Show: false,
			},
			SplitLine: &opts.SplitLine{
				Show: true,
				LineStyle: &opts.LineStyle{
					Type:  "dashed",
					Color: "#262626",
				},
			},
		}),
	)

	xAxis, yAxis := barAxisValues(views)

	// Put data into instance
	bar.SetXAxis(xAxis).AddSeries("Data", yAxis).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "top",
			}),
		)

	bar.PageTitle = "Rickard Natt och Dag"

	return bar.Render(c)
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

	bar := charts.NewBar()

	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "250px",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: false,
		}),
		charts.WithGridOpts(opts.Grid{
			Left:   "8%",
			Right:  "2%",
			Bottom: "8%",
			Top:    "8%",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "item",
			Formatter: "{b}: {c}",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			AxisLabel: &opts.AxisLabel{
				Show: false,
			},
			SplitLine: &opts.SplitLine{
				Show: true,
				LineStyle: &opts.LineStyle{
					Type:  "dashed",
					Color: "#262626",
				},
			},
		}),
	)

	xAxis, yAxis := barAxisValues(posts)

	// Put data into instance
	bar.SetXAxis(xAxis).AddSeries("Data", yAxis).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:     true,
				Position: "top",
			}),
		)

	bar.PageTitle = "Rickard Natt och Dag"

	return bar.Render(c)
}
