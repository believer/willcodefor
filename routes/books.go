package routes

import (
	"time"

	"github.com/believer/willcodefor-go/data"
	"github.com/believer/willcodefor-go/model"
	"github.com/gofiber/fiber/v2"
)

func BooksHandler(c *fiber.Ctx) error {
	var books []model.Book
	var currentBook model.Book

	err := data.Dot.Select(data.DB, &books, "get-books")

	if err != nil {
		return err
	}

	err = data.Dot.Get(data.DB, &currentBook, "currently-reading")

	if err != nil {
		return err
	}

	totalWords := 0

	for _, book := range books {
		totalWords += book.WordCount
	}

	totalWords += (currentBook.WordCount / currentBook.PageCount) * currentBook.CurrentPage

	totalBooks := len(books) + 1

	now := time.Now()
	dayOfYear := now.YearDay()
	wordsPerDay := totalWords / dayOfYear

	return c.Render("books", fiber.Map{
		"Path":        "/books",
		"Books":       books,
		"CurrentBook": currentBook,
		"TotalWords":  totalWords,
		"TotalBooks":  totalBooks,
		"WordsPerDay": wordsPerDay,
	})
}
