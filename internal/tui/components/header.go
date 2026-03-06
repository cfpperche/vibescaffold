package components

import (
	"github.com/cfpperche/vibescaffold/internal/tui/styles"
)

func Header() string {
	logo := styles.Logo.Render(styles.LogoASCII)
	version := styles.Subtle.Render("  v0.1.0")
	return "\n " + logo + version + "\n\n"
}
