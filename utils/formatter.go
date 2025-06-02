package utils

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

func FormatNumber(n int) string {
	p := message.NewPrinter(language.English)

	return p.Sprintf("%v", number.Decimal(n))
}

func FormatFloat(n float64, noDecimals bool) string {
	p := message.NewPrinter(language.English)

	if noDecimals {
		return p.Sprintf("%v", n)
	}

	return p.Sprintf("%.2f", n)
}
