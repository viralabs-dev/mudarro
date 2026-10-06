package terminal

import "strings"

// Presentation only affects UI; action identifiers and execution are unchanged.
type Presentation struct{ Locale, Theme, Density, Lettering, Mouse string }

func (v *Menu) Configure(p Presentation) {
	v.locale = p.Locale
	v.style.theme = p.Theme
	v.density = p.Density
	v.lettering = p.Lettering
	v.mouse = p.Mouse
}
func (v *Menu) text(en, pt string) string {
	if v.locale == "pt-BR" {
		return pt
	}
	return en
}
func (v *Menu) mark() []string {
	if v.lettering == "text" || v.density == "compact" {
		return []string{" " + v.style.heading(fitText(v.project, v.style.width-2))}
	}
	if v.lettering == "ascii" {
		return splitBanner(asciiWidth(v.project, v.style.width))
	}
	return wordmark(v.project, v.style)
}
func splitBanner(s string) []string { return strings.Split(strings.TrimSuffix(s, "\n"), "\n") }

func (v *Menu) SetPreviews(previews []string) { v.previews = previews }
