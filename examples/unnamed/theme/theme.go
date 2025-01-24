package theme

import (
	"github.com/metafates/bento"
	"github.com/muesli/termenv"
)

var (
	white     = termenv.RGBColor("#FFFFFF")
	darkBlue  = termenv.RGBColor("#101830")
	lightBlue = termenv.RGBColor("#4060c0")
	darkGray  = termenv.RGBColor("#444444")
	midGray   = termenv.RGBColor("#808080")
	lightGray = termenv.RGBColor("#bcbcbc")
)

var Global = Theme{
	Root:             bento.NewStyle().WithBackground(darkBlue),
	Tabs:             bento.NewStyle().WithForeground(midGray).WithBackground(darkBlue),
	TabsSelected:     bento.NewStyle().WithForeground(white).WithBackground(darkBlue).Bold().Reversed(),
	AppTitle:         bento.NewStyle().WithForeground(white).WithBackground(darkBlue).Bold(),
	Borders:          bento.NewStyle().WithForeground(lightGray),
	Description:      bento.NewStyle().WithForeground(lightGray).WithBackground(darkBlue),
	DescriptionTitle: bento.NewStyle().WithForeground(lightGray).Bold(),
	Content:          bento.NewStyle().WithForeground(lightGray).WithBackground(darkBlue),
}

type Theme struct {
	Root             bento.Style
	Content          bento.Style
	AppTitle         bento.Style
	Tabs             bento.Style
	TabsSelected     bento.Style
	Borders          bento.Style
	Description      bento.Style
	DescriptionTitle bento.Style
	KeyBinding       ThemeKeyBinding
}

type ThemeKeyBinding struct {
	Key         bento.Style
	Description bento.Style
}

type ThemeEmail struct {
	Tabs         bento.Style
	TabsSelected bento.Style
	Inbox        bento.Style
	Item         bento.Style
	SelectedItem bento.Style
	Header       bento.Style
	HeaderValue  bento.Style
	Body         bento.Style
}

type ThemeTraceroute struct {
	Header   bento.Style
	Selected bento.Style
	Ping     bento.Style
	Map      ThemeMap
}

type ThemeMap struct {
	Style           bento.Style
	Color           bento.Color
	Path            bento.Color
	Source          bento.Color
	Destination     bento.Color
	BackgroundColor bento.Color
}

type ThemeRecipe struct {
	Ingredients       bento.Style
	IngredientsHeader bento.Style
}
