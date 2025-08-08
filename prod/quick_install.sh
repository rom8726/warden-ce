#!/bin/bash

# Quick installer launcher for Warden Platform
# This script checks dependencies and launches the main installer

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Function to check if command exists
check_command() {
    local cmd="$1"
    local name="$2"
    
    if command -v "$cmd" &> /dev/null; then
        print_success "$name is available"
        return 0
    else
        print_error "$name is not available"
        return 1
    fi
}

# Function to check system requirements
check_requirements() {
    print_info "Checking system requirements..."
    
    local missing_deps=0
    
    # Check for bash
    if ! check_command "bash" "Bash"; then
        ((missing_deps++))
    fi
    
    # Check for docker
    if ! check_command "docker" "Docker"; then
        ((missing_deps++))
    fi
    
    # Check for docker compose
    if ! check_command "docker-compose" "Docker Compose" && ! docker compose version &> /dev/null; then
        print_error "Docker Compose is not available"
        ((missing_deps++))
    else
        print_success "Docker Compose is available"
    fi
    
    # Check for openssl
    if ! check_command "openssl" "OpenSSL"; then
        print_warning "OpenSSL is not available - SSL certificate generation will be skipped"
    fi
    
    # Check for make
    if ! check_command "make" "Make"; then
        print_warning "Make is not available - you'll need to use docker compose commands directly"
    fi
    
    # Check if running as root
    if [[ $EUID -ne 0 ]]; then
        print_error "This installer must be run as root (sudo)"
        print_info "Please run: sudo $0"
        exit 1
    fi
    
    if [[ $missing_deps -gt 0 ]]; then
        print_error "Missing required dependencies. Please install them before running the installer."
        exit 1
    fi
    
    print_success "All requirements met!"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -t, --test     Run tests only"
    echo "  -v, --version  Show version information"
    echo
    echo "Examples:"
    echo "  sudo $0                    # Run installer"
    echo "  sudo $0 --test            # Run tests only"
    echo "  sudo $0 --help            # Show this help"
}

# Function to show version
show_version() {
    echo "Warden Platform Installer (Bash Version)"
    echo "Version: 1.0.0"
    echo "Based on Go installer from internal/installer"
}

# Main function
main() {
    # Parse command line arguments
    case "${1:-}" in
        -h|--help)
            show_usage
            exit 0
            ;;
        -v|--version)
            show_version
            exit 0
            ;;
        -t|--test)
            print_info "Running tests only..."
            ./test_installer.sh
            exit $?
            ;;
        "")
            # No arguments, run installer
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
    
    # Show version
    show_version
    echo
    
    # Check requirements
    check_requirements
    echo
    
    # Confirm installation
    print_info "Ready to start Warden Platform installation."
    echo -n "Press Enter to continue or Ctrl+C to cancel: "
    read -r
    
    # Run installer
    print_info "Starting installer..."
    ./install.sh
}

# Run main function
main "$@"
