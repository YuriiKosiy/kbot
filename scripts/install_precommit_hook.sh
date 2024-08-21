#!/bin/bash

# Directory of the script
SCRIPT_DIR=$(dirname "$0")

# Destination for the pre-commit hook
HOOK_DESTINATION=".git/hooks/pre-commit"

# Check if the hook already exists
if [ -f "$HOOK_DESTINATION" ]; then
    echo "Pre-commit hook already exists. Skipping installation."
else
    echo "Installing pre-commit hook..."

    # Copy the pre-commit hook from the scripts directory
    cp "$SCRIPT_DIR/pre-commit" "$HOOK_DESTINATION"
    chmod +x "$HOOK_DESTINATION"

    echo "Pre-commit hook installed successfully."
fi
