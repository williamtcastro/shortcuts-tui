package tui

import (
	"strings"
	"testing"

	"github.com/williamtcastro/shortcuts-tui/internal/models"
)

func TestItemMethods(t *testing.T) {
	item := models.Item{
		ItemTitle:   "gs",
		ItemDesc:    "Git Status",
		ItemContent: "Full command content here",
		Category:    "Git",
		Subdivision: "Work",
		ViewName:    "Aliases",
		IsAlias:     true,
		Command:     "git status",
	}

	// Test Title
	if item.Title() != "gs" {
		t.Errorf("Expected Title() 'gs', got %s", item.Title())
	}

	// Test Description with Category and Subdivision
	expectedDesc := "[Work > Git] Git Status"
	if item.Description() != expectedDesc {
		t.Errorf("Expected Description() '%s', got %s", expectedDesc, item.Description())
	}

	// Test FilterValue (Deep Search)
	filter := item.FilterValue()
	if !strings.Contains(filter, "gs") || !strings.Contains(filter, "Git Status") || !strings.Contains(filter, "git status") || !strings.Contains(filter, "Work") {
		t.Error("FilterValue should contain Title, Desc, Category, Subdivision, and Command for deep search")
	}

	// Test without Subdivision
	item.Subdivision = ""
	expectedDescNoSub := "[Git] Git Status"
	if item.Description() != expectedDescNoSub {
		t.Errorf("Expected Description() '%s', got %s", expectedDescNoSub, item.Description())
	}

	// Test without Category
	item.Category = ""
	expectedDescNoCat := "Git Status"
	if item.Description() != expectedDescNoCat {
		t.Errorf("Expected Description() '%s', got %s", expectedDescNoCat, item.Description())
	}
}
