package footerwidget

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
