package models

// Theme represents a color scheme for the map poster
type Theme struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Background    string `json:"bg"`
	Text          string `json:"text"`
	GradientColor string `json:"gradient_color"`
	Water         string `json:"water"`
	Parks         string `json:"parks"`
	RoadMotorway  string `json:"road_motorway"`
	RoadPrimary   string `json:"road_primary"`
	RoadSecondary string `json:"road_secondary"`
	RoadTertiary  string `json:"road_tertiary"`
	RoadResidential string `json:"road_residential"`
	RoadDefault   string `json:"road_default"`
}

// GenerateRequest represents a request to generate a map poster
type GenerateRequest struct {
	City     string `json:"city" binding:"required"`
	Country  string `json:"country" binding:"required"`
	Theme    string `json:"theme"`
	Distance int    `json:"distance"`
}

// GenerateResponse represents the response after generating a poster
type GenerateResponse struct {
	Success  bool   `json:"success"`
	Filename string `json:"filename,omitempty"`
	URL      string `json:"url,omitempty"`
	Message  string `json:"message,omitempty"`
}

// PreviewRequest represents a request to generate a preview
type PreviewRequest struct {
	City     string `json:"city" binding:"required"`
	Country  string `json:"country" binding:"required"`
	Theme    string `json:"theme"`
	Distance int    `json:"distance"`
}

// Location represents geographic coordinates
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

// ThemeListResponse represents a list of available themes
type ThemeListResponse struct {
	Themes []string `json:"themes"`
}

// ThemeResponse represents a single theme
type ThemeResponse struct {
	Theme Theme `json:"theme"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}
