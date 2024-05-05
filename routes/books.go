package routes

import (
	"time"

	"github.com/believer/willcodefor-go/data"
	"github.com/believer/willcodefor-go/model"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

func BooksHandler(c *fiber.Ctx) error {
	var books []model.Book
	var currentBooks []model.Book

	p := message.NewPrinter(language.Swedish)
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

	booksRead := len(books)
	now := time.Now()
	dayOfYear := now.YearDay()
	wordsPerDay := totalWords / dayOfYear

	formattedTotalWords := p.Sprintf("%v", number.Decimal(totalWords))
	formattedWordsPerDay := p.Sprintf("%v", number.Decimal(wordsPerDay))

	return c.Render("books", fiber.Map{
		"Path":                 "/books",
		"Books":                books,
		"CurrentBooks":         currentBooks,
		"FormattedTotalWords":  formattedTotalWords,
		"FormattedWordsPerDay": formattedWordsPerDay,
		"TotalWords":           totalWords,
		"BooksRead":            booksRead,
		"WordsPerDay":          wordsPerDay,
	})
}
