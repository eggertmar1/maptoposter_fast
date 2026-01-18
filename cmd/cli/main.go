package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	defaultAPIURL = "http://localhost:8080"
	version       = "1.0.0"
)

type GenerateRequest struct {
	City     string `json:"city"`
	Country  string `json:"country"`
	Theme    string `json:"theme"`
	Distance int    `json:"distance"`
}

type GenerateResponse struct {
	Success  bool   `json:"success"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Message  string `json:"message"`
}

type ThemeListResponse struct {
	Themes []string `json:"themes"`
}

func main() {
	// Define flags
	city := flag.String("city", "", "City name (required)")
	country := flag.String("country", "", "Country name (required)")
	theme := flag.String("theme", "feature_based", "Theme name")
	distance := flag.Int("distance", 29000, "Map radius in meters")
	preview := flag.Bool("preview", false, "Generate preview instead of full poster")
	listThemes := flag.Bool("list-themes", false, "List available themes")
	apiURL := flag.String("api", defaultAPIURL, "API base URL")
	showVersion := flag.Bool("version", false, "Show version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Map Poster Fast CLI v%s\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -city Paris -country France\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -city Tokyo -country Japan -theme ocean -distance 15000\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -city London -country UK -preview\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -list-themes\n", os.Args[0])
	}

	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Printf("Map Poster Fast CLI v%s\n", version)
		os.Exit(0)
	}

	// List themes
	if *listThemes {
		listAvailableThemes(*apiURL)
		return
	}

	// Validate required arguments
	if *city == "" || *country == "" {
		fmt.Fprintf(os.Stderr, "Error: -city and -country are required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Generate poster
	endpoint := "/api/generate"
	if *preview {
		endpoint = "/api/preview"
		fmt.Println("Generating preview...")
	} else {
		fmt.Println("Generating poster...")
	}

	req := GenerateRequest{
		City:     *city,
		Country:  *country,
		Theme:    *theme,
		Distance: *distance,
	}

	if err := generatePoster(*apiURL, endpoint, req); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func listAvailableThemes(apiURL string) {
	resp, err := http.Get(apiURL + "/api/themes")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch themes: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "API returned status %d\n", resp.StatusCode)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read response: %v\n", err)
		os.Exit(1)
	}

	var themeList ThemeListResponse
	if err := json.Unmarshal(body, &themeList); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse response: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Available themes:")
	for _, theme := range themeList.Themes {
		fmt.Printf("  - %s\n", theme)
	}
}

func generatePoster(apiURL, endpoint string, req GenerateRequest) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		apiURL+endpoint,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var result GenerateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("%s", result.Message)
	}

	fmt.Println("✓ Success!")
	fmt.Printf("  Filename: %s\n", result.Filename)
	fmt.Printf("  URL: %s%s\n", apiURL, result.URL)
	fmt.Printf("  Download: wget %s%s\n", apiURL, result.URL)

	return nil
}
