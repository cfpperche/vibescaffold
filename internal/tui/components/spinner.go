package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/cfpperche/vibeforge/internal/tui/styles"
)

func NewSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.Success
	return s
}
