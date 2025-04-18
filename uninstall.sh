#!/bin/bash

# Fort-Go Uninstallation Script
# This script removes Fort-Go from your system

set -e  # Exit on error

# Text formatting
BOLD="\033[1m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
RED="\033[0;31m"
YELLOW="\033[0;33m"
NC="\033[0m" # No Color

echo -e "${BOLD}${BLUE}"
echo "  ______          _   _    _____      "
echo " |  ____|        | | (_)  / ____|     "
echo " | |__ ___  _ __ | |_ _  | |  __  ___ "
echo " |  __/ _ \| '_ \| __| | | | |_ |/ _ \\"
echo " | | | (_) | | | | |_| | | |__| | (_) |"
echo " |_|  \___/|_| |_|\__|_|  \_____|\___/"
echo -e "${NC}"
echo -e "${BOLD}Fort-Go Uninstallation Script${NC}"
echo

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
  echo -e "${YELLOW}Notice: Not running as root. Uninstallation may require sudo privileges.${NC}"
  USE_SUDO="sudo"
else
  USE_SUDO=""
fi

echo -e "${YELLOW}Warning: This will remove Fort-Go from your system.${NC}"
read -p "Are you sure you want to continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    echo "Uninstallation cancelled."
    exit 0
fi

# Remove binary
if [ -f "/usr/local/bin/fort" ]; then
    echo "Removing Fort-Go binary..."
    $USE_SUDO rm -f /usr/local/bin/fort
    echo -e "${GREEN}✓ Binary removed${NC}"
else
    echo -e "${YELLOW}! Binary not found in /usr/local/bin${NC}"
fi

# Remove repository
INSTALL_DIR="$HOME/.fort-go"
if [ -d "$INSTALL_DIR" ]; then
    echo "Removing Fort-Go repository..."
    rm -rf "$INSTALL_DIR"
    echo -e "${GREEN}✓ Repository removed${NC}"
else
    echo -e "${YELLOW}! Repository not found at $INSTALL_DIR${NC}"
fi

# Remove shell completion
SHELL_TYPE=$(basename "$SHELL")
case "$SHELL_TYPE" in
    bash)
        COMPLETION_FILE="$HOME/.bash_completion"
        COMPLETION_SCRIPT="$HOME/.bash_completion.d/fort.bash"
        
        if [ -f "$COMPLETION_SCRIPT" ]; then
            echo "Removing Bash completion..."
            rm -f "$COMPLETION_SCRIPT"
            
            # Remove from .bash_completion if it exists
            if [ -f "$COMPLETION_FILE" ] && grep -q "fort.bash" "$COMPLETION_FILE"; then
                sed -i '/# Fort-Go completion/d' "$COMPLETION_FILE"
                sed -i '/fort.bash/d' "$COMPLETION_FILE"
            fi
            
            echo -e "${GREEN}✓ Bash completion removed${NC}"
        else
            echo -e "${YELLOW}! Bash completion not found${NC}"
        fi
        ;;
    zsh)
        COMPLETION_SCRIPT="$HOME/.zsh/completion/_fort"
        
        if [ -f "$COMPLETION_SCRIPT" ]; then
            echo "Removing Zsh completion..."
            rm -f "$COMPLETION_SCRIPT"
            
            # Remove from .zshrc if it exists
            if grep -q "completion/_fort" "$HOME/.zshrc"; then
                sed -i '/# Fort-Go completion/d' "$HOME/.zshrc"
                sed -i '/completion\/_fort/d' "$HOME/.zshrc"
            fi
            
            echo -e "${GREEN}✓ Zsh completion removed${NC}"
        else
            echo -e "${YELLOW}! Zsh completion not found${NC}"
        fi
        ;;
    *)
        echo -e "${YELLOW}! Shell completion for $SHELL_TYPE not supported. Skipping.${NC}"
        ;;
esac

echo
echo -e "${GREEN}${BOLD}Fort-Go has been successfully uninstalled!${NC}" 