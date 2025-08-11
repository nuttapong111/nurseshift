#!/bin/bash

# Swagger Documentation Generator Script
PROJECT_ROOT="/Volumes/NO NAME/project/nurseshift"
BACKEND_DIR="$PROJECT_ROOT/backend"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

# Services with their details
declare -a SERVICES=(
    "auth-service:8081:Authentication microservice for NurseShift application"
    "user-service:8082:User management and profile microservice"
    "department-service:8083:Department and employee management microservice"
    "schedule-service:8084:Shift scheduling and management microservice"
    "setting-service:8085:System settings and configuration microservice"
    "priority-service:8086:Scheduling priority management microservice"
    "notification-service:8087:Notification and alert microservice"
    "package-service:8088:Membership package management microservice"
    "payment-service:8089:Payment processing and history microservice"
)

# Function to create basic swagger docs for a service
create_swagger_docs() {
    local service_name=$1
    local port=$2
    local description=$3
    local service_dir="$BACKEND_DIR/$service_name"
    
    echo -e "${BLUE}📚 Creating Swagger docs for $service_name...${NC}"
    
    # Create docs directory
    mkdir -p "$service_dir/docs"
    
    # Create basic swagger.yaml
    cat > "$service_dir/docs/swagger.yaml" << EOF
swagger: "2.0"
info:
  title: "NurseShift ${service_name^} API"
  description: "$description"
  version: "1.0.0"
  contact:
    name: "NurseShift Team"
    email: "support@nurseshift.com"
host: "localhost:$port"
basePath: "/api/v1"
schemes:
  - "http"
  - "https"

securityDefinitions:
  BearerAuth:
    type: apiKey
    in: header
    name: Authorization
    description: "JWT token in format: Bearer {token}"

consumes:
  - "application/json"
produces:
  - "application/json"

tags:
  - name: "${service_name^}"
    description: "${service_name^} operations"
  - name: "Health"
    description: "Service health check"

paths:
  /health:
    get:
      tags:
        - "Health"
      summary: "Service health check"
      description: "Check if the ${service_name} is running properly"
      produces:
        - "application/json"
      responses:
        200:
          description: "Service is healthy"
          schema:
            type: object
            properties:
              status:
                type: string
                example: "ok"
              service:
                type: string
                example: "$service_name"
              timestamp:
                type: string
                format: date-time

definitions:
  ErrorResponse:
    type: object
    properties:
      status:
        type: string
        example: "error"
      message:
        type: string
        description: "Error message in Thai"
      error:
        type: string
        description: "Technical error details (optional)"
    required:
      - status
      - message

  SuccessResponse:
    type: object
    properties:
      status:
        type: string
        example: "success"
      message:
        type: string
        description: "Success message in Thai"
      data:
        type: object
        description: "Response data"
    required:
      - status
      - message
EOF

    echo -e "${GREEN}✅ Created swagger.yaml for $service_name${NC}"
}

# Main function
echo -e "${BLUE}🚀 Generating Swagger Documentation for All Services...${NC}"
echo "================================================================"

for service_entry in "${SERVICES[@]}"; do
    IFS=':' read -r service_name port description <<< "$service_entry"
    create_swagger_docs "$service_name" "$port" "$description"
done

echo "================================================================"
echo -e "${GREEN}✅ All Swagger documentation generated successfully!${NC}"
echo ""
echo -e "${BLUE}📚 Access Swagger UI:${NC}"

for service_entry in "${SERVICES[@]}"; do
    IFS=':' read -r service_name port description <<< "$service_entry"
    echo "  🔗 $service_name: http://localhost:$port/swagger/"
done

echo ""
echo -e "${BLUE}💡 Note: Auth Service has the most complete documentation${NC}"
echo -e "${BLUE}💡 Other services have basic documentation templates${NC}"


