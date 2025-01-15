package footerwidget

import (
	"github.com/metafates/bento"
	"github.com/metafates/bento/listwidget"
	"github.com/metafates/bento/textwidget"
)

var _ listwidget.Item = (*Binding)(nil)

type Binding struct {
	Key         string
	Action      string
	Description string
}

func NewBinding(key, action string) Binding {
	return Binding{
		Key:         key,
		Action:      action,
		Description: "",
	}
}

func (b Binding) WithDescription(description string) Binding {
	b.Description = description
	return b
}

func (b Binding) Title() textwidget.Text {
	text := textwidget.NewText(
		textwidget.NewLine(
			textwidget.NewSpan(b.Key).WithStyle(bento.NewStyle().Bold()),
			textwidget.NewSpan("  "),
			textwidget.NewSpan(b.Action),
		),
	)

	if b.Description != "" {
		description := textwidget.NewLineStr(b.Description).WithStyle(bento.NewStyle().Italic().Dim())

		text.Lines = append(text.Lines, description)
	}

	return text
}
