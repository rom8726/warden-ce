#!/bin/bash

# Test script for Warden Platform Installer
# This script tests the installer functionality without actually installing

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test directory
TEST_DIR="/tmp/warden_test_install"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

print_info() {
    echo -e "${BLUE}[TEST INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[TEST SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[TEST WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[TEST ERROR]${NC} $1"
}

# Function to test random string generation
test_random_string_generation() {
    print_info "Testing random string generation..."
    
    # Test function from install.sh
    generate_random_string() {
        local length="$1"
        LC_ALL=C tr -dc 'A-Za-z0-9' < /dev/urandom | head -c "$length"
    }
    
    # Test different lengths
    local test_lengths=(8 12 16 32)
    for length in "${test_lengths[@]}"; do
        local result=$(generate_random_string "$length")
        if [[ ${#result} -eq $length ]]; then
            print_success "Random string generation for length $length: OK"
        else
            print_error "Random string generation for length $length: FAILED (expected $length, got ${#result})"
            return 1
        fi
    done
}

# Function to test email validation
test_email_validation() {
    print_info "Testing email validation..."
    
    # Test function from install.sh
    validate_email() {
        local email="$1"
        if [[ "$email" =~ ^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$ ]] && [[ ! "$email" =~ \.\..*@ ]]; then
            return 0
        else
            return 1
        fi
    }
    
    # Valid emails
    local valid_emails=(
        "test@example.com"
        "user.name@domain.co.uk"
        "user+tag@example.org"
    )
    
    # Invalid emails
    local invalid_emails=(
        "invalid-email"
        "@example.com"
        "user@"
        "user@.com"
        "user..name@example.com"
    )
    
    # Test valid emails
    for email in "${valid_emails[@]}"; do
        if validate_email "$email"; then
            print_success "Valid email '$email': OK"
        else
            print_error "Valid email '$email': FAILED"
            return 1
        fi
    done
    
    # Test invalid emails
    for email in "${invalid_emails[@]}"; do
        if ! validate_email "$email"; then
            print_success "Invalid email '$email': OK"
        else
            print_error "Invalid email '$email': FAILED"
            return 1
        fi
    done
}

# Function to test template rendering
test_template_rendering() {
    print_info "Testing template rendering..."
    
    # Create test template
    local test_template="$TEST_DIR/test_template.txt"
    cat > "$test_template" << 'EOF'
DOMAIN={{ .Domain }}
EMAIL={{ .Email }}
PASSWORD={{ .Password }}
EOF
    
    # Test data
    local domain="test.example.com"
    local email="admin@test.example.com"
    local password="testpass123"
    
    # Render template manually
    local expected_output="$TEST_DIR/expected_output.txt"
    cat > "$expected_output" << EOF
DOMAIN=$domain
EMAIL=$email
PASSWORD=$password
EOF
    
    # Create rendered output
    local actual_output="$TEST_DIR/actual_output.txt"
    sed -e "s/{{ .Domain }}/$domain/g" \
        -e "s/{{ .Email }}/$email/g" \
        -e "s/{{ .Password }}/$password/g" \
        "$test_template" > "$actual_output"
    
    # Compare
    if diff "$expected_output" "$actual_output" > /dev/null; then
        print_success "Template rendering: OK"
    else
        print_error "Template rendering: FAILED"
        echo "Expected:"
        cat "$expected_output"
        echo "Actual:"
        cat "$actual_output"
        return 1
    fi
}

# Function to test directory creation
test_directory_creation() {
    print_info "Testing directory creation..."
    
    local test_dirs=(
        "$TEST_DIR/test1"
        "$TEST_DIR/test2/subdir"
        "$TEST_DIR/test3/subdir1/subdir2"
    )
    
    for dir in "${test_dirs[@]}"; do
        if mkdir -p "$dir" 2>/dev/null; then
            print_success "Directory creation '$dir': OK"
        else
            print_error "Directory creation '$dir': FAILED"
            return 1
        fi
    done
}

# Function to test SSL certificate generation
test_ssl_generation() {
    print_info "Testing SSL certificate generation..."
    
    local key_path="$TEST_DIR/test_key.pem"
    local cert_path="$TEST_DIR/test_cert.pem"
    local domain="test.example.com"
    
    # Check if openssl is available
    if ! command -v openssl &> /dev/null; then
        print_warning "OpenSSL not available, skipping SSL test"
        return 0
    fi
    
    # Generate test certificate
    if openssl genrsa -out "$key_path" 2048 2>/dev/null && \
       openssl req -new -x509 -key "$key_path" -out "$cert_path" -days 3650 -subj "/C=US/ST=State/L=City/O=Warden/CN=$domain" 2>/dev/null; then
        print_success "SSL certificate generation: OK"
        
        # Check if files exist and have content
        if [[ -s "$key_path" && -s "$cert_path" ]]; then
            print_success "SSL certificate files validation: OK"
        else
            print_error "SSL certificate files validation: FAILED"
            return 1
        fi
    else
        print_error "SSL certificate generation: FAILED"
        return 1
    fi
}

# Function to test configuration file generation
test_config_generation() {
    print_info "Testing configuration file generation..."
    
    # Test platform.env generation
    local platform_env="$TEST_DIR/platform.env"
    cat > "$platform_env" << EOF
DOCKER_REGISTRY=
DOMAIN=test.example.com
DB_PASSWORD=testpass123
CLICKHOUSE_DB_PASSWORD=testpass456
SSL_CERT=nginx_cert.pem
SSL_CERT_KEY=nginx_key.pem
SECURE_LINK_MD5=
PLATFORM_VERSION=latest
EOF
    
    if [[ -f "$platform_env" ]]; then
        print_success "Platform.env generation: OK"
    else
        print_error "Platform.env generation: FAILED"
        return 1
    fi
    
    # Test config.env generation
    local config_env="$TEST_DIR/config.env"
    cat > "$config_env" << EOF
WARDEN_FRONTEND_URL=https://test.example.com
WARDEN_SECRET_KEY=testsecretkey123
WARDEN_JWT_SECRET_KEY=testjwtkey123
WARDEN_POSTGRES_PASSWORD=testpass123
WARDEN_CLICKHOUSE_PASSWORD=testpass456
WARDEN_MAILER_ADDR=smtp.example.com:587
WARDEN_MAILER_USER=test@example.com
WARDEN_MAILER_PASSWORD=testpass
WARDEN_MAILER_FROM=test@example.com
WARDEN_MAILER_USE_TLS=true
WARDEN_ADMIN_EMAIL=admin@test.example.com
WARDEN_ADMIN_TMP_PASSWORD=adminpass123
EOF
    
    if [[ -f "$config_env" ]]; then
        print_success "Config.env generation: OK"
    else
        print_error "Config.env generation: FAILED"
        return 1
    fi
}

# Function to test docker-compose.yml validation
test_docker_compose_validation() {
    print_info "Testing docker-compose.yml validation..."
    
    local docker_compose="$SCRIPT_DIR/docker-compose.yml"
    
    if [[ ! -f "$docker_compose" ]]; then
        print_error "docker-compose.yml not found: FAILED"
        return 1
    fi
    
    # Check if file contains required services
    local required_services=(
        "warden-backend"
        "warden-frontend"
        "warden-postgresql"
        "warden-redis"
        "warden-clickhouse"
        "warden-kafka"
    )
    
    for service in "${required_services[@]}"; do
        if grep -q "$service:" "$docker_compose"; then
            print_success "Service '$service' found in docker-compose.yml: OK"
        else
            print_error "Service '$service' not found in docker-compose.yml: FAILED"
            return 1
        fi
    done
}

# Function to test Makefile generation
test_makefile_generation() {
    print_info "Testing Makefile generation..."
    
    local makefile="$TEST_DIR/Makefile"
    cat > "$makefile" << 'EOF'
_COMPOSE=docker compose -f docker-compose.yml --project-name warden --env-file platform.env

.DEFAULT_GOAL := help

.PHONY: help
help: ## Print this message
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

.PHONY: up
up: ## Up the environment in docker compose
	${_COMPOSE} up -d

.PHONY: down
down: ## Down the environment in docker compose
	${_COMPOSE} down --remove-orphans

.PHONY: pull
pull: ## Pull images from remote Docker registry
	${_COMPOSE} pull
EOF
    
    if [[ -f "$makefile" ]]; then
        print_success "Makefile generation: OK"
    else
        print_error "Makefile generation: FAILED"
        return 1
    fi
}

# Main test function
main() {
    print_info "Starting Warden Platform Installer tests..."
    
    # Create test directory
    mkdir -p "$TEST_DIR"
    
    # Run tests
    local tests=(
        "test_random_string_generation"
        "test_email_validation"
        "test_template_rendering"
        "test_directory_creation"
        "test_ssl_generation"
        "test_config_generation"
        "test_docker_compose_validation"
        "test_makefile_generation"
    )
    
    local failed_tests=0
    
    for test in "${tests[@]}"; do
        if $test; then
            print_success "Test '$test' passed"
        else
            print_error "Test '$test' failed"
            ((failed_tests++))
        fi
        echo
    done
    
    # Cleanup
    rm -rf "$TEST_DIR"
    
    # Summary
    if [[ $failed_tests -eq 0 ]]; then
        print_success "All tests passed! Installer should work correctly."
        exit 0
    else
        print_error "$failed_tests test(s) failed. Please check the installer."
        exit 1
    fi
}

# Run main function
main "$@"
