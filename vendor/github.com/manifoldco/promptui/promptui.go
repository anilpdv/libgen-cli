package promptui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type SelectTemplates struct {
	Active   string
	Inactive string
	Selected string
}

type Key struct {
	Code    rune
	Display string
}

type SelectKeys struct {
	Next     Key
	Prev     Key
	PageUp   Key
	PageDown Key
}

type Select struct {
	Label     string
	Items     interface{}
	Templates *SelectTemplates
	Size      int
	IsVimMode bool
	Keys      *SelectKeys
}

func (s *Select) Run() (int, string, error) {
	var items []string
	switch v := s.Items.(type) {
	case []string:
		items = v
	default:
		return 0, "", fmt.Errorf("unsupported items type")
	}

	if len(items) == 0 {
		return 0, "", fmt.Errorf("no items to select")
	}

	fmt.Printf("\n%s (Total: %d):\n", s.Label, len(items))
	for i, item := range items {
		fmt.Printf("  [%d] %s\n", i+1, item)
	}

	fmt.Printf("\nEnter choice [1-%d] (default 1, 'n' next page, 'p' prev page, 'q' to quit): ", len(items))

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			return 0, items[0], nil
		}
		if text == "q" || text == "quit" || text == "exit" {
			return 0, "", fmt.Errorf("selection cancelled by user")
		}
		if strings.EqualFold(text, "n") || strings.EqualFold(text, "next") {
			for idx, item := range items {
				if strings.Contains(item, "Next Page") {
					return idx, item, nil
				}
			}
		}
		if strings.EqualFold(text, "p") || strings.EqualFold(text, "prev") || strings.EqualFold(text, "previous") {
			for idx, item := range items {
				if strings.Contains(item, "Previous Page") {
					return idx, item, nil
				}
			}
		}
		num, err := strconv.Atoi(text)
		if err == nil && num >= 1 && num <= len(items) {
			return num - 1, items[num-1], nil
		}
		fmt.Println("Invalid selection, defaulting to option 1.")
		return 0, items[0], nil
	}

	// Non-interactive or EOF: default to 1st item
	return 0, items[0], nil
}
