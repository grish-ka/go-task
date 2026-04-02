#!/bin/bash

# Colors
G='\033[0;32m'; B='\033[0;34m'; C='\033[0;36m'; NC='\033[0m'

# 1. Top-Level Environment Detection
UNAME_S=$(uname -s)
IS_WSL=false
grep -qi microsoft /proc/version 2>/dev/null && IS_WSL=true

# 2. Logic for Windows (Git Bash / MinGW / MSYS)
if [[ "$UNAME_S" == MINGW* ]] || [[ "$UNAME_S" == MSYS* ]] || [[ "$UNAME_S" == CYGWIN* ]]; then
    OS_ID="windows"
    OS_NAME="Windows (Git Bash/MinGW)"
    ICON="\ue8e5" #  Windows 11 Icon
    
# 3. Logic for Linux (WSL or Native)
elif [[ "$UNAME_S" == "Linux" ]]; then
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS_ID=$ID
        OS_NAME=$NAME
    else
        OS_ID="linux"; OS_NAME="Linux"
    fi

    # Distro-specific Icons
    case "$OS_ID" in
        ubuntu)   ICON="\uf31b" ;; # 
        fedora)   ICON="\uf30a" ;; # 
        arch)     ICON="\uf303" ;; # 󰣇
        debian)   ICON="\uf306" ;; # 
        "IotaOS") ICON="\uf444" ;; # Custom Icon for your OS!
        *)        ICON="\uf17c" ;; #  (Tux)
    esac

    # WSL Multi-Icon Override
    [ "$IS_WSL" = true ] && ICON="\ue8e5 $ICON"

# 4. Logic for macOS
elif [[ "$UNAME_S" == "Darwin" ]]; then
    OS_ID="darwin"
    OS_NAME="macOS"
    ICON="\uf179" # 🍎
fi

# 5. First-Time vs Update Detection
# Note: Windows Git Bash uses $HOME too!
CONFIG_DIR="$HOME/.config/go-task"
if [ ! -d "$CONFIG_DIR" ]; then
    MSG="Welcome to your first run of Go-Task!"
else
    MSG="Go-Task has been updated successfully."
fi

# 6. Direct-to-TTY Output
{
    printf "${B}%b${NC}  ${G}Go-Task identified: %s${NC}\n" "$ICON" "$OS_NAME"
    [ "$IS_WSL" = true ] && printf "${C}(Running on Windows 11 WSL)${NC}\n"
    printf "${G}%s${NC}\n" "$MSG"
    printf "Type ${B}go-task${NC} to begin.\n"
} > /dev/tty 2>/dev/null || true
