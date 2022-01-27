package printer

import (
	"github.com/pterm/pterm"
)

func (p *ConsolePrinter) BigLettersWithColor(colorStr, text string) (string, error) {
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
