#!/bin/bash

# NurseShift Microservices Stop Script
# This script stops all microservices

PROJECT_ROOT="/Volumes/NO NAME/project/nurseshift_final"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Services configuration (service_name:port)
SERVICES=(
    "auth-service:8081"
    "user-service:8082"
    "department-service:8083"
    "schedule-service:8084"
    "setting-service:8085"
    "priority-service:8086"
    "notification-service:8087"
    "package-service:8088"
    "payment-service:8089"
    "employee-leave-service:8090"
)

echo -e "${YELLOW}🛑 Stopping all NurseShift Microservices...${NC}"
echo "=================================="

# Stop each service by port
for service_entry in "${SERVICES[@]}"; do
    IFS=':' read -r service_name port <<< "$service_entry"
    
    echo -e "${BLUE}🛑 Stopping $service_name on port $port...${NC}"
    
    # Kill process using the port
    if lsof -ti:$port >/dev/null 2>&1; then
        pids=$(lsof -ti:$port)
        for pid in $pids; do
            echo -e "${YELLOW}  Killing process PID: $pid${NC}"
            kill -9 $pid 2>/dev/null
        done
        echo -e "${GREEN}✅ $service_name stopped${NC}"
    else
        echo -e "${YELLOW}⚠️  No process found on port $port${NC}"
    fi
done

# Clean up PID files
echo -e "${BLUE}🧹 Cleaning up PID files...${NC}"
for service_entry in "${SERVICES[@]}"; do
    IFS=':' read -r service_name port <<< "$service_entry"
    pid_file="/tmp/$service_name.pid"
    if [ -f "$pid_file" ]; then
        rm -f "$pid_file"
        echo -e "${GREEN}✅ Cleaned up $service_name PID file${NC}"
    fi
done

echo "=================================="
echo -e "${GREEN}✅ All services stopped successfully!${NC}"
