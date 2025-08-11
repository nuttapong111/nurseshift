#!/bin/bash

# Microservice Generator Script
# Usage: ./create-microservice.sh service-name port

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <service-name> <port>"
    echo "Example: $0 user-service 8082"
    exit 1
fi

SERVICE_NAME=$1
PORT=$2
PROJECT_ROOT="/Volumes/NO NAME/project/nurseshift"
BACKEND_DIR="$PROJECT_ROOT/backend"
AUTH_SERVICE_DIR="$BACKEND_DIR/auth-service"
NEW_SERVICE_DIR="$BACKEND_DIR/$SERVICE_NAME"

echo "🚀 Creating $SERVICE_NAME on port $PORT..."

# Create directory structure
mkdir -p "$NEW_SERVICE_DIR"/{cmd/server,internal/{domain/{entities,usecases,repositories},infrastructure/{config,database,services},interfaces/http/{handlers,routes,middleware}}}

# Copy base files from auth-service
echo "📂 Copying base structure..."
cp -r "$AUTH_SERVICE_DIR/internal/infrastructure" "$NEW_SERVICE_DIR/internal/"
cp -r "$AUTH_SERVICE_DIR/internal/interfaces" "$NEW_SERVICE_DIR/internal/"

# Create go.mod
echo "📦 Creating go.mod..."
cat > "$NEW_SERVICE_DIR/go.mod" << EOF
module nurseshift/$SERVICE_NAME

go 1.21

require (
	github.com/gofiber/fiber/v2 v2.50.0
	github.com/google/uuid v1.5.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.13.0
)
EOF

# Create config.env
echo "⚙️ Creating config.env..."
cat > "$NEW_SERVICE_DIR/config.env" << EOF
# Server Configuration
PORT=$PORT
ENV=development

# Database Configuration  
DB_HOST=localhost
DB_PORT=5432
DB_USER=nuttapong2
DB_PASSWORD=
DB_NAME=nurseshift
DB_SSLMODE=disable
DB_SCHEMA=nurse_shift

# Auth Service Configuration
AUTH_SERVICE_URL=http://localhost:8081

# Security
BCRYPT_COST=12

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:3002
CORS_CREDENTIALS=true
EOF

echo "✅ $SERVICE_NAME created successfully!"
echo "📍 Location: $NEW_SERVICE_DIR"
echo "🌐 Port: $PORT"


