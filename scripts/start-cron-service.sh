#!/bin/bash

# NurseShift Cron Service Startup Script
# This script starts the cron service for updating user days remaining

PROJECT_ROOT="/Volumes/NO NAME/project/nurseshift_final"
CRON_SERVICE_DIR="$PROJECT_ROOT/backend/cron-service"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting NurseShift Cron Service...${NC}"

# Check if cron service directory exists
if [ ! -d "$CRON_SERVICE_DIR" ]; then
    echo -e "${RED}❌ Cron service directory not found: $CRON_SERVICE_DIR${NC}"
    exit 1
fi

# Change to cron service directory
cd "$CRON_SERVICE_DIR"

# Check if binary exists
if [ ! -f "cron-service" ]; then
    echo -e "${YELLOW}⚠️  Cron service binary not found, building...${NC}"
    if go build -o cron-service cmd/server/main.go; then
        echo -e "${GREEN}✅ Cron service built successfully${NC}"
    else
        echo -e "${RED}❌ Failed to build cron service${NC}"
        exit 1
    fi
fi

# Check if cron service is already running
if pgrep -f "cron-service" > /dev/null; then
    echo -e "${YELLOW}⚠️  Cron service is already running${NC}"
    echo -e "${BLUE}📊 Current cron service processes:${NC}"
    ps aux | grep "cron-service" | grep -v grep
    exit 0
fi

# Start cron service in background
echo -e "${BLUE}🚀 Starting cron service...${NC}"
./cron-service > "/tmp/cron-service.log" 2>&1 &
CRON_PID=$!

# Wait a moment and check if service started successfully
sleep 3
if kill -0 $CRON_PID 2>/dev/null; then
    echo -e "${GREEN}✅ Cron service started successfully (PID: $CRON_PID)${NC}"
    echo $CRON_PID > "/tmp/cron-service.pid"
    echo -e "${BLUE}📝 Logs are being written to /tmp/cron-service.log${NC}"
    echo -e "${BLUE}📅 Service will update user days remaining daily at midnight UTC (7 AM Thailand time)${NC}"
    echo -e "${BLUE}📊 Service will log status every 6 hours${NC}"
else
    echo -e "${RED}❌ Failed to start cron service${NC}"
    echo -e "${BLUE}📝 Check logs at /tmp/cron-service.log${NC}"
    exit 1
fi

echo -e "${GREEN}🎉 Cron service is now running!${NC}"
echo -e "${BLUE}💡 To stop the service, run: kill $CRON_PID${NC}"
echo -e "${BLUE}💡 To view logs, run: tail -f /tmp/cron-service.log${NC}"
