#!/bin/bash
#
# PR News - TUI Module (Graceful Fallback)
# ========================================
#
# 터미널 UI 렌더링을 담당합니다.
# gum > fzf > bash 순으로 fallback하여 어떤 환경에서도 동작합니다.
#
# Functions:
#   ui::header   - 앱 헤더 출력
#   ui::section  - 섹션 제목 출력
#   ui::spinner  - 로딩 스피너 (gum) 또는 텍스트
#   ui::success  - 성공 메시지 (✓)
#   ui::error    - 에러 메시지 (✗)
#   ui::info     - 정보 메시지 (●)
#   ui::choose   - 대화형 선택 UI
#   ui::box      - 박스 스타일 출력
#   ui::progress - 진행 상황 표시
#
# Environment Detection:
#   HAS_GUM - gum 설치 여부 (1/0)
#   HAS_FZF - fzf 설치 여부 (1/0)
#
# Vim Keybindings (fzf):
#   j/k - up/down navigation
#   Ctrl+d/u - half page scroll
#   Ctrl+f/b - full page scroll
#   / - toggle search
#

# ANSI Color codes
readonly C_RESET='\033[0m'
readonly C_BOLD='\033[1m'
readonly C_DIM='\033[2m'
readonly C_CYAN='\033[36m'
readonly C_GREEN='\033[32m'
readonly C_YELLOW='\033[33m'
readonly C_MAGENTA='\033[35m'
readonly C_RED='\033[31m'

# Detect available tools
HAS_GUM=$(command -v gum &>/dev/null && echo 1 || echo 0)
HAS_FZF=$(command -v fzf &>/dev/null && echo 1 || echo 0)

# Print app header
ui::header() {
  if [[ $HAS_GUM -eq 1 ]]; then
    gum style \
      --border rounded \
      --border-foreground 6 \
      --padding "0 2" \
      --margin "1 0" \
      "$(gum style --bold '🗞️  PR News') - GitHub PR Learning Tool"
  else
    echo ""
    echo -e "${C_CYAN}╭──────────────────────────────────────────────────────────╮${C_RESET}"
    echo -e "${C_CYAN}│${C_RESET} ${C_BOLD}🗞️  PR News${C_RESET} - GitHub PR Learning Tool                   ${C_CYAN}│${C_RESET}"
    echo -e "${C_CYAN}╰──────────────────────────────────────────────────────────╯${C_RESET}"
    echo ""
  fi
}

# Print section header
ui::section() {
  local title="$1"
  echo ""
  if [[ $HAS_GUM -eq 1 ]]; then
    gum style --foreground 5 --bold "█ $title"
  else
    echo -e "${C_MAGENTA}█${C_RESET} ${C_BOLD}$title${C_RESET}"
  fi
  echo -e "${C_DIM}$(printf '─%.0s' {1..50})${C_RESET}"
}

# Show spinner while running command
ui::spinner() {
  local message="$1"
  shift
  local cmd="$@"

  if [[ $HAS_GUM -eq 1 ]]; then
    gum spin --spinner dot --title "$message" -- bash -c "$cmd"
  else
    echo -en "${C_DIM}▸${C_RESET} ${message}... "
    if eval "$cmd" &>/dev/null; then
      echo -e "${C_GREEN}✓${C_RESET}"
    else
      echo -e "${C_RED}✗${C_RESET}"
      return 1
    fi
  fi
}

# Success message
ui::success() {
  local message="$1"
  if [[ $HAS_GUM -eq 1 ]]; then
    gum style --foreground 2 "✓ $message"
  else
    echo -e "${C_GREEN}✓${C_RESET} $message"
  fi
}

# Error message
ui::error() {
  local message="$1"
  if [[ $HAS_GUM -eq 1 ]]; then
    gum style --foreground 1 "✗ $message"
  else
    echo -e "${C_RED}✗${C_RESET} $message"
  fi
}

# Info message
ui::info() {
  local message="$1"
  if [[ $HAS_GUM -eq 1 ]]; then
    gum style --foreground 3 "● $message"
  else
    echo -e "${C_YELLOW}●${C_RESET} $message"
  fi
}

# Interactive choice with search (reads from stdin)
# Supports vim keybindings: j/k for navigation, / for search
ui::choose() {
  local header="$1"

  if [[ $HAS_GUM -eq 1 ]]; then
    # gum filter: 검색 + 선택 지원
    gum filter --header "$header" --placeholder "Type to search..." --height=15
  elif [[ $HAS_FZF -eq 1 ]]; then
    # fzf with vim keybindings
    fzf --prompt="$header " \
        --height=15 \
        --reverse \
        --bind='j:down,k:up,ctrl-j:down,ctrl-k:up' \
        --bind='ctrl-d:half-page-down,ctrl-u:half-page-up' \
        --bind='ctrl-f:page-down,ctrl-b:page-up' \
        --bind='/:toggle-search'
  else
    # Fallback to bash select - read stdin into array
    local items=()
    while IFS= read -r line; do
      items+=("$line")
    done

    echo "$header" >&2
    PS3="Select number: "
    select item in "${items[@]}"; do
      if [[ -n "$item" ]]; then
        echo "$item"
        break
      fi
    done < /dev/tty
  fi
}

# Print styled box (for summary output)
ui::box() {
  local title="$1"
  local content="$2"

  if [[ $HAS_GUM -eq 1 ]]; then
    echo "$content" | gum style \
      --border rounded \
      --border-foreground 6 \
      --padding "1 2" \
      --margin "1 0"
  else
    local width=70
    echo -e "${C_CYAN}╭$(printf '─%.0s' $(seq 1 $((width-2))))╮${C_RESET}"
    echo -e "${C_CYAN}│${C_RESET} ${C_BOLD}$title${C_RESET}$(printf ' %.0s' $(seq 1 $((width-4-${#title}))))${C_CYAN}│${C_RESET}"
    echo -e "${C_CYAN}├$(printf '─%.0s' $(seq 1 $((width-2))))┤${C_RESET}"
    while IFS= read -r line; do
      local padding=$((width - 4 - ${#line}))
      [[ $padding -lt 0 ]] && padding=0
      echo -e "${C_CYAN}│${C_RESET} $line$(printf ' %.0s' $(seq 1 $padding)) ${C_CYAN}│${C_RESET}"
    done <<< "$content"
    echo -e "${C_CYAN}╰$(printf '─%.0s' $(seq 1 $((width-2))))╯${C_RESET}"
  fi
}

# Progress indicator with bar
ui::progress() {
  local current="$1"
  local total="$2"
  local message="$3"
  local width=30
  local percent=$((current * 100 / total))
  local filled=$((current * width / total))
  local empty=$((width - filled))

  # Build progress bar
  local bar=""
  for ((i=0; i<filled; i++)); do bar+="█"; done
  for ((i=0; i<empty; i++)); do bar+="░"; done

  if [[ $HAS_GUM -eq 1 ]]; then
    # Clear line and print
    echo -ne "\r\033[K"
    gum style --foreground 6 "[$current/$total] $bar ${percent}% $message"
  else
    # Clear line and print
    echo -ne "\r\033[K"
    echo -ne "${C_CYAN}[$current/$total]${C_RESET} ${C_DIM}${bar}${C_RESET} ${percent}% $message"
  fi
}

# Progress indicator - finish (new line)
ui::progress_done() {
  echo ""
}

# Fetch indicator with spinner animation
ui::fetch_start() {
  local message="$1"
  if [[ $HAS_GUM -eq 1 ]]; then
    echo -e "${C_DIM}▸${C_RESET} $message"
  else
    echo -en "${C_DIM}▸${C_RESET} $message "
  fi
}

# Spinner frames for fetch animation
readonly SPINNER_FRAMES=('⠋' '⠙' '⠹' '⠸' '⠼' '⠴' '⠦' '⠧' '⠇' '⠏')

# Animated fetch indicator (call in loop)
ui::fetch_tick() {
  local frame_idx="$1"
  local message="$2"
  local frame="${SPINNER_FRAMES[$((frame_idx % ${#SPINNER_FRAMES[@]}))]}"

  echo -ne "\r\033[K${C_CYAN}${frame}${C_RESET} $message"
}

# Fetch complete
ui::fetch_done() {
  local message="$1"
  echo -e "\r\033[K${C_GREEN}✓${C_RESET} $message"
}
