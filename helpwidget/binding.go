package helpwidget

import (
	"strings"

	"github.com/metafates/bento"
	"github.com/metafates/bento/filterablelistwidget"
	"github.com/metafates/bento/textwidget"
)

var _ filterablelistwidget.Item = (*Binding)(nil)

var _helpBinding = NewBinding("help", "?").WithDescription("Show help")

type Action func() bento.Cmd

type Binding struct {
	DisplayKey  string
	Key         string
	Aliases     []string
	Name        string
	Description string
	Action      Action
	IsHidden    bool
}

func NewBinding(action, key string) Binding {
	return Binding{
		DisplayKey:  strings.ReplaceAll(key, "ctrl+", "^"),
		Key:         key,
		Name:        action,
		Description: "",
		Action:      nil,
	}
}

func (b Binding) WithAliases(aliases ...string) Binding {
	b.Aliases = aliases
	return b
}

func (b Binding) Call() bento.Cmd {
	if b.Action == nil {
		return nil
	}

	return b.Action()
}

func (b Binding) String() string {
	if b.DisplayKey != "" {
		return b.DisplayKey
	}

	return b.Key
}

func (b Binding) WithDisplayKey(displayKey string) Binding {
	b.DisplayKey = displayKey
	return b
}

func (b Binding) WithAction(action Action) Binding {
	b.Action = action
	return b
}

func (b Binding) WithDescription(description string) Binding {
	b.Description = description
	return b
}

func (b Binding) Hidden() Binding {
	b.IsHidden = true
	return b
}

func (b Binding) Matches(key bento.Key) bool {
	keyStr := key.String()

	if b.Key == keyStr {
		return true
	}

	for _, alias := range b.Aliases {
		if alias == keyStr {
			return true
		}
	}

	return false
}

func (b Binding) Text() textwidget.Text {
	text := textwidget.NewText(
		textwidget.NewLine(
			textwidget.NewSpan(b.String()).WithStyle(bento.NewStyle().Bold()),
			textwidget.NewSpan("  "),
			textwidget.NewSpan(b.Name),
		),
	)

	if b.Description != "" {
		description := textwidget.NewLineStr(b.Description).WithStyle(bento.NewStyle().Italic().Dim())

		text.Lines = append(text.Lines, description)
	}

	return text
}
