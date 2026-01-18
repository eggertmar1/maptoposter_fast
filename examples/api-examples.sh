#!/bin/bash

# Example script demonstrating the Map Poster Fast API

API_BASE="http://localhost:8080"

echo "========================================="
echo "Map Poster Fast - API Examples"
echo "========================================="
echo ""

# 1. Health Check
echo "1. Health Check"
echo "GET $API_BASE/health"
curl -s "$API_BASE/health" | jq .
echo ""
echo ""

# 2. List Available Themes
echo "2. List Available Themes"
echo "GET $API_BASE/api/themes"
curl -s "$API_BASE/api/themes" | jq .
echo ""
echo ""

# 3. Get Theme Details
echo "3. Get Theme Details (noir)"
echo "GET $API_BASE/api/themes/noir"
curl -s "$API_BASE/api/themes/noir" | jq .
echo ""
echo ""

# 4. Get Coordinates
echo "4. Get Coordinates"
echo "GET $API_BASE/api/coordinates?city=Paris&country=France"
curl -s "$API_BASE/api/coordinates?city=Paris&country=France" | jq .
echo ""
echo ""

# 5. Generate Preview
echo "5. Generate Preview"
echo "POST $API_BASE/api/preview"
curl -s -X POST "$API_BASE/api/preview" \
  -H "Content-Type: application/json" \
  -d '{
    "city": "Paris",
    "country": "France",
    "theme": "noir",
    "distance": 5000
  }' | jq .
echo ""
echo ""

# 6. Generate Full Poster
echo "6. Generate Full Poster"
echo "POST $API_BASE/api/generate"
curl -s -X POST "$API_BASE/api/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "city": "Tokyo",
    "country": "Japan",
    "theme": "ocean",
    "distance": 12000
  }' | jq .
echo ""
echo ""

echo "========================================="
echo "Examples complete!"
echo "========================================="
echo ""
echo "You can view generated posters at:"
echo "  $API_BASE/posters/"
echo ""
echo "Or use the web interface at:"
echo "  $API_BASE/"
