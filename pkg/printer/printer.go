package printer

import "github.com/pterm/pterm"

type Printer interface {
	HighlightText(s string) string
	BigTextWithColor(colorStr, text string) (string, error)
}

type ConsolePrinter struct {
}

func (p *ConsolePrinter) HighlightText(s string) string {
	return highlightFgColor.Sprint(s)
}

func (p *ConsolePrinter) BigTextWithColor(colorStr, text string) (string, error) {
	var color pterm.Color
	switch colorStr {
	case green:
		color = greenFgColor
	case yellow:
		color = yellowFgColor
	case red:
		color = redFgColor
	default:
		color = defaultFgColor

	}

	return pterm.DefaultBigText.WithLetters(pterm.NewLettersFromStringWithStyle(text, pterm.NewStyle(color))).Srender()

}
