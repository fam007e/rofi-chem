package search

import (
	"sort"
	"strings"

	"rofi-chem/internal/db"
)

type Result struct {
	Type  string
	Data  map[string]interface{}
	Score int
}

func PerformSearch(d *db.Database, query string, threshold int) ([]Result, error) {
	var results []Result

	// Get all elements for fuzzy searching
	elements, err := d.GetAllElements()
	if err != nil {
		return nil, err
	}

	for _, e := range elements {
		score := 0
		symbol, _ := e["symbol"].(string)
		name, _ := e["name"].(string)


		// Show all elements if query is empty
		if query == "" {
			score = 100
		} else if strings.EqualFold(query, symbol) {
			// Exact match bonus
			score = 100
		} else if strings.EqualFold(query, name) {
			score = 100
		} else {
			// Simple containment check for now (faster than full Levenshtein for all)
			// For a true fuzzy search, we'll implement Levenshtein below
			dist := levenshtein(strings.ToLower(query), strings.ToLower(name))
			maxLen := max(len(query), len(name))
			if maxLen > 0 {
				ratio := 100 - (dist * 100 / maxLen)
				if ratio >= threshold {
					score = ratio
				}
			}
		}

		if score > 0 {
			results = append(results, Result{
				Type:  "element",
				Data:  e,
				Score: score,
			})
		}
	}

	// Compounds (Using SQL LIKE for now as list can be huge)
	compounds, err := d.SearchCompounds(query)
	if err != nil {
		return nil, err
	}
	for _, c := range compounds {
		results = append(results, Result{
			Type:  "compound",
			Data:  c,
			Score: 90, // Inherently relevant if returned by DB
		})
	}

	// Sort by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

// Simple Levenshtein distance implementation
func levenshtein(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	n, m := len(r1), len(r2)
	if n > m {
		r1, r2 = r2, r1
		n, m = m, n
	}

	currentRow := make([]int, n+1)
	for i := 0; i <= n; i++ {
		currentRow[i] = i
	}

	for i := 1; i <= m; i++ {
		previousRow := currentRow
		currentRow = make([]int, n+1)
		currentRow[0] = i

		for j := 1; j <= n; j++ {
			add, del, change := previousRow[j]+1, currentRow[j-1]+1, previousRow[j-1]
			if r1[j-1] != r2[i-1] {
				change++
			}
			currentRow[j] = min(add, min(del, change))
		}
	}
	return currentRow[n]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
