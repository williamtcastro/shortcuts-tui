package models

import (
	"fmt"
	"strings"
)

type Item struct {
	ItemTitle   string
	ItemDesc    string
	ItemContent string
	Category    string
	Subdivision string
	ViewName    string
	IsAlias     bool
	Command     string
}

func (i Item) Title() string { return i.ItemTitle }

func (i Item) Description() string {
	prefix := ""
	if i.Subdivision != "" {
		prefix = i.Subdivision
	}
	if i.Category != "" {
		category := strings.Title(strings.ReplaceAll(i.Category, "_", " "))
		if prefix != "" {
			prefix += " > " + category
		} else {
			prefix = category
		}
	}

	if prefix != "" {
		return fmt.Sprintf("[%s] %s", prefix, i.ItemDesc)
	}
	return i.ItemDesc
}

func (i Item) FilterValue() string {
	// Search by Title, Description, Category, Subdivision, AND the full command/content
	return i.ItemTitle + " " + i.ItemDesc + " " + i.Category + " " + i.Subdivision + " " + i.Command + " " + i.ItemContent
}
