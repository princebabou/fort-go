#!/bin/bash

# Fort-Go Installation Script
# This script installs Fort-Go and its dependencies

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
echo -e "${BOLD}Fort-Go Installation Script${NC}"
echo

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
  echo -e "${YELLOW}Notice: Not running as root. Installation may require sudo privileges.${NC}"
  USE_SUDO="sudo"
else
  USE_SUDO=""
fi

# Check for required tools
check_required() {
  echo -e "${BOLD}Checking required tools...${NC}"
  
  # Check for Go
  if command -v go >/dev/null 2>&1; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "✅ Go is installed: ${GO_VERSION}"
  else
    echo -e "${RED}❌ Go is not installed${NC}"
    echo -e "Installing Go..."
    
    # Determine OS
    OS="$(uname -s)"
    case "${OS}" in
      Linux*)
        if command -v apt-get >/dev/null 2>&1; then
          $USE_SUDO apt-get update
          $USE_SUDO apt-get install -y golang-go
        elif command -v yum >/dev/null 2>&1; then
          $USE_SUDO yum install -y golang
        elif command -v pacman >/dev/null 2>&1; then
          $USE_SUDO pacman -S go --noconfirm
        else
          echo -e "${RED}Unsupported Linux distribution. Please install Go manually.${NC}"
          exit 1
        fi
        ;;
      Darwin*)
        if command -v brew >/dev/null 2>&1; then
          brew install go
        else
          echo -e "${RED}Homebrew not found. Please install Go manually.${NC}"
          exit 1
        fi
        ;;
      MINGW*|MSYS*|CYGWIN*)
        echo -e "${RED}Windows detected. Please install Go manually from https://golang.org/dl/${NC}"
        exit 1
        ;;
      *)
        echo -e "${RED}Unsupported operating system. Please install Go manually.${NC}"
        exit 1
        ;;
    esac
    
    # Verify Go installation
    if command -v go >/dev/null 2>&1; then
      GO_VERSION=$(go version | awk '{print $3}')
      echo -e "✅ Go installed successfully: ${GO_VERSION}"
    else
      echo -e "${RED}❌ Failed to install Go. Please install it manually.${NC}"
      exit 1
    fi
  fi
  
  # Check for git
  if command -v git >/dev/null 2>&1; then
    GIT_VERSION=$(git --version | awk '{print $3}')
    echo -e "✅ Git is installed: ${GIT_VERSION}"
  else
    echo -e "${RED}❌ Git is not installed${NC}"
    echo -e "Installing Git..."
    
    # Determine OS
    OS="$(uname -s)"
    case "${OS}" in
      Linux*)
        if command -v apt-get >/dev/null 2>&1; then
          $USE_SUDO apt-get update
          $USE_SUDO apt-get install -y git
        elif command -v yum >/dev/null 2>&1; then
          $USE_SUDO yum install -y git
        elif command -v pacman >/dev/null 2>&1; then
          $USE_SUDO pacman -S git --noconfirm
        else
          echo -e "${RED}Unsupported Linux distribution. Please install Git manually.${NC}"
          exit 1
        fi
        ;;
      Darwin*)
        if command -v brew >/dev/null 2>&1; then
          brew install git
        else
          echo -e "${RED}Homebrew not found. Please install Git manually.${NC}"
          exit 1
        fi
        ;;
      MINGW*|MSYS*|CYGWIN*)
        echo -e "${RED}Windows detected. Please install Git manually from https://git-scm.com/download/win${NC}"
        exit 1
        ;;
      *)
        echo -e "${RED}Unsupported operating system. Please install Git manually.${NC}"
        exit 1
        ;;
    esac
    
    # Verify Git installation
    if command -v git >/dev/null 2>&1; then
      GIT_VERSION=$(git --version | awk '{print $3}')
      echo -e "✅ Git installed successfully: ${GIT_VERSION}"
    else
      echo -e "${RED}❌ Failed to install Git. Please install it manually.${NC}"
      exit 1
    fi
  fi
  
  echo -e "${GREEN}All required tools are installed!${NC}"
}

# Clone or update repository
setup_repo() {
  echo -e "\n${BOLD}Setting up Fort-Go repository...${NC}"
  
  # Define installation directory
  INSTALL_DIR="$HOME/.fort-go"
  
  # Check if directory exists
  if [ -d "$INSTALL_DIR" ]; then
    echo "Fort-Go repository already exists. Updating..."
    cd "$INSTALL_DIR"
    git pull
  else
    echo "Cloning Fort-Go repository..."
    git clone https://github.com/princebabou/fort-go.git "$INSTALL_DIR"
    cd "$INSTALL_DIR"
  fi
  
  echo -e "${GREEN}Repository setup complete!${NC}"
}

# Build the application
build_app() {
  echo -e "\n${BOLD}Building Fort-Go...${NC}"
  
  cd "$HOME/.fort-go"
  go build -o fort ./cmd/fort
  
  echo -e "${GREEN}Build complete!${NC}"
}

# Install the application
install_app() {
  echo -e "\n${BOLD}Installing Fort-Go system-wide...${NC}"
  
  # Create bin directory if it doesn't exist
  BIN_DIR="/usr/local/bin"
  if [ ! -d "$BIN_DIR" ]; then
    $USE_SUDO mkdir -p "$BIN_DIR"
  fi
  
  # Copy the binary
  $USE_SUDO cp "$HOME/.fort-go/fort" "$BIN_DIR/fort"
  
  # Make it executable
  $USE_SUDO chmod +x "$BIN_DIR/fort"
  
  echo -e "${GREEN}Installation complete!${NC}"
}

# Setup shell completion
setup_completion() {
  echo -e "\n${BOLD}Setting up shell completion...${NC}"
  
  # Detect shell
  SHELL_TYPE=$(basename "$SHELL")
  
  case "$SHELL_TYPE" in
    bash)
      COMPLETION_FILE="$HOME/.bash_completion"
      COMPLETION_DIR="$HOME/.bash_completion.d"
      
      # Create completion directory if it doesn't exist
      mkdir -p "$COMPLETION_DIR"
      
      # Generate completion script
      "$HOME/.fort-go/fort" completion bash > "$COMPLETION_DIR/fort.bash"
      
      # Add to .bash_completion if it doesn't exist
      if [ ! -f "$COMPLETION_FILE" ] || ! grep -q "fort.bash" "$COMPLETION_FILE"; then
        echo "# Fort-Go completion" >> "$COMPLETION_FILE"
        echo "source $COMPLETION_DIR/fort.bash" >> "$COMPLETION_FILE"
      fi
      
      echo "Bash completion installed. Please restart your shell or run 'source $COMPLETION_FILE'."
      ;;
    zsh)
      COMPLETION_DIR="$HOME/.zsh/completion"
      
      # Create completion directory if it doesn't exist
      mkdir -p "$COMPLETION_DIR"
      
      # Generate completion script
      "$HOME/.fort-go/fort" completion zsh > "$COMPLETION_DIR/_fort"
      
      # Add to .zshrc if it doesn't exist
      if ! grep -q "completion/_fort" "$HOME/.zshrc"; then
        echo "# Fort-Go completion" >> "$HOME/.zshrc"
        echo "fpath=($COMPLETION_DIR \$fpath)" >> "$HOME/.zshrc"
        echo "autoload -U compinit" >> "$HOME/.zshrc"
        echo "compinit" >> "$HOME/.zshrc"
      fi
      
      echo "Zsh completion installed. Please restart your shell or run 'source ~/.zshrc'."
      ;;
    *)
      echo "Shell completion for $SHELL_TYPE is not supported. Skipping."
      ;;
  esac
  
  echo -e "${GREEN}Shell completion setup complete!${NC}"
}

# Run the install
main() {
  check_required
  setup_repo
  build_app
  install_app
  setup_completion
  
  echo -e "\n${BOLD}${GREEN}Fort-Go has been successfully installed!${NC}"
  echo -e "You can now use the '${BOLD}fort${NC}' command from anywhere."
  echo 
  echo -e "Example commands:"
  echo -e "  ${BOLD}fort scan -t example.com${NC}              # Perform a full scan"
  echo -e "  ${BOLD}fort exploit -t example.com${NC}           # Attempt safe exploitation"
  echo -e "  ${BOLD}fort report -i results.json -f pdf${NC}    # Generate a PDF report"
  echo 
  echo -e "For more information, run: ${BOLD}fort --help${NC}"
}

# Run the main function
main 