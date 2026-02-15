package routes

import (
	"strconv"
	"time"

	"github.com/believer/willcodefor-go/data"
	"github.com/believer/willcodefor-go/model"
	"github.com/believer/willcodefor-go/utils"
	"github.com/gofiber/fiber/v2"
)

func BooksHandler(c *fiber.Ctx) error {
	var books []model.Book
	var currentBooks []model.Book
	var nextBooks []model.Book

	currentYear := time.Now().Format("2006")
	year := c.Query("year", currentYear)

	err := data.Dot.Select(data.DB, &books, "get-books", year)

	if err != nil {
		return err
	}

	err = data.Dot.Select(data.DB, &currentBooks, "currently-reading")

	if err != nil {
		return err
	}

	err = data.Dot.Select(data.DB, &nextBooks, "next-books")

	if err != nil {
		return err
	}

	totalWords := 0

	years := make([]string, 0)
	startYear := 2024
	endYear := time.Now().Year()

	for year := startYear; year <= endYear; year++ {
		y := strconv.Itoa(year)
		years = append(years, y)
	}

	for _, book := range books {
		totalWords += book.WordCount
	}

	for _, book := range currentBooks {
		totalWords += (book.WordCount / book.PageCount) * book.CurrentPage
	}

	booksRead := len(books)
	now := time.Now()
	dayOfYear := now.YearDay()

	if year != currentYear {
		dayOfYear = 365
	}

	wordsPerDay := totalWords / dayOfYear

	formattedTotalWords := utils.FormatNumber(totalWords)
	formattedWordsPerDay := utils.FormatNumber(wordsPerDay)

	yearlyProgress := float64(booksRead) / 20 * 100

	return c.Render("books", fiber.Map{
		"Path":                 "/books",
		"Books":                books,
		"HasCurrentBooks":      len(currentBooks) > 0,
		"HasPreviousBooks":     len(books) > 0,
		"CurrentBooks":         currentBooks,
		"NextBooks":            nextBooks,
		"HasNextBooks":         len(nextBooks) > 0,
		"FormattedTotalWords":  formattedTotalWords,
		"FormattedWordsPerDay": formattedWordsPerDay,
		"BooksRead":            booksRead,
		"YearlyProgress":       yearlyProgress,
		"SelectedYear":         year,
		"Years":                years,
	})
}
