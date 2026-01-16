#!/usr/bin/env bash
set -eu

# Function to check if command exists in PATH
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to install via brew if available, otherwise pip
install_dependency() {
    local dep="$1"
    local pip_package="${2:-$1}"

    if command_exists "$dep"; then
        echo "$dep is already installed"
        return 0
    fi

    echo "Installing $dep..."

    if command_exists brew; then
        brew install "$dep"
    elif command_exists pip || command_exists pip3; then
        local pip_cmd=$(command_exists pip3 && echo "pip3" || echo "pip")
        $pip_cmd install --user "$pip_package"
    else
        echo "Error: Neither brew nor pip found. Please install one of them first."
        exit 1
    fi
}

# Install dev dependencies
echo "Installing development dependencies..."

install_dependency "pre-commit"
install_dependency "detect-secrets"

echo "All development dependencies installed successfully!"
