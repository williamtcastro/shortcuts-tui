package tui

type Item struct {
	ItemTitle   string
	ItemDesc    string
	ItemContent string
	IsAlias     bool
	Command     string
}

func (i Item) Title() string       { return i.ItemTitle }
func (i Item) Description() string { return i.ItemDesc }
func (i Item) FilterValue() string { return i.ItemTitle + " " + i.ItemDesc }
