package utils

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

func FormatNumber(n int) string {
	p := message.NewPrinter(language.Swedish)

	return p.Sprintf("%v", number.Decimal(n))
}
