package components

import (
	"github.com/cfpperche/vibeforge/internal/tui/styles"
)

func Footer(hints string) string {
	return styles.Footer.Render(hints)
}
