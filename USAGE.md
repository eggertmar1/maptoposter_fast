# Usage Guide

This guide shows you how to use the Map Poster Fast service in different ways.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Using the Web Interface](#using-the-web-interface)
3. [Using the REST API](#using-the-rest-api)
4. [Integrating with Your Application](#integrating-with-your-application)
5. [Adding Custom Themes](#adding-custom-themes)

## Quick Start

### Start the Service

**Option 1: Using Docker (Recommended)**
```bash
docker-compose up -d
```

**Option 2: Using Go**
```bash
go run cmd/server/main.go
```

**Option 3: Using the Binary**
```bash
./server
```

The service will be available at `http://localhost:8080`

## Using the Web Interface

1. Open your browser and navigate to `http://localhost:8080`
2. Enter a city name (e.g., "Paris")
3. Enter a country name (e.g., "France")
4. Select a theme from the dropdown
5. Adjust the distance (radius in meters) if needed
6. Click "Generate Preview" for a quick preview or "Generate Poster" for full resolution

The generated poster will appear below the form, and you can download it using the download link.

## Using the REST API

### 1. Check Service Health

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "version": "1.0.0"
}
```

### 2. List Available Themes

```bash
curl http://localhost:8080/api/themes
```

Response:
```json
{
  "themes": ["feature_based", "noir", "ocean"]
}
```

### 3. Get Theme Details

```bash
curl http://localhost:8080/api/themes/noir
```

Response:
```json
{
  "theme": {
    "name": "Noir",
    "description": "Classic black and white high contrast",
    "bg": "#000000",
    "text": "#FFFFFF",
    ...
  }
}
```

### 4. Get Coordinates for a Location

```bash
curl "http://localhost:8080/api/coordinates?city=Paris&country=France"
```

Response:
```json
{
  "latitude": 48.8566,
  "longitude": 2.3522,
  "address": "Paris, Île-de-France, France"
}
```

### 5. Generate a Preview

```bash
curl -X POST http://localhost:8080/api/preview \
  -H "Content-Type: application/json" \
  -d '{
    "city": "Paris",
    "country": "France",
    "theme": "noir",
    "distance": 5000
  }'
```

Response:
```json
{
  "success": true,
  "filename": "paris_noir_20260118_120000.png",
  "url": "/posters/paris_noir_20260118_120000.png",
  "message": "Preview generated successfully"
}
```

### 6. Generate a Full Poster

```bash
curl -X POST http://localhost:8080/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "city": "New York",
    "country": "USA",
    "theme": "feature_based",
    "distance": 12000
  }'
```

Response:
```json
{
  "success": true,
  "filename": "new_york_feature_based_20260118_120000.png",
  "url": "/posters/new_york_feature_based_20260118_120000.png",
  "message": "Poster generated successfully"
}
```

### 7. Download Generated Poster

```bash
wget http://localhost:8080/posters/paris_noir_20260118_120000.png
```

Or open in browser:
```
http://localhost:8080/posters/paris_noir_20260118_120000.png
```

## Integrating with Your Application

### JavaScript/React Example

```javascript
async function generateMapPoster(city, country, theme = 'noir', distance = 10000) {
  try {
    const response = await fetch('http://localhost:8080/api/generate', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        city,
        country,
        theme,
        distance
      })
    });

    const data = await response.json();
    
    if (data.success) {
      console.log('Poster generated:', data.url);
      return `http://localhost:8080${data.url}`;
    } else {
      throw new Error(data.message);
    }
  } catch (error) {
    console.error('Failed to generate poster:', error);
    throw error;
  }
}

// Usage
generateMapPoster('Paris', 'France', 'noir', 10000)
  .then(url => {
    document.getElementById('poster').src = url;
  })
  .catch(error => {
    console.error('Error:', error);
  });
```

### Python Example

```python
import requests
import json

def generate_map_poster(city, country, theme='noir', distance=10000):
    url = 'http://localhost:8080/api/generate'
    payload = {
        'city': city,
        'country': country,
        'theme': theme,
        'distance': distance
    }
    
    response = requests.post(url, json=payload)
    data = response.json()
    
    if data['success']:
        print(f"Poster generated: {data['filename']}")
        return data['url']
    else:
        raise Exception(data['message'])

# Usage
try:
    poster_url = generate_map_poster('Tokyo', 'Japan', 'ocean', 15000)
    print(f"Download from: http://localhost:8080{poster_url}")
except Exception as e:
    print(f"Error: {e}")
```

### Go Example

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
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

func generatePoster(city, country, theme string, distance int) (string, error) {
    req := GenerateRequest{
        City:     city,
        Country:  country,
        Theme:    theme,
        Distance: distance,
    }

    jsonData, err := json.Marshal(req)
    if err != nil {
        return "", err
    }

    resp, err := http.Post(
        "http://localhost:8080/api/generate",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var result GenerateResponse
    if err := json.Unmarshal(body, &result); err != nil {
        return "", err
    }

    if !result.Success {
        return "", fmt.Errorf(result.Message)
    }

    return result.URL, nil
}

func main() {
    url, err := generatePoster("Paris", "France", "noir", 10000)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Poster available at: http://localhost:8080%s\n", url)
}
```

## Adding Custom Themes

You can create your own color schemes by adding JSON files to the `themes/` directory.

### Create a New Theme

1. Create a new file: `themes/my_theme.json`
2. Add the color scheme:

```json
{
  "name": "My Custom Theme",
  "description": "A beautiful custom color scheme",
  "bg": "#FFFFFF",
  "text": "#000000",
  "gradient_color": "#FFFFFF",
  "water": "#C0C0C0",
  "parks": "#F0F0F0",
  "road_motorway": "#0A0A0A",
  "road_primary": "#1A1A1A",
  "road_secondary": "#2A2A2A",
  "road_tertiary": "#3A3A3A",
  "road_residential": "#4A4A4A",
  "road_default": "#3A3A3A"
}
```

3. The theme will be automatically available as "my_theme"

### Color Scheme Guide

- `bg`: Background color
- `text`: Text color for city/country labels
- `gradient_color`: Color for gradient fades (usually same as bg)
- `water`: Color for water bodies (rivers, lakes, ocean)
- `parks`: Color for parks and green spaces
- `road_motorway`: Color for highways/motorways (thickest roads)
- `road_primary`: Color for primary roads
- `road_secondary`: Color for secondary roads
- `road_tertiary`: Color for tertiary roads
- `road_residential`: Color for residential streets
- `road_default`: Default color for other roads

## Distance Guidelines

Choose the appropriate distance based on the city size and area you want to cover:

| Distance | Best For | Examples |
|----------|----------|----------|
| 4000-6000m | Small/dense cities, historic centers | Venice, Amsterdam center |
| 8000-12000m | Medium cities, focused downtown | Paris, Barcelona |
| 15000-20000m | Large metros, full city view | Tokyo, Mumbai |
| 25000-35000m | Metropolitan areas | Greater London, NYC metro |

## Tips and Best Practices

1. **Start with a preview**: Always generate a preview first to check the composition
2. **Adjust distance**: If the map looks too sparse or too dense, adjust the distance parameter
3. **Theme selection**: Choose themes that complement your interior design
4. **High contrast themes**: Work best for printing (noir, feature_based)
5. **API rate limiting**: Be respectful of OSM services; don't make too many rapid requests

## Troubleshooting

### Issue: "Location not found"
- Check spelling of city and country names
- Try different name variants (e.g., "NYC" vs "New York")
- Make sure the location exists in OpenStreetMap

### Issue: "Generation takes too long"
- Reduce the distance parameter
- Try generating a preview first
- Check your internet connection

### Issue: "Empty or sparse map"
- Increase the distance parameter
- Try a different location
- Check if the area has good OSM coverage

## Support

For more help, check:
- Main README.md for setup instructions
- GitHub issues for known problems
- API documentation in README.md

Enjoy creating beautiful map posters! 🗺️
