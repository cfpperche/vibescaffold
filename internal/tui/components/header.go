package components

import (
	"github.com/cfpperche/vibeforge/internal/i18n"
	"github.com/cfpperche/vibeforge/internal/tui/styles"
)

func Header() string {
	logo := styles.Logo.Render(styles.LogoASCII())
	version := styles.Subtle.Render("  " + i18n.T("app.version"))
	return "\n " + logo + version + "\n\n"
}
