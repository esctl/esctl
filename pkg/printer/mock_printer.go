package printer

import "github.com/stretchr/testify/mock"

type MockPrinter struct {
	mock.Mock
}

func (p *MockPrinter) BigTextWithColor(colorStr, text string) (string, error) {
	args := p.Called(colorStr, text)
	return args.String(0), args.Error(1)
}

func (p *MockPrinter) HighlightText(s string) string {
	args := p.Called(s)
	return args.String(0)
}
