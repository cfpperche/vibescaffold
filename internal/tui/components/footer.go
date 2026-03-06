package components

import (
	"github.com/cfpperche/vibescaffold/internal/tui/styles"
)

func Footer(hints string) string {
	return styles.Footer.Render(hints)
}
