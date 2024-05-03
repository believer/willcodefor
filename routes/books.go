package routes

import (
	"time"

	"github.com/believer/willcodefor-go/data"
	"github.com/believer/willcodefor-go/model"
	"github.com/gofiber/fiber/v2"
)

func BooksHandler(c *fiber.Ctx) error {
	var books []model.Book
	var currentBooks []model.Book

	err := data.Dot.Select(data.DB, &books, "get-books")

	if err != nil {
		return err
	}

	err = data.Dot.Select(data.DB, &currentBooks, "currently-reading")

	if err != nil {
		return err
	}

	totalWords := 0

	for _, book := range books {
		totalWords += book.WordCount
	}

	for _, book := range currentBooks {
		totalWords += (book.WordCount / book.PageCount) * book.CurrentPage
	}

	totalBooks := len(books) + 1

	now := time.Now()
	dayOfYear := now.YearDay()
	wordsPerDay := totalWords / dayOfYear

	for _, book := range currentBooks {
		book.Progress = (float64(book.CurrentPage) / float64(book.PageCount)) * 100
	}

	return c.Render("books", fiber.Map{
		"Path":         "/books",
		"Books":        books,
		"CurrentBooks": currentBooks,
		"TotalWords":   totalWords,
		"TotalBooks":   totalBooks,
		"WordsPerDay":  wordsPerDay,
	})
}
