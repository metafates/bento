package main

import (
	"fmt"
	"log"
	"time"

	"github.com/metafates/bento"
	"github.com/metafates/bento/blockwidget"
	"github.com/metafates/bento/examples/demo/gradient"
	"github.com/metafates/bento/examples/demo/tabs"
	"github.com/metafates/bento/examples/demo/theme"
	"github.com/metafates/bento/tabswidget"
	"github.com/metafates/bento/textwidget"
	"github.com/muesli/termenv"
)

var _ bento.Model = (*Model)(nil)

type Tab int

const (
	TabAbout Tab = iota
	TabRecipe
	TabEmail
	TabTraceroute
	TabWeather
	_tabsCount
)

func (t Tab) String() string {
	switch t {
	case TabAbout:
		return "About"
	case TabEmail:
		return "Email"
	case TabRecipe:
		return "Recipe"
	case TabTraceroute:
		return "Traceroute"
	case TabWeather:
		return "Weather"
	default:
		return ""
	}
}

func (t Tab) Title() string {
	if t == TabAbout {
		return ""
	}

	return " " + t.String() + " "
}

func (t Tab) Next() Tab {
	return min(_tabsCount-1, t+1)
}

func (t Tab) Prev() Tab {
	return max(0, t-1)
}

type Model struct {
	frameCount int
	destroy    bool
	tab        Tab

	aboutTab tabs.AboutTab

	recipeTabState tabs.RecipeTabState
	recipeTab      tabs.RecipeTab
}

// Init implements bento.Model.
func (m *Model) Init() bento.Cmd {
	return nil
}

// Render implements bento.Model.
func (m *Model) Render(area bento.Rect, buffer *bento.Buffer) {
	m.frameCount++

	var titleBar, tab, bottomBar bento.Rect

	bento.
		NewLayout(
			bento.ConstraintLen(1),
			bento.ConstraintMin(0),
			bento.ConstraintLen(1),
		).
		Vertical().
		Split(area).
		Assign(&titleBar, &tab, &bottomBar)

	blockwidget.New().WithStyle(theme.Global.Root).Render(area, buffer)

	m.renderTitleBar(titleBar, buffer)
	m.renderSelectedTab(tab, buffer)
	m.renderBottomBar(bottomBar, buffer)

	if m.destroy {
		destroy(m.frameCount, buffer)
	}
}

func (m *Model) renderTitleBar(area bento.Rect, buffer *bento.Buffer) {
	var title, tabs bento.Rect

	bento.
		NewLayout(
			bento.ConstraintMin(0),
			bento.ConstraintLen(43),
		).
		Horizontal().
		Split(area).
		Assign(&title, &tabs)

	textwidget.NewSpan("Bento").WithStyle(theme.Global.AppTitle).Render(title, buffer)

	titles := make([]textwidget.Line, 0, _tabsCount)
	for _, tab := range []Tab{
		TabAbout,
		TabRecipe,
		TabEmail,
		TabTraceroute,
		TabWeather,
	} {
		titles = append(titles, textwidget.NewLineStr(tab.Title()))
	}

	tabswidget.
		New(titles...).
		WithStyle(theme.Global.Tabs).
		WithHighlightStyle(theme.Global.TabsSelected).
		Select(int(m.tab)).
		WithDividerStr("").
		WithPaddingLeftStr("").
		WithPaddingRightStr("").
		Render(tabs, buffer)
}

func (m *Model) renderSelectedTab(area bento.Rect, buffer *bento.Buffer) {
	gradient.New().Render(area, buffer)

	switch m.tab {
	case TabAbout:
		m.aboutTab.Render(area, buffer)
	case TabRecipe:
		m.recipeTab.RenderStateful(area, buffer, m.recipeTabState)
	}
}

func (m *Model) renderBottomBar(area bento.Rect, buffer *bento.Buffer) {
	var spans []textwidget.Span

	for _, tuple := range [][]string{
		{"H/←", "Left"},
		{"L/→", "Right"},
		{"K/↑", "Up"},
		{"J/↓", "Down"},
		{"D/Del", "Destroy"},
		{"Q/Esc", "Quit"},
	} {
		key, desc := tuple[0], tuple[1]

		keySpan := textwidget.
			NewSpan(fmt.Sprintf(" %s ", key)).
			WithStyle(theme.Global.KeyBinding.Key)

		descSpan := textwidget.
			NewSpan(fmt.Sprintf(" %s ", desc)).
			WithStyle(theme.Global.KeyBinding.Description)

		spans = append(spans, keySpan, descSpan)
	}

	fg := termenv.ANSI256.Color("236")
	bg := termenv.ANSI256.Color("232")

	style := bento.NewStyle().WithForeground(fg).WithBackground(bg)

	textwidget.NewLine(spans...).Center().WithStyle(style).Render(area, buffer)
}

type TickMsg time.Time

func (m *Model) tick() bento.Cmd {
	return bento.Tick(20*time.Millisecond, func(t time.Time) bento.Msg {
		return TickMsg(t)
	})
}

// Update implements bento.Model.
func (m *Model) Update(msg bento.Msg) (bento.Model, bento.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		return m, m.tick()
	case bento.KeyMsg:
		switch msg.String() {
		case "d":
			if m.destroy {
				return m, nil
			}

			m.destroy = true
			return m, m.tick()
		case "ctrl+c", "esc", "q":
			return m, bento.Quit

		case "h":
			m.prevTab()
			return m, nil

		case "l":
			m.nextTab()
			return m, nil
		}
	}

	return m, nil
}

func (m *Model) prevTab() {
	m.tab = m.tab.Prev()
}

func (m *Model) nextTab() {
	m.tab = m.tab.Next()
}

func newModel() Model {
	return Model{
		tab:            TabAbout,
		aboutTab:       tabs.NewAboutTab(),
		recipeTabState: tabs.NewRecipeTabState(),
		recipeTab:      tabs.NewRecipeTab(tabs.SalmonNigiriRecipe),
	}
}

func run() error {
	model := newModel()

	_, err := bento.NewApp(&model).Run()
	if err != nil {
		return fmt.Errorf("new app: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
