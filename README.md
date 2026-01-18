# Map Poster Fast 🚀

A blazing-fast map poster generation service built in Go. Generate beautiful, minimalist map posters for any city in the world through a simple REST API.

## Features

✨ **High Performance**: Built in Go for speed and efficiency  
🌍 **Global Coverage**: Works with any city using OpenStreetMap data  
🎨 **Multiple Themes**: Choose from various color schemes  
🔌 **REST API**: Easy to integrate with any frontend  
📦 **Docker Ready**: Simple deployment with Docker  
🖼️ **Preview Support**: Generate previews before full-resolution images  

## Quick Start

### Using Docker (Recommended)

```bash
# Build and run
docker-compose up -d

# The service will be available at http://localhost:8080
```

### Manual Installation

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go
```

## API Documentation

### Base URL
```
http://localhost:8080
```

### Endpoints

#### 1. Health Check
Check if the service is running.

```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0"
}
```

#### 2. List Available Themes
Get a list of all available themes.

```http
GET /api/themes
```

**Response:**
```json
{
  "themes": ["feature_based", "noir", "ocean"]
}
```

#### 3. Get Theme Details
Get the color scheme for a specific theme.

```http
GET /api/themes/:name
```

**Response:**
```json
{
  "theme": {
    "name": "Noir",
    "description": "Classic black and white high contrast",
    "bg": "#000000",
    "text": "#FFFFFF",
    "gradient_color": "#000000",
    "water": "#1A1A1A",
    "parks": "#0F0F0F",
    "road_motorway": "#FFFFFF",
    "road_primary": "#E0E0E0",
    "road_secondary": "#C0C0C0",
    "road_tertiary": "#A0A0A0",
    "road_residential": "#808080",
    "road_default": "#808080"
  }
}
```

#### 4. Get Coordinates
Get coordinates for a city.

```http
GET /api/coordinates?city=Paris&country=France
```

**Response:**
```json
{
  "latitude": 48.8566,
  "longitude": 2.3522,
  "address": "Paris, Île-de-France, France"
}
```

#### 5. Generate Poster
Generate a full-resolution map poster.

```http
POST /api/generate
Content-Type: application/json

{
  "city": "Paris",
  "country": "France",
  "theme": "noir",
  "distance": 10000
}
```

**Parameters:**
- `city` (required): City name
- `country` (required): Country name
- `theme` (optional): Theme name (default: "feature_based")
- `distance` (optional): Map radius in meters (default: 29000)

**Response:**
```json
{
  "success": true,
  "filename": "paris_noir_20260118_120000.png",
  "url": "/posters/paris_noir_20260118_120000.png",
  "message": "Poster generated successfully"
}
```

#### 6. Generate Preview
Generate a preview of the map poster (smaller area for faster generation).

```http
POST /api/preview
Content-Type: application/json

{
  "city": "Tokyo",
  "country": "Japan",
  "theme": "ocean",
  "distance": 5000
}
```

**Response:**
```json
{
  "success": true,
  "filename": "tokyo_ocean_20260118_120000.png",
  "url": "/posters/tokyo_ocean_20260118_120000.png",
  "message": "Preview generated successfully"
}
```

#### 7. Download Generated Posters
Access generated posters via the `/posters` endpoint.

```http
GET /posters/:filename
```

## Distance Guidelines

- **4000-6000m**: Small/dense cities (Venice, Amsterdam center)
- **8000-12000m**: Medium cities, focused downtown (Paris, Barcelona)
- **15000-20000m**: Large metros, full city view (Tokyo, Mumbai)
- **25000-35000m**: Metropolitan areas

## Available Themes

1. **feature_based** - Different shades for road hierarchy (default)
2. **noir** - Classic black and white high contrast
3. **ocean** - Deep ocean blue with lighter road network

### Adding Custom Themes

Create a JSON file in the `themes/` directory:

```json
{
  "name": "My Theme",
  "description": "My custom color scheme",
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

## Example Frontend Integration

### JavaScript/Fetch

```javascript
async function generatePoster() {
  const response = await fetch('http://localhost:8080/api/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      city: 'Paris',
      country: 'France',
      theme: 'noir',
      distance: 10000
    })
  });
  
  const data = await response.json();
  if (data.success) {
    console.log('Poster URL:', data.url);
    // Display the image
    document.getElementById('poster').src = 'http://localhost:8080' + data.url;
  }
}
```

### cURL

```bash
# Generate a poster
curl -X POST http://localhost:8080/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "city": "New York",
    "country": "USA",
    "theme": "noir",
    "distance": 12000
  }'

# List themes
curl http://localhost:8080/api/themes

# Get coordinates
curl "http://localhost:8080/api/coordinates?city=London&country=UK"
```

## Development

### Project Structure

```
.
├── cmd/
│   └── server/          # Main application entry point
│       └── main.go
├── pkg/
│   ├── api/            # HTTP handlers and routes
│   │   └── handlers.go
│   ├── generator/      # Map generation logic
│   │   └── generator.go
│   ├── models/         # Data models
│   │   └── models.go
│   └── themes/         # Theme management
│       └── themes.go
├── themes/             # Theme JSON files
├── fonts/              # Font files (optional)
├── posters/            # Generated posters (output)
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

### Building

```bash
# Build the binary
go build -o server ./cmd/server

# Run
./server

# Build with Docker
docker build -t maptoposter-fast .
```

### Environment Variables

- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin mode (debug/release)

## Performance

- **Fast Generation**: Go's concurrency and efficient OSM data handling
- **Low Memory**: Optimized image processing
- **Scalable**: Easy to deploy multiple instances behind a load balancer

## Credits

Based on the original [maptoposter](https://github.com/originalankur/maptoposter) by originalankur, reimplemented in Go for performance and service deployment.

## License

MIT License - Feel free to use in your projects!

## Support

For issues, questions, or contributions, please open an issue on GitHub.
