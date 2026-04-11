package tui

import "fmt"

type Item struct {
	ItemTitle   string
	ItemDesc    string
	ItemContent string
	Category    string
	IsAlias     bool
	Command     string
}

func (i Item) Title() string { return i.ItemTitle }

func (i Item) Description() string {
	if i.Category != "" {
		return fmt.Sprintf("[%s] %s", i.Category, i.ItemDesc)
	}
	return i.ItemDesc
}

func (i Item) FilterValue() string { return i.ItemTitle + " " + i.ItemDesc + " " + i.Category }
