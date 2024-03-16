package routes

import (
	"fmt"
	"time"

	"github.com/believer/willcodefor-go/data"
	"github.com/believer/willcodefor-go/model"
	"github.com/believer/willcodefor-go/utils"
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

	timeQuery := c.Query("time", "today")
	timeQueryString := timeToQuery(c.Query("time", "today"))

	err := data.Dot.Get(data.DB, &totalViews, "stats-views-for-period", timeQueryString)

	if err != nil {
		return err
	}

	err = data.Dot.Get(data.DB, &bots, "stats-bots")

	if err != nil {
		return err
	}

	return c.Render("stats", fiber.Map{
		"Bots": bots,
		"Time": timeQuery,
		"Path": "/stats",
	})
}

func StatsPostHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	var totalViews int
	var biggestDay struct {
		Date  time.Time `db:"date"`
		Count int       `db:"count"`
	}

	err := data.Dot.Get(data.DB, &totalViews, "stats-post-total-views", id)

	if err != nil {
		return err
	}

	err = data.Dot.Get(data.DB, &biggestDay, "stats-post-biggest-day", id)

	if err != nil {
		return err
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
	var views []utils.CountData

	err := data.Dot.Select(data.DB, &views, "stats-post-views", id)

	if err != nil {
		return err
	}

	p, err := utils.LineChart(views)

	if err != nil {
		return err
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
		err = data.Dot.Get(data.DB, &totalViews, "stats-views-per-week")
	} else {
		err = data.Dot.Get(data.DB, &totalViews, "stats-views-for-period", timeQuerySQL)
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

	err := data.Dot.Select(data.DB, &posts, "stats-most-viewed-posts")

	if err != nil {
		return err
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

	err := data.Dot.Select(data.DB, &userAgents, "stats-browsers", timeQuery)

	if err != nil {
		return err
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

	err := data.Dot.Select(data.DB, &os, "stats-os", timeQuery)

	if err != nil {
		return err
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
		err = data.Dot.Get(data.DB, &count, "stats-views-per-week")
	} else {
		err = data.Dot.Get(data.DB, &count, "stats-views-for-period", timeQuerySQL)
	}

	if err != nil {
		return err
	}

	return c.SendString(fmt.Sprint(count))
}

func MostViewedTodayHandler(c *fiber.Ctx) error {
	var posts []model.Post

	err := data.Dot.Select(data.DB, &posts, "stats-most-viewed-posts-today")

	if err != nil {
		return err
	}

	return c.Render("partials/postList", fiber.Map{
		"Posts":     posts,
		"SortOrder": "views",
		"Path":      "stats",
	}, "")
}

func ChartHandler(c *fiber.Ctx) error {
	var views []utils.CountData
	var err error

	time := c.Query("time", "today")

	query := "stats-chart-today"

	switch time {
	case "week":
		query = "stats-chart-week"
	case "thirty-days":
		query = "stats-chart-thirty-days"
	case "this-year":
		query = "stats-chart-this-year"
	case "cumulative":
		query = "stats-chart-all-time"
	}

	err = data.Dot.Select(data.DB, &views, query)

	if err != nil {
		return err
	}

	data, err := utils.BarChart(views)

	if err != nil {
		return err
	}

	return c.Render("partials/graph", fiber.Map{
		"Bars":     data,
		"Animated": true,
	}, "")
}

func PostsStatsHandler(c *fiber.Ctx) error {
	var posts []utils.CountData

	err := data.Dot.Select(data.DB, &posts, "stats-chart-posts-per-month")

	if err != nil {
		return err
	}

	data, err := utils.BarChart(posts)

	if err != nil {
		return err
	}

	return c.Render("partials/graph", fiber.Map{
		"Bars":     data,
		"Animated": true,
	}, "")
}
