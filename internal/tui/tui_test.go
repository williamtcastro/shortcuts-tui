package tui

import (
	"strings"
	"testing"
)

func TestItemMethods(t *testing.T) {
	item := Item{
		ItemTitle:   "gs",
		ItemDesc:    "Git Status",
		ItemContent: "Full command content here",
		Category:    "Git",
		ViewName:    "Aliases",
		IsAlias:     true,
		Command:     "git status",
	}

	// Test Title
	if item.Title() != "gs" {
		t.Errorf("Expected Title() 'gs', got %s", item.Title())
	}

	// Test Description with Category
	expectedDesc := "[Git] Git Status"
	if item.Description() != expectedDesc {
		t.Errorf("Expected Description() '%s', got %s", expectedDesc, item.Description())
	}

	// Test FilterValue (Deep Search)
	filter := item.FilterValue()
	if !strings.Contains(filter, "gs") || !strings.Contains(filter, "Git Status") || !strings.Contains(filter, "git status") {
		t.Error("FilterValue should contain Title, Desc, and Command for deep search")
	}
}
