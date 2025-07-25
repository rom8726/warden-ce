#!/bin/bash

# Load Test Runner for Envelope Endpoint
# Usage: ./run_load_test.sh [profile] [duration]

set -e

# Default values
DEFAULT_PROFILE="quick"
DEFAULT_DURATION="30s"
DEFAULT_PROJECT_ID="1"
DEFAULT_SENTRY_KEY="4678fd4d2cb5500cff8f33d02a041f96d2e05bf5eb467e7d10a9bddbfc298bab"
DEFAULT_ADDR="http://127.0.0.1:8098"

# Environment variables (can be overridden)
PROJECT_ID=${PROJECT_ID:-$DEFAULT_PROJECT_ID}
SENTRY_KEY=${SENTRY_KEY:-$DEFAULT_SENTRY_KEY}
APP_ADDR=${APP_ADDR:-$DEFAULT_ADDR}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Load Test Runner for Envelope Endpoint"
    echo ""
    echo "Usage: $0 [profile] [duration]"
    echo ""
    echo "Environment Variables:"
    echo "  PROJECT_ID    - Project ID (default: $DEFAULT_PROJECT_ID)"
    echo "  SENTRY_KEY    - Sentry public key (default: $DEFAULT_SENTRY_KEY)"
    echo ""
    echo "Profiles:"
    echo "  quick     - Quick test (5 threads, 30s)"
    echo "  light     - Light load (10 threads, 1m)"
    echo "  medium    - Medium load (25 threads, 2m)"
    echo "  heavy     - Heavy load (50 threads, 5m)"
    echo "  stress    - Stress test (100 threads, 10m)"
    echo "  custom    - Custom parameters (use -h for options)"
    echo ""
    echo "Examples:"
    echo "  $0 quick"
    echo "  $0 medium 3m"
    echo "  PROJECT_ID=123 SENTRY_KEY=mykey $0 heavy"
    echo "  $0 custom -threads 20 -duration 5m -rps 10"
    echo ""
}

# Function to check if load_test_envelope exists
check_binary() {
    if [ ! -f "./load_test_envelope" ]; then
        print_error "load_test_envelope binary not found!"
        print_info "Building binary..."
        go build -o load_test_envelope load_test_envelope.go
        if [ $? -ne 0 ]; then
            print_error "Failed to build load_test_envelope"
            exit 1
        fi
        print_success "Binary built successfully"
    fi
}

# Function to run quick test
run_quick() {
    print_info "Running quick test (5 threads, $1)"
    ./load_test_envelope -threads 5 -duration "$1" -rps 5 -project "$PROJECT_ID" -key "$SENTRY_KEY" -url "$APP_ADDR"
}

# Function to run light load test
run_light() {
    print_info "Running light load test (10 threads, $1)"
    ./load_test_envelope -threads 10 -duration "$1" -rps 10 -project "$PROJECT_ID" -key "$SENTRY_KEY" -url "$APP_ADDR"
}

# Function to run medium load test
run_medium() {
    print_info "Running medium load test (25 threads, $1)"
    ./load_test_envelope -threads 25 -duration "$1" -rps 15 -project "$PROJECT_ID" -key "$SENTRY_KEY" -url "$APP_ADDR"
}

# Function to run heavy load test
run_heavy() {
    print_info "Running heavy load test (50 threads, $1)"
    ./load_test_envelope -threads 50 -duration "$1" -rps 20 -project "$PROJECT_ID" -key "$SENTRY_KEY" -url "$APP_ADDR"
}

# Function to run stress test
run_stress() {
    print_info "Running stress test (100 threads, $1)"
    ./load_test_envelope -threads 100 -duration "$1" -rps 25 -project "$PROJECT_ID" -key "$SENTRY_KEY" -url "$APP_ADDR"
}

# Function to run custom test
run_custom() {
    print_info "Running custom test with parameters: $*"
    ./load_test_envelope "$@"
}

# Main script logic
main() {
    # Check if help is requested
    if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
        show_usage
        exit 0
    fi
    
    # Check binary
    check_binary
    
    # Get profile and duration
    PROFILE=${1:-$DEFAULT_PROFILE}
    DURATION=${2:-$DEFAULT_DURATION}
    
    print_info "Starting load test with profile: $PROFILE"
    print_info "Target: http://localhost:8080/api/$PROJECT_ID/envelope/"
    print_info "Project ID: $PROJECT_ID"
    print_info "Sentry Key: ${SENTRY_KEY:0:16}..."
    echo ""
    
    # Run test based on profile
    case $PROFILE in
        "quick")
            run_quick "$DURATION"
            ;;
        "light")
            run_light "$DURATION"
            ;;
        "medium")
            run_medium "$DURATION"
            ;;
        "heavy")
            run_heavy "$DURATION"
            ;;
        "stress")
            run_stress "$DURATION"
            ;;
        "custom")
            shift
            run_custom "$@"
            ;;
        *)
            print_error "Unknown profile: $PROFILE"
            show_usage
            exit 1
            ;;
    esac
    
    print_success "Load test completed!"
}

# Run main function with all arguments
main "$@" 
