package themes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eggertmar1/maptoposter_fast/pkg/models"
)

const ThemesDir = "themes"

// LoadTheme loads a theme from a JSON file
func LoadTheme(themeName string) (*models.Theme, error) {
	themeFile := filepath.Join(ThemesDir, themeName+".json")

	// Check if theme file exists
	if _, err := os.Stat(themeFile); os.IsNotExist(err) {
		// Return default theme if file doesn't exist
		return getDefaultTheme(), nil
	}

	// Read theme file
	data, err := os.ReadFile(themeFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file: %w", err)
	}

	// Parse theme
	var theme models.Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("failed to parse theme: %w", err)
	}

	return &theme, nil
}

// GetAvailableThemes returns a list of available theme names
func GetAvailableThemes() ([]string, error) {
	// Check if themes directory exists
	if _, err := os.Stat(ThemesDir); os.IsNotExist(err) {
		if err := os.MkdirAll(ThemesDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create themes directory: %w", err)
		}
		return []string{}, nil
	}

	// Read directory
	entries, err := os.ReadDir(ThemesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read themes directory: %w", err)
	}

	// Filter for .json files
	themes := []string{}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			themeName := strings.TrimSuffix(entry.Name(), ".json")
			themes = append(themes, themeName)
		}
	}

	return themes, nil
}

// getDefaultTheme returns a default theme
func getDefaultTheme() *models.Theme {
	return &models.Theme{
		Name:            "Feature-Based Shading",
		Description:     "Different shades for different road types and features with clear hierarchy",
		Background:      "#FFFFFF",
		Text:            "#000000",
		GradientColor:   "#FFFFFF",
		Water:           "#C0C0C0",
		Parks:           "#F0F0F0",
		RoadMotorway:    "#0A0A0A",
		RoadPrimary:     "#1A1A1A",
		RoadSecondary:   "#2A2A2A",
		RoadTertiary:    "#3A3A3A",
		RoadResidential: "#4A4A4A",
		RoadDefault:     "#3A3A3A",
	}
}

// SaveTheme saves a theme to a JSON file
func SaveTheme(themeName string, theme *models.Theme) error {
	// Ensure themes directory exists
	if err := os.MkdirAll(ThemesDir, 0755); err != nil {
		return fmt.Errorf("failed to create themes directory: %w", err)
	}

	themeFile := filepath.Join(ThemesDir, themeName+".json")

	// Marshal theme to JSON
	data, err := json.MarshalIndent(theme, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal theme: %w", err)
	}

	// Write to file
	if err := os.WriteFile(themeFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write theme file: %w", err)
	}

	return nil
}
