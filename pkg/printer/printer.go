package printer

type Printer interface {
	HighlightText(s string) string
	BigLettersWithColor(colorStr, text string) (string, error)
}

type ConsolePrinter struct {
}

func (p *ConsolePrinter) HighlightText(s string) string {
	return highlightFgColor.Sprint(s)
}
