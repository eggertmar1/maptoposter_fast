package api

import (
	"fmt"
	"net/http"

	"github.com/eggertmar1/maptoposter_fast/pkg/generator"
	"github.com/eggertmar1/maptoposter_fast/pkg/models"
	"github.com/eggertmar1/maptoposter_fast/pkg/themes"
	"github.com/gin-gonic/gin"
)

const Version = "1.0.0"

const (
	DefaultDistance        = 29000 // Default 29km radius for full poster
	DefaultPreviewDistance = 10000 // Default 10km radius for preview
)

// Handler handles HTTP requests
type Handler struct {
	generator *generator.Generator
}

// NewHandler creates a new API handler
func NewHandler() *Handler {
	return &Handler{
		generator: generator.NewGenerator(),
	}
}

// SetupRoutes sets up all API routes
func (h *Handler) SetupRoutes(r *gin.Engine) {
	// Health check
	r.GET("/health", h.Health)

	// Serve static frontend
	r.Static("/public", "./public")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/public/index.html")
	})

	// API routes
	api := r.Group("/api")
	{
		api.GET("/themes", h.ListThemes)
		api.GET("/themes/:name", h.GetTheme)
		api.POST("/generate", h.Generate)
		api.POST("/preview", h.Preview)
		api.GET("/coordinates", h.GetCoordinates)
	}

	// Serve generated posters
	r.Static("/posters", "./posters")
}

// Health returns the health status of the service
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthResponse{
		Status:  "healthy",
		Version: Version,
	})
}

// ListThemes returns all available themes
func (h *Handler) ListThemes(c *gin.Context) {
	themeList, err := themes.GetAvailableThemes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get themes: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, models.ThemeListResponse{
		Themes: themeList,
	})
}

// GetTheme returns a specific theme by name
func (h *Handler) GetTheme(c *gin.Context) {
	themeName := c.Param("name")

	theme, err := themes.LoadTheme(themeName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Theme not found: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, models.ThemeResponse{
		Theme: *theme,
	})
}

// Generate generates a full-resolution map poster
func (h *Handler) Generate(c *gin.Context) {
	var req models.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	// Set defaults
	if req.Theme == "" {
		req.Theme = "feature_based"
	}
	if req.Distance == 0 {
		req.Distance = DefaultDistance
	}

	// Load theme
	theme, err := themes.LoadTheme(req.Theme)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to load theme: %v", err),
		})
		return
	}

	// Generate poster
	filename, err := h.generator.GeneratePoster(req.City, req.Country, theme, req.Distance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.GenerateResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to generate poster: %v", err),
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, models.GenerateResponse{
		Success:  true,
		Filename: filename,
		URL:      fmt.Sprintf("/posters/%s", filename),
		Message:  "Poster generated successfully",
	})
}

// Preview generates a preview (smaller resolution) of the map poster
func (h *Handler) Preview(c *gin.Context) {
	var req models.PreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	// Set defaults
	if req.Theme == "" {
		req.Theme = "feature_based"
	}
	if req.Distance == 0 {
		req.Distance = DefaultPreviewDistance
	}

	// Load theme
	theme, err := themes.LoadTheme(req.Theme)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to load theme: %v", err),
		})
		return
	}

	// Generate preview (using same generator for now, could be optimized)
	filename, err := h.generator.GeneratePoster(req.City, req.Country, theme, req.Distance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.GenerateResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to generate preview: %v", err),
		})
		return
	}

	// Return response
	c.JSON(http.StatusOK, models.GenerateResponse{
		Success:  true,
		Filename: filename,
		URL:      fmt.Sprintf("/posters/%s", filename),
		Message:  "Preview generated successfully",
	})
}

// GetCoordinates returns coordinates for a city/country
func (h *Handler) GetCoordinates(c *gin.Context) {
	city := c.Query("city")
	country := c.Query("country")

	if city == "" || country == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "city and country parameters are required",
		})
		return
	}

	location, err := h.generator.GetCoordinates(city, country)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Failed to get coordinates: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, location)
}
