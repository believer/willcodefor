package data

import (
	"os"
	"time"

	"github.com/lib/pq"
)

// Views before this come from the frozen post_view table and views after it
// from Umami.
var UmamiStart = time.Date(2026, 11, 9, 16, 0, 0, 0, time.UTC)

type umamiPathViews struct {
	Path  string `db:"path"`
	Views int64  `db:"views"`
}

// Page views per post path from Umami, as arrays for the view count queries.
// Empty when Umami isn't configured.
func UmamiPostViews() (pq.StringArray, pq.Int64Array, error) {
	var (
		rows      []umamiPathViews
		paths     = pq.StringArray{}
		views     = pq.Int64Array{}
		websiteID = os.Getenv("UMAMI_WEBSITE_ID")
	)

	if UmamiDB == nil || websiteID == "" {
		return paths, views, nil
	}

	if err := Dot.Select(UmamiDB, &rows, "umami-post-views", websiteID, UmamiStart); err != nil {
		return nil, nil, err
	}

	for _, row := range rows {
		paths = append(paths, row.Path)
		views = append(views, row.Views)
	}

	return paths, views, nil
}
