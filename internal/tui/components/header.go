package components

import (
	"github.com/cfpperche/vibescaffold/internal/tui/styles"
)

func Header() string {
	logo := styles.Logo.Render(styles.LogoASCII)
	version := styles.Subtle.Render("  vibescaffold v0.1.0")
	return logo + version + "\n"
}
