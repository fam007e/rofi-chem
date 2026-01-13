package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Display        DisplayConfig  `yaml:"display"`
	ElementFields  []string       `yaml:"element_fields"`
	CompoundFields []string       `yaml:"compound_fields"`
	Colors         ColorsConfig   `yaml:"colors"`
	Search         SearchConfig   `yaml:"search"`
}

type DisplayConfig struct {
	MaxResults  int  `yaml:"max_results"`
	UseColors   bool `yaml:"use_colors"`
	ShowIcons   bool `yaml:"show_icons"`
	CompactMode bool `yaml:"compact_mode"`
}

type ColorsConfig struct {
	ElementName   string `yaml:"element_name"`
	CompoundName  string `yaml:"compound_name"`
	PropertyLabel string `yaml:"property_label"`
	PropertyValue string `yaml:"property_value"`
	Separator     string `yaml:"separator"`
}

type SearchConfig struct {
	FuzzyThreshold int  `yaml:"fuzzy_threshold"`
	EnableFuzzy    bool `yaml:"enable_fuzzy"`
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".config", "rofi", "rofi-chem", "config.yaml")

	// Default configuration
	cfg := &Config{
		Display: DisplayConfig{
			MaxResults: 50,
			UseColors:  true,
		},
		ElementFields: []string{
			"symbol", "name", "atomic_number", "atomic_mass", "density",
		},
		CompoundFields: []string{
			"formula", "name", "molecular_weight",
		},
		Colors: ColorsConfig{
			ElementName:   "#61AFEF",
			CompoundName:  "#C678DD",
			PropertyLabel: "#98C379",
			PropertyValue: "#E5C07B",
			Separator:     "#5C6370",
		},
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// If file doesn't exist, return defaults
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
