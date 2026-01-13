package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"rofi-chem/internal/config"
	"rofi-chem/internal/db"
	"rofi-chem/internal/display"
	"rofi-chem/internal/search"
)

const BackOption = "⬅ Back to Menu"

func main() {
	// Debug logging
	f, _ := os.OpenFile("/tmp/rofi-chem-debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if f != nil {
		fmt.Fprintf(f, "Called with args: %v\n", os.Args)
		f.Close()
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		// Fallback to defaults if loading fails is handled inside LoadConfig mostly,
		// but if we get a fatal error here, just stderr it for debugging.
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		// Proceed with potentially partial config
	}

	// Connect to database
	database, err := db.NewDatabase()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Initialize components
	formatter := display.NewFormatter(cfg)

	// Check CLI arguments (selection made)
	if len(os.Args) > 1 {
		selection := os.Args[1]

		// Handle Back Option: just fall through to main search
		if selection == BackOption {
			// Do nothing here, allowing code to proceed to main search below
		} else {


		// Regex to parse "<b>Name (Symbol)</b>" format from the markup
		// e.g. <span ...><b>Hydrogen (H)</b></span>
		re := regexp.MustCompile(`<b>(.*?) \((.+?)\)</b>`)
		matches := re.FindStringSubmatch(selection)

		if len(matches) == 3 {
			name := matches[1]
			symbolOrFormula := matches[2]

			// Try finding element by Symbol
			if data, err := database.GetElementBySymbol(symbolOrFormula); err == nil {
				lines := formatter.FormatDetailList(data)
				fmt.Println(BackOption)
				for _, l := range lines {
					fmt.Println(l)
				}
				os.Exit(0)
			}

			// Try finding compound by Name
			if data, err := database.GetCompoundByName(name); err == nil {
				lines := formatter.FormatDetailList(data)
				fmt.Println(BackOption)
				for _, l := range lines {
					fmt.Println(l)
				}
				os.Exit(0)
			}
		}

		// Fallback for direct "type:id" format (if info worked magically for some users)
		if strings.Contains(selection, "\x1f") || strings.Contains(selection, "element:") || strings.Contains(selection, "compound:") {
			parts := strings.Split(selection, ":")
			if len(parts) >= 2 {
				itemType := parts[0]
				itemId := strings.Join(parts[1:], ":")
				// ... existing logic simplified ...
				// Assuming if above regex didn't catch it, and this is "element:H", it's clean
				var data map[string]interface{}
				var err error
				if strings.Contains(itemType, "element") {
					data, err = database.GetElementBySymbol(itemId)
				} else if strings.Contains(itemType, "compound") {
					data, err = database.GetCompoundByName(itemId)
				}
				if err == nil {
					lines := formatter.FormatDetailList(data)
					fmt.Println(BackOption)
					for _, l := range lines {
						fmt.Println(l)
					}
					os.Exit(0)
				}
			}
		}

		// Selection is a property (Copy)
		// Ensure it's NOT a markup string (which starts with <) to avoid loop
		if strings.Contains(selection, ":") && !strings.HasPrefix(selection, "<") {
			parts := strings.SplitN(selection, ":", 2)
			if len(parts) == 2 {
				val := strings.TrimSpace(parts[1])
				copyToClipboard(val)
				os.Exit(0)
			}
		}

			os.Exit(0)
		}
	}

	// Check environment variables for Rofi state
	// ROFI_RETV: 0 = initial, 1 = selected, etc.
	// For script mode, we usually just output the list.
	// If we want dynamic filtering based on user typing (for very large datasets),
	// we'd check ROFI_INFO or arguments.
	// But our dataset is small enough (< 200 items usually) to dump all at once
	// unless we implement fuzzy search inside the tool itself triggered by rofi.
	//
	// Current Python implementation dumped everything.
	// We can add a search query argument later if we switch to `rofi -modi "chem:rofi-chem"`
	// where rofi handles filtering.

	// If the user provided a query via some mechanism (like a wrapper script calling this with args)
	// we could filter. But standard rofi script usage is:
	// 1. Run script -> output list
	// 2. User types -> rofi filters list
	// 3. User selects -> script called with selection

	// So we just output everything relevant.

	// Enable Pango markup in Rofi
	fmt.Print("\x00markup-rows\x1ftrue\n")
	// Also can set prompt
	fmt.Print("\x00prompt\x1fChem\n")

	results, err := search.PerformSearch(database, "", 0) // Empty query returns all elements + no compounds usually?
	// Actually, PerformSearch with empty query in our implementation:
	// Elements: returns all (good)
	// Compounds: returns all matching "" -> all compounds (good)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching: %v\n", err)
		os.Exit(1)
	}


	for _, result := range results {
		var line string
		if result.Type == "element" {
			line = formatter.FormatElement(result.Data)
		} else {
			line = formatter.FormatCompound(result.Data)
		}

		// In Rofi, if we want to pass data back invisible to user, we can use \0info or similar features
		// but standard simple rows are fine for now.
		fmt.Println(line)
	}
}

func copyToClipboard(text string) error {
	f, err := os.OpenFile("/tmp/rofi-chem-debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		fmt.Fprintf(f, "Attempting to copy: '%s'\n", text)
	}

	// Try xclip first
	cmd := exec.Command("xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		if f != nil {
			fmt.Fprintf(f, "xclip failed: %v. Trying wl-copy...\n", err)
		}
		// Try wl-copy as fallback
		cmd = exec.Command("wl-copy")
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err != nil {
			if f != nil {
				fmt.Fprintf(f, "wl-copy failed: %v\n", err)
			}
			return err
		}
	}

	if f != nil {
		fmt.Fprintf(f, "Copy successful.\n")
	}

	// Notify user
	notify := exec.Command("notify-send", "Rofi Chem", fmt.Sprintf("Copied '%s' to clipboard", text))
	if err := notify.Run(); err != nil {
		if f != nil {
			fmt.Fprintf(f, "notify-send failed: %v\n", err)
		}
		return err
	}

	return nil
}
