package data

import "time"

type Project struct {
	Description string
	Link        string
	Name        string
	Tech        []string
}

type Work struct {
	Company         string
	Description     string
	EndDate         string
	Link            string
	LinkDescription string
	StartDate       string
	Title           string
}

var Projects = [8]Project{
	{
		Name:        "My Movies",
		Description: "This is where I keep track of all the movies I've watched. I've been doing it for over 20 years, first using lists in IMDb and from 2011 using my own database.",
		Tech:        []string{"go", "htmx", "tailwind", "hyperscript"},
		Link:        "https://movies.willcodefor.beer",
	},
	{
		Name:        "Supreme",
		Description: "Supreme is a command line tool that helps you get up and running fast with new apps. It can currently generate rescript-react apps with Tailwind CSS, GraphQL APIs with examples for queries, mutations and subscriptions using TypeScript and React apps with both TypeScript and JavaScript. It can also help you install and generate commonly used configs for things like prettier, husky and jest. ",
		Tech:        []string{"rust", "github actions"},
		Link:        "https://github.com/opendevtools/supreme",
	},
	{
		Name:        "rescript-intl",
		Description: "re-intl helps you with date, number and currency formatting in ReasonML (BuckleScript). Everything is built on top of Intl which comes built-in with browsers >= IE11 as well as Node.",
		Tech:        []string{"rescript", "github actions"},
		Link:        "https://github.com/opendevtools/rescript-intl",
	},
	{
		Name:        "Clearingnummer",
		Description: "Sort codes, clearingnummer in Swedish, are four or five digit identifiers for Swedish banks. This package helps you find the bank related to a specific number. ",
		Tech:        []string{"typescript", "github actions"},
		Link:        "https://github.com/believer/clearingnummer",
	},
	{
		Name:        "Telefonnummer",
		Description: "Telefonnummer is phone number in Swedish. This package formats all Swedish phone numbers, both mobile and landline, to a standard format. ",
		Tech:        []string{"typescript", "github actions"},
		Link:        "https://github.com/believer/telefonnummer",
	},
	{
		Name:        "WCAG Color",
		Description: "<p>According to the WHO an <a href=\"https://www.who.int/en/news-room/fact-sheets/detail/blindness-and-visual-impairment\">estimated 1.3 billion</a> people live with some form of visual impairment. This includes people who are legally blind and people with less than 20/20 vision.</p>  <p>This library helps you achieve the accessibility standards for color contrast outlined in the WCAG 2.0 specification.</p> ",
		Tech:        []string{"rescript", "github actions"},
		Link:        "https://github.com/opendevtools/wcag-color",
	},
	{
		Name:        "Wejay",
		Description: "A Slack bot that controls a Sonos system. We use it at Iteam as a collaborative music player. It can do pretty much everything from managing the play queue, control playback, list most played songs and even contains some hidden easter eggs. ",
		Tech:        []string{"reasonml", "docker", "elasticsearch", "github actions", "slack"},
		Link:        "https://github.com/Iteam1337/sonos-wejay",
	},
	{
		Name:        "Workout of the Day",
		Description: "<p>A collection of competition and benchmark CrossFit workouts but also workouts that I\"ve made. A combination of two of my passions code and CrossFit.</p><p>I\"ve also made a version of the app in <a href=\"https://github.com/believer/wod-elm\">Elm</a>.</p> ",
		Tech:        []string{"rescript", "vercel", "github actions"},
		Link:        "https://github.com/believer/wod",
	},
}

var Positions = [5]Work{
	{
		Company:     "SEB",
		StartDate:   "2023-06-05T00:00:00.000Z",
		EndDate:     "",
		Title:       "Senior Fullstack Developer",
		Description: "<p>SEB is one of the largest Swedish banks.</p>",
	},
	{
		Company:     "Arizon",
		StartDate:   "2022-01-10T00:00:00.000Z",
		EndDate:     "2023-05-30T00:00:00.000Z",
		Title:       "Developer Consultant",
		Description: "<p>Arizon is a IT consultancy and startup incubator.</p>",
	},
	{
		Company:     "Hemnet",
		StartDate:   "2020-04-20T00:00:00.000Z",
		EndDate:     "2022-12-01T00:00:00.000Z",
		Title:       "Frontend Developer",
		Description: "<p>With 2.8 miljon unique vistors each week, Hemnet is Sweden\"s biggest website when you\"re looking to buy or sell your appartment or house.</p><p>I was part of the Seller's Experience team. This team handles the \"behind the scenes\" of a sale. Everything from the broker adding your listing, you purchasing additional packages for better exposure of your listing to a dashboard where you can follow statistics on the sale.</p>",
	},
	{
		Company:         "Iteam",
		StartDate:       "2012-11-05T00:00:00.000Z",
		EndDate:         "2020-03-01T00:00:00.000Z",
		Title:           "Developer / Head of Tech",
		Description:     "<p>Iteam is a development consultancy working mostly in-house.</p><p>My work focused on front-end, but also backend (Node) whenever there's a need. We use React and React Native with TypeScript, but recently we've also started using ReasonML. We write all code using TDD and Jest. API integrations are made using GraphQL, with some REST.</p>",
		Link:            "/iteam",
		LinkDescription: "Here's a list of all the projects I've a been a part of at Iteam",
	},
	{
		Company:     "MatHem",
		StartDate:   "2011-12-01T00:00:00.000Z",
		EndDate:     "2012-06-01T00:00:00.000Z",
		Title:       "Interaction designer",
		Description: "<p>MatHem delivers groceries directly to your door, either as a prepackaged concept with recipes or as individual products of your choosing. MatHem has been selected as one of the best Swedish online stores two years running by Internetworld.</p><p>My job was mostly front-end development. I made mockups in Photoshop and then implemented the HTML, CSS and some jQuery on the website. I also made flash banners for advertising campaigns.</p>",
	},
}

func init() {
	// For each position parse the dates and update the data
	for i, position := range Positions {
		startDate, _ := time.Parse(time.RFC3339, position.StartDate)
		Positions[i].StartDate = startDate.Format("2006")

		if position.EndDate != "" {
			endDate, _ := time.Parse(time.RFC3339, position.EndDate)
			Positions[i].EndDate = endDate.Format("06")
		}
	}
}
