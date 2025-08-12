#!/bin/bash

echo "🚀 Starting NurseShift User Service..."
echo "======================================"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Please run this script from the user-service directory"
    exit 1
fi

# Check if config.env exists
if [ ! -f "config.env" ]; then
    echo "⚠️  config.env not found. Creating from example..."
    if [ -f "config.env.example" ]; then
        cp config.env.example config.env
        echo "✅ config.env created from example. Please edit it with your database credentials."
        echo "   Then run this script again."
        exit 1
    else
        echo "❌ config.env.example not found. Please create config.env manually."
        exit 1
    fi
fi

# Check if database is accessible
echo "🔌 Checking database connection..."
DB_HOST=$(grep DB_HOST config.env | cut -d'=' -f2)
DB_PORT=$(grep DB_PORT config.env | cut -d'=' -f2)
DB_NAME=$(grep DB_NAME config.env | cut -d'=' -f2)

if command -v psql &> /dev/null; then
    if PGPASSWORD=$(grep DB_PASSWORD config.env | cut -d'=' -f2) psql -h $DB_HOST -p $DB_PORT -U $(grep DB_USER config.env | cut -d'=' -f2) -d $DB_NAME -c "SELECT 1;" &> /dev/null; then
        echo "✅ Database connection test successful!"
    else
        echo "❌ Database connection failed. Please check your config.env"
        exit 1
    fi
else
    echo "⚠️  psql not found. Skipping database connection test."
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod tidy

# Build the service
echo "🔨 Building service..."
go build -o user-service cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed!"
    exit 1
fi

# Start the service
echo "🚀 Starting User Service..."
./user-service
