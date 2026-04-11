package tui

import "fmt"

type Item struct {
	ItemTitle   string
	ItemDesc    string
	ItemContent string
	Category    string
	ViewName    string
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

func (i Item) FilterValue() string { 
	// Search by Title, Description, Category, AND the full command/content
	return i.ItemTitle + " " + i.ItemDesc + " " + i.Category + " " + i.Command + " " + i.ItemContent
}
