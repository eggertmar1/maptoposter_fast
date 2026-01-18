package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eggertmar1/maptoposter_fast/pkg/models"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

const (
	PostersDir = "posters"
	FontsDir   = "fonts"
	Width      = 2400  // 12 inches * 200 DPI
	Height     = 3200  // 16 inches * 200 DPI
)

// Generator handles map poster generation
type Generator struct {
	httpClient *http.Client
}

// NewGenerator creates a new generator instance
func NewGenerator() *Generator {
	return &Generator{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetCoordinates fetches coordinates for a given city and country using Nominatim
func (g *Generator) GetCoordinates(city, country string) (*models.Location, error) {
	// Use Nominatim geocoding API
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Add("q", fmt.Sprintf("%s, %s", city, country))
	params.Add("format", "json")
	params.Add("limit", "1")

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set user agent as required by Nominatim
	req.Header.Set("User-Agent", "MapPosterFast/1.0")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch coordinates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nominatim returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var results []struct {
		Lat         string `json:"lat"`
		Lon         string `json:"lon"`
		DisplayName string `json:"display_name"`
	}

	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("location not found for %s, %s", city, country)
	}

	// Parse coordinates
	var lat, lon float64
	fmt.Sscanf(results[0].Lat, "%f", &lat)
	fmt.Sscanf(results[0].Lon, "%f", &lon)

	return &models.Location{
		Latitude:  lat,
		Longitude: lon,
		Address:   results[0].DisplayName,
	}, nil
}

// GeneratePoster generates a map poster for the given parameters
func (g *Generator) GeneratePoster(city, country string, theme *models.Theme, distance int) (string, error) {
	// Ensure posters directory exists
	if err := os.MkdirAll(PostersDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create posters directory: %w", err)
	}

	// Get coordinates
	location, err := g.GetCoordinates(city, country)
	if err != nil {
		return "", fmt.Errorf("failed to get coordinates: %w", err)
	}

	// Fetch OSM data
	osmData, err := g.fetchOSMData(location.Latitude, location.Longitude, distance)
	if err != nil {
		return "", fmt.Errorf("failed to fetch OSM data: %w", err)
	}

	// Generate image
	img, err := g.renderMap(osmData, theme, city, country, location)
	if err != nil {
		return "", fmt.Errorf("failed to render map: %w", err)
	}

	// Generate filename
	timestamp := time.Now().Format("20060102_150405")
	citySlug := strings.ToLower(strings.ReplaceAll(city, " ", "_"))
	themeSlug := strings.ToLower(strings.ReplaceAll(theme.Name, " ", "_"))
	filename := fmt.Sprintf("%s_%s_%s.png", citySlug, themeSlug, timestamp)
	filepath := filepath.Join(PostersDir, filename)

	// Save image
	f, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	return filename, nil
}

// OSMData represents OpenStreetMap data
type OSMData struct {
	Elements []OSMElement `json:"elements"`
}

// OSMElement represents a single OSM element
type OSMElement struct {
	Type      string             `json:"type"`
	ID        int64              `json:"id"`
	Lat       float64            `json:"lat,omitempty"`
	Lon       float64            `json:"lon,omitempty"`
	Nodes     []int64            `json:"nodes,omitempty"`
	Members   []OSMMember        `json:"members,omitempty"`
	Tags      map[string]string  `json:"tags,omitempty"`
	Geometry  []OSMNode          `json:"geometry,omitempty"`
}

// OSMMember represents a member of a relation
type OSMMember struct {
	Type string `json:"type"`
	Ref  int64  `json:"ref"`
	Role string `json:"role"`
}

// OSMNode represents a node with coordinates
type OSMNode struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// fetchOSMData fetches map data from Overpass API
func (g *Generator) fetchOSMData(lat, lon float64, distance int) (*OSMData, error) {
	// Calculate bounding box
	// Approximate: 1 degree latitude ≈ 111km
	distKm := float64(distance) / 1000.0
	latDelta := distKm / 111.0
	lonDelta := distKm / (111.0 * math.Cos(lat*math.Pi/180.0))

	south := lat - latDelta
	north := lat + latDelta
	west := lon - lonDelta
	east := lon + lonDelta

	// Build Overpass query
	query := fmt.Sprintf(`[out:json][timeout:25];
(
  way["highway"](%f,%f,%f,%f);
  way["natural"="water"](%f,%f,%f,%f);
  way["waterway"](%f,%f,%f,%f);
  way["leisure"="park"](%f,%f,%f,%f);
  way["landuse"="grass"](%f,%f,%f,%f);
);
out geom;`, 
		south, west, north, east,
		south, west, north, east,
		south, west, north, east,
		south, west, north, east,
		south, west, north, east)

	// Use Overpass API
	apiURL := "https://overpass-api.de/api/interpreter"
	
	resp, err := g.httpClient.Post(apiURL, "application/x-www-form-urlencoded", 
		bytes.NewBufferString("data="+url.QueryEscape(query)))
	if err != nil {
		return nil, fmt.Errorf("failed to query Overpass API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("overpass API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var osmData OSMData
	if err := json.Unmarshal(body, &osmData); err != nil {
		return nil, fmt.Errorf("failed to parse OSM data: %w", err)
	}

	return &osmData, nil
}

// renderMap renders the map to an image
func (g *Generator) renderMap(osmData *OSMData, theme *models.Theme, city, country string, location *models.Location) (image.Image, error) {
	// Create image
	img := image.NewRGBA(image.Rect(0, 0, Width, Height))

	// Parse background color
	bgColor, err := parseHexColor(theme.Background)
	if err != nil {
		return nil, fmt.Errorf("invalid background color: %w", err)
	}

	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Calculate bounds for all elements
	minLat, maxLat, minLon, maxLon := calculateBounds(osmData)
	if minLat == maxLat || minLon == maxLon {
		return nil, fmt.Errorf("invalid bounds: no map data available")
	}

	// Render features with layering
	g.renderFeatures(img, osmData, theme, minLat, maxLat, minLon, maxLon)

	// Add text overlays
	if err := g.addTextOverlay(img, city, country, location, theme); err != nil {
		return nil, fmt.Errorf("failed to add text overlay: %w", err)
	}

	return img, nil
}

// calculateBounds calculates the bounding box for all OSM elements
func calculateBounds(osmData *OSMData) (minLat, maxLat, minLon, maxLon float64) {
	minLat, minLon = 90.0, 180.0
	maxLat, maxLon = -90.0, -180.0

	for _, element := range osmData.Elements {
		if element.Lat != 0 && element.Lon != 0 {
			if element.Lat < minLat {
				minLat = element.Lat
			}
			if element.Lat > maxLat {
				maxLat = element.Lat
			}
			if element.Lon < minLon {
				minLon = element.Lon
			}
			if element.Lon > maxLon {
				maxLon = element.Lon
			}
		}
		for _, node := range element.Geometry {
			if node.Lat < minLat {
				minLat = node.Lat
			}
			if node.Lat > maxLat {
				maxLat = node.Lat
			}
			if node.Lon < minLon {
				minLon = node.Lon
			}
			if node.Lon > maxLon {
				maxLon = node.Lon
			}
		}
	}

	return
}

// latLonToPixel converts lat/lon coordinates to pixel coordinates
func latLonToPixel(lat, lon, minLat, maxLat, minLon, maxLon float64, width, height int) (int, int) {
	// Add margin
	margin := 100
	availWidth := width - 2*margin
	availHeight := height - 400 // Reserve space for text

	// Normalize coordinates
	x := (lon - minLon) / (maxLon - minLon)
	y := (lat - minLat) / (maxLat - minLat)

	// Convert to pixel coordinates (invert Y axis)
	px := margin + int(x*float64(availWidth))
	py := margin + int((1-y)*float64(availHeight))

	return px, py
}

// renderFeatures renders all map features
func (g *Generator) renderFeatures(img *image.RGBA, osmData *OSMData, theme *models.Theme, minLat, maxLat, minLon, maxLon float64) {
	// Render in layers: water -> parks -> roads

	// Layer 1: Water features
	waterColor, _ := parseHexColor(theme.Water)
	for _, element := range osmData.Elements {
		if isWater(element.Tags) {
			g.renderPolygon(img, element.Geometry, waterColor, minLat, maxLat, minLon, maxLon)
		}
	}

	// Layer 2: Parks
	parksColor, _ := parseHexColor(theme.Parks)
	for _, element := range osmData.Elements {
		if isPark(element.Tags) {
			g.renderPolygon(img, element.Geometry, parksColor, minLat, maxLat, minLon, maxLon)
		}
	}

	// Layer 3: Roads
	for _, element := range osmData.Elements {
		if isRoad(element.Tags) {
			roadColor, width := g.getRoadStyle(element.Tags, theme)
			g.renderWay(img, element.Geometry, roadColor, width, minLat, maxLat, minLon, maxLon)
		}
	}
}

// isWater checks if an element is a water feature
func isWater(tags map[string]string) bool {
	return tags["natural"] == "water" || tags["waterway"] != ""
}

// isPark checks if an element is a park
func isPark(tags map[string]string) bool {
	return tags["leisure"] == "park" || tags["landuse"] == "grass"
}

// isRoad checks if an element is a road
func isRoad(tags map[string]string) bool {
	return tags["highway"] != ""
}

// getRoadStyle returns the color and width for a road based on its type
func (g *Generator) getRoadStyle(tags map[string]string, theme *models.Theme) (color.Color, int) {
	highway := tags["highway"]

	var colorHex string
	var width int

	switch highway {
	case "motorway", "motorway_link":
		colorHex = theme.RoadMotorway
		width = 4
	case "trunk", "trunk_link", "primary", "primary_link":
		colorHex = theme.RoadPrimary
		width = 3
	case "secondary", "secondary_link":
		colorHex = theme.RoadSecondary
		width = 2
	case "tertiary", "tertiary_link":
		colorHex = theme.RoadTertiary
		width = 2
	case "residential", "living_street", "unclassified":
		colorHex = theme.RoadResidential
		width = 1
	default:
		colorHex = theme.RoadDefault
		width = 1
	}

	c, _ := parseHexColor(colorHex)
	return c, width
}

// renderWay renders a way (line) on the image
func (g *Generator) renderWay(img *image.RGBA, geometry []OSMNode, color color.Color, width int, minLat, maxLat, minLon, maxLon float64) {
	if len(geometry) < 2 {
		return
	}

	for i := 0; i < len(geometry)-1; i++ {
		x1, y1 := latLonToPixel(geometry[i].Lat, geometry[i].Lon, minLat, maxLat, minLon, maxLon, Width, Height)
		x2, y2 := latLonToPixel(geometry[i+1].Lat, geometry[i+1].Lon, minLat, maxLat, minLon, maxLon, Width, Height)

		drawThickLine(img, x1, y1, x2, y2, color, width)
	}
}

// renderPolygon renders a filled polygon on the image
func (g *Generator) renderPolygon(img *image.RGBA, geometry []OSMNode, color color.Color, minLat, maxLat, minLon, maxLon float64) {
	if len(geometry) < 3 {
		return
	}

	// Simple polygon fill using scanline algorithm
	points := make([]image.Point, len(geometry))
	for i, node := range geometry {
		x, y := latLonToPixel(node.Lat, node.Lon, minLat, maxLat, minLon, maxLon, Width, Height)
		points[i] = image.Point{x, y}
	}

	// Draw filled polygon (simplified - just draw lines for now)
	for i := 0; i < len(points)-1; i++ {
		drawThickLine(img, points[i].X, points[i].Y, points[i+1].X, points[i+1].Y, color, 1)
	}
	// Close the polygon
	if len(points) > 0 {
		drawThickLine(img, points[len(points)-1].X, points[len(points)-1].Y, points[0].X, points[0].Y, color, 1)
	}
}

// drawThickLine draws a line with specified thickness using Bresenham's algorithm
func drawThickLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color, thickness int) {
	// Simple line drawing - can be improved with anti-aliasing
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := -1
	if x1 < x2 {
		sx = 1
	}
	sy := -1
	if y1 < y2 {
		sy = 1
	}
	err := dx - dy

	for {
		// Draw thick point
		for tx := -thickness / 2; tx <= thickness/2; tx++ {
			for ty := -thickness / 2; ty <= thickness/2; ty++ {
				px, py := x1+tx, y1+ty
				if px >= 0 && px < Width && py >= 0 && py < Height {
					img.Set(px, py, c)
				}
			}
		}

		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// addTextOverlay adds city name, country, and coordinates to the image
func (g *Generator) addTextOverlay(img *image.RGBA, city, country string, location *models.Location, theme *models.Theme) error {
	// For now, we'll use a simple text rendering
	// In production, you'd want to use a proper font rendering library

	textColor, err := parseHexColor(theme.Text)
	if err != nil {
		return err
	}

	// Try to load font
	fontPath := filepath.Join(FontsDir, "Roboto-Bold.ttf")
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		// Font not available, skip text rendering for now
		return nil
	}

	f, err := truetype.Parse(fontData)
	if err != nil {
		return nil
	}

	// Draw city name
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(100)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(&image.Uniform{textColor})

	// City name (spaced)
	spacedCity := strings.Join(strings.Split(strings.ToUpper(city), ""), "  ")
	pt := freetype.Pt(Width/2-len(spacedCity)*20, Height-500)
	c.DrawString(spacedCity, pt)

	// Country
	c.SetFontSize(40)
	pt = freetype.Pt(Width/2-len(country)*10, Height-400)
	c.DrawString(strings.ToUpper(country), pt)

	// Coordinates
	c.SetFontSize(24)
	coords := fmt.Sprintf("%.4f° N / %.4f° E", location.Latitude, location.Longitude)
	if location.Latitude < 0 {
		coords = fmt.Sprintf("%.4f° S / %.4f° E", math.Abs(location.Latitude), location.Longitude)
	}
	if location.Longitude < 0 {
		coords = strings.Replace(coords, "E", "W", 1)
	}
	pt = freetype.Pt(Width/2-len(coords)*5, Height-300)
	c.DrawString(coords, pt)

	return nil
}

// parseHexColor converts a hex color string to color.RGBA
func parseHexColor(s string) (color.RGBA, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid color format: %s", s)
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{r, g, b, 255}, nil
}
