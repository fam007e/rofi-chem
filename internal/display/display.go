package display

import (
	"fmt"
	"strings"

	"rofi-chem/internal/config"
)

type Formatter struct {
	cfg *config.Config
}

func NewFormatter(cfg *config.Config) *Formatter {
	return &Formatter{cfg: cfg}
}

func (f *Formatter) FormatElement(data map[string]interface{}) string {
	nameColor := f.cfg.Colors.ElementName


	symbol := getString(data, "symbol")
	name := getString(data, "name")

	mainText := fmt.Sprintf("<span color='%s'><b>%s (%s)</b></span>", nameColor, name, symbol)

	return mainText + "\x00info\x1felement:" + symbol
}

func (f *Formatter) FormatCompound(data map[string]interface{}) string {
	nameColor := f.cfg.Colors.CompoundName


	name := getString(data, "name")
	formula := getString(data, "formula")

	mainText := fmt.Sprintf("<span color='%s'><b>%s (%s)</b></span>", nameColor, name, formula)

	return mainText + "\x00info\x1fcompound:" + name
}

func formatLabel(field string) string {
	switch field {
	case "atomic_number":
		return "Z"
	case "atomic_mass":
		return "A"
	case "melting_point":
		return "Mp"
	case "boiling_point":
		return "Bp"
	case "molecular_weight":
		return "MW"
	case "electron_configuration":
		return "EC"
	case "electronegativity":
		return "EN"
	}
	return strings.Title(strings.ReplaceAll(field, "_", " "))
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func (f *Formatter) FormatDetailList(data map[string]interface{}) []string {
	var lines []string

	// Map of units for properties
	units := map[string]string{
		"atomic_mass":      " u",
		"atomic_weight":    " u",
		"melting_point":    " K",
		"boiling_point":    " K",
		"density":          " g/cm³",
		"atomic_radius":    " pm",
		"covalent_radius":  " pm",
		"vdw_radius":       " pm",
		"electron_affinity": " eV",
	}

	name := getString(data, "name")
	if name != "" {
		lines = append(lines, fmt.Sprintf("Name: %s", name))
	}

	seen := make(map[string]bool)
	seen["name"] = true

	// Helper to process fields
	processField := func(field string, val string) {
		// Sanitize: skip Python object representations
		if strings.HasPrefix(val, "<") || strings.Contains(val, "bound method") {
			return
		}

		// Add unit if available and value is numeric-ish
		if unit, ok := units[field]; ok {
			val += unit
		}

		lines = append(lines, fmt.Sprintf("%s: %s", formatLabel(field), val))
		seen[field] = true
	}

	for _, field := range f.cfg.ElementFields {
		if !seen[field] {
			val := getString(data, field)
			if val != "" && val != "<nil>" {
				processField(field, val)
			}
		}
	}

	for _, field := range f.cfg.CompoundFields {
		if !seen[field] {
			val := getString(data, field)
			if val != "" && val != "<nil>" {
				processField(field, val)
			}
		}
	}

	for k, v := range data {
		if !seen[k] {
			val := fmt.Sprintf("%v", v)
			if val != "" && val != "<nil>" {
				processField(k, val)
			}
		}
	}

	return lines
}
