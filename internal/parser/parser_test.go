package parser

import (
	"reflect"
	"testing"
)

func TestAliasRegex(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected []string
	}{
		{
			name:     "Basic alias",
			line:     `alias ll="ls -la"`,
			expected: []string{`alias ll="ls -la"`, "ll", "ls -la", ""},
		},
		{
			name:     "Alias with comment",
			line:     `alias gs="git status" # Show git status`,
			expected: []string{`alias gs="git status" # Show git status`, "gs", "git status", "Show git status"},
		},
		{
			name:     "Alias with trailing spaces",
			line:     `alias p="pnpm"   # Pnpm manager  `,
			expected: []string{`alias p="pnpm"   # Pnpm manager  `, "p", "pnpm", "Pnpm manager  "},
		},
		{
			name:     "Alias with complex command",
			line:     `alias brewup="brew update && brew upgrade" # Update everything`,
			expected: []string{`alias brewup="brew update && brew upgrade" # Update everything`, "brewup", "brew update && brew upgrade", "Update everything"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := aliasRegex.FindStringSubmatch(tt.line)
			if !reflect.DeepEqual(matches, tt.expected) {
				t.Errorf("FindStringSubmatch() = %v, want %v", matches, tt.expected)
			}
		})
	}
}
