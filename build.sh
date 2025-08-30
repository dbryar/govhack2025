#!/bin/bash

# Build script for Transliterate service with embedded frontend

set -e

echo "Building Transliterate Service..."

# Check if Hugo is installed
if ! command -v hugo &> /dev/null; then
    echo "Error: Hugo is not installed. Please install Hugo first."
    echo "Visit: https://gohugo.io/getting-started/installing/"
    exit 1
fi

# Build the frontend
echo "Building frontend with Hugo..."
cd frontend
hugo --minify --baseURL /app/
cd ..

# Copy dist folder to service directory for embedding
echo "Copying frontend dist to service directory..."
rm -rf transliterate/dist
cp -r frontend/dist transliterate/dist

# Build the service
echo "Building Encore service..."
encore build

echo "Build complete!"
echo ""
echo "To run locally: encore run"
echo "To deploy: encore deploy"