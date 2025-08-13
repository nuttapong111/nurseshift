#!/bin/bash

# NurseShift Microservices Startup Script
# This script starts all microservices for local development

PROJECT_ROOT="/Volumes/NO NAME/project/nurseshift_final"
BACKEND_DIR="$PROJECT_ROOT/backend"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
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

# Function to check if port is available
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
        return 1  # Port is in use
    else
        return 0  # Port is available
    fi
}

# Function to start a service
start_service() {
    local service_name=$1
    local port=$2
    local service_dir="$BACKEND_DIR/$service_name"
    
    echo -e "${BLUE}🚀 Starting $service_name on port $port...${NC}"
    
    # Check if service directory exists
    if [ ! -d "$service_dir" ]; then
        echo -e "${RED}❌ Service directory not found: $service_dir${NC}"
        return 1
    fi
    
    # Check if port is available
    if ! check_port $port; then
        echo -e "${YELLOW}⚠️  Port $port is already in use, skipping $service_name${NC}"
        return 1
    fi
    
    # Start the service in background
    cd "$service_dir"
    
    # Export environment variables from config.env if it exists
    if [ -f "config.env" ]; then
        echo -e "${CYAN}📋 Loading environment variables from config.env...${NC}"
        export $(grep -v "^#" config.env | grep -v "^$" | xargs)
    fi
    
    # Build and run the service
    if go mod tidy && go build -o "$service_name" cmd/server/main.go; then
        ./"$service_name" > "/tmp/$service_name.log" 2>&1 &
        local pid=$!
        echo $pid > "/tmp/$service_name.pid"
        
        # Wait a moment and check if service started successfully
        sleep 2
        if kill -0 $pid 2>/dev/null; then
            echo -e "${GREEN}✅ $service_name started successfully (PID: $pid)${NC}"
            return 0
        else
            echo -e "${RED}❌ Failed to start $service_name${NC}"
            return 1
        fi
    else
        echo -e "${RED}❌ Failed to build $service_name${NC}"
        return 1
    fi
}

# Function to check service health
check_health() {
    local service_name=$1
    local port=$2
    
    echo -e "${CYAN}🔍 Checking health of $service_name...${NC}"
    
    local response=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:$port/health" 2>/dev/null)
    
    if [ "$response" = "200" ]; then
        echo -e "${GREEN}✅ $service_name is healthy${NC}"
        return 0
    else
        echo -e "${RED}❌ $service_name health check failed (HTTP: $response)${NC}"
        return 1
    fi
}

# Function to stop all services
stop_all_services() {
    echo -e "${YELLOW}🛑 Stopping all services...${NC}"
    
    for service_entry in "${SERVICES[@]}"; do
        IFS=':' read -r service_name port <<< "$service_entry"
        local pid_file="/tmp/$service_name.pid"
        if [ -f "$pid_file" ]; then
            local pid=$(cat "$pid_file")
            if kill -0 $pid 2>/dev/null; then
                echo -e "${YELLOW}🛑 Stopping $service_name (PID: $pid)...${NC}"
                kill $pid
                rm -f "$pid_file"
            fi
        fi
    done
    
    echo -e "${GREEN}✅ All services stopped${NC}"
}

# Function to show service status
show_status() {
    echo -e "${PURPLE}📊 Service Status:${NC}"
    echo "=================================="
    
    for service_entry in "${SERVICES[@]}"; do
        IFS=':' read -r service_name port <<< "$service_entry"
        local pid_file="/tmp/$service_name.pid"
        
        printf "%-20s Port %-4s " "$service_name" "$port"
        
        if [ -f "$pid_file" ]; then
            local pid=$(cat "$pid_file")
            if kill -0 $pid 2>/dev/null; then
                if check_port $port; then
                    echo -e "${RED}❌ Not Running${NC}"
                else
                    echo -e "${GREEN}✅ Running (PID: $pid)${NC}"
                fi
            else
                echo -e "${RED}❌ Dead Process${NC}"
                rm -f "$pid_file"
            fi
        else
            echo -e "${RED}❌ Not Started${NC}"
        fi
    done
    
    echo "=================================="
}

# Function to show logs
show_logs() {
    local service_name=$1
    local log_file="/tmp/$service_name.log"
    
    if [ -f "$log_file" ]; then
        echo -e "${CYAN}📋 Logs for $service_name:${NC}"
        tail -20 "$log_file"
    else
        echo -e "${RED}❌ No logs found for $service_name${NC}"
    fi
}

# Main script logic
case "${1:-start}" in
    "start")
        echo -e "${BLUE}🚀 Starting NurseShift Microservices...${NC}"
        echo "=================================="
        
        # Change to project root
        cd "$PROJECT_ROOT"
        
        # Start each service
        for service_entry in "${SERVICES[@]}"; do
            IFS=':' read -r service_name port <<< "$service_entry"
            start_service "$service_name" "$port"
        done
        
        echo "=================================="
        echo -e "${BLUE}⏳ Waiting for services to initialize...${NC}"
        sleep 5
        
        # Check health of all services
        echo "=================================="
        healthy_count=0
        for service_entry in "${SERVICES[@]}"; do
            IFS=':' read -r service_name port <<< "$service_entry"
            if check_health "$service_name" "$port"; then
                ((healthy_count++))
            fi
        done
        
        echo "=================================="
        echo -e "${GREEN}🎉 $healthy_count/${#SERVICES[@]} services are running${NC}"
        
        if [ $healthy_count -eq ${#SERVICES[@]} ]; then
            echo -e "${GREEN}✅ All services started successfully!${NC}"
            echo -e "${CYAN}📚 API Endpoints:${NC}"
            for service_entry in "${SERVICES[@]}"; do
                IFS=':' read -r service_name port <<< "$service_entry"
                echo "  🔗 $service_name: http://localhost:$port/health"
            done
        else
            echo -e "${YELLOW}⚠️  Some services failed to start. Check logs with: $0 logs <service-name>${NC}"
        fi
        ;;
        
    "stop")
        stop_all_services
        ;;
        
    "status")
        show_status
        ;;
        
    "logs")
        if [ -n "$2" ]; then
            show_logs "$2"
        else
            echo -e "${RED}❌ Please specify service name: $0 logs <service-name>${NC}"
            echo "Available services:"
            for service_entry in "${SERVICES[@]}"; do
                IFS=':' read -r service_name port <<< "$service_entry"
                echo "  - $service_name"
            done
        fi
        ;;
        
    "restart")
        stop_all_services
        sleep 3
        $0 start
        ;;
        
    "health")
        echo -e "${CYAN}🔍 Health Check for All Services:${NC}"
        echo "=================================="
        for service_entry in "${SERVICES[@]}"; do
            IFS=':' read -r service_name port <<< "$service_entry"
            check_health "$service_name" "$port"
        done
        ;;
        
    *)
        echo "Usage: $0 {start|stop|restart|status|health|logs <service-name>}"
        echo ""
        echo "Commands:"
        echo "  start    - Start all microservices"
        echo "  stop     - Stop all microservices"
        echo "  restart  - Restart all microservices"
        echo "  status   - Show status of all services"
        echo "  health   - Check health of all services"
        echo "  logs     - Show logs for specific service"
        echo ""
        echo "Available services:"
        for service_entry in "${SERVICES[@]}"; do
            IFS=':' read -r service_name port <<< "$service_entry"
            echo "  - $service_name (port $port)"
        done
        exit 1
        ;;
esac
