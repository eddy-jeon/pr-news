#!/bin/bash
#
# PR News - Configuration Module
# ==============================
#
# 설정 파일 로드 및 의존성 체크를 담당합니다.
#
# Functions:
#   config::load      - 설정 파일 로드 (기본: ~/.pr-news.conf)
#   config::check_deps - 필수 의존성 확인 (gh, claude, jq)
#
# Configuration Variables:
#   DAYS                    - 조회 기간 (일)
#   THRESHOLD_FILES         - 큰 PR 기준 파일 수
#   THRESHOLD_CHANGES       - 큰 PR 기준 변경 라인 수
#   INCLUDE_REVIEW_COMMENTS - 리뷰 코멘트 포함 여부
#   BOT_FILTER              - 제외할 봇 목록 (쉼표 구분)
#

# Default values
DEFAULT_DAYS=7
DEFAULT_THRESHOLD_FILES=10
DEFAULT_THRESHOLD_CHANGES=500
DEFAULT_INCLUDE_REVIEW_COMMENTS=true
DEFAULT_BOT_FILTER="coderabbitai,snyk-io-us,dependabot,github-actions,codecov"

# Load configuration
config::load() {
  local config_file="${1:-$HOME/.pr-news.conf}"

  # Set defaults
  DAYS="${DAYS:-$DEFAULT_DAYS}"
  THRESHOLD_FILES="${THRESHOLD_FILES:-$DEFAULT_THRESHOLD_FILES}"
  THRESHOLD_CHANGES="${THRESHOLD_CHANGES:-$DEFAULT_THRESHOLD_CHANGES}"
  INCLUDE_REVIEW_COMMENTS="${INCLUDE_REVIEW_COMMENTS:-$DEFAULT_INCLUDE_REVIEW_COMMENTS}"
  BOT_FILTER="${BOT_FILTER:-$DEFAULT_BOT_FILTER}"

  # Load from config file if exists
  if [[ -f "$config_file" ]]; then
    source "$config_file"
  fi

  # Export for subshells
  export DAYS THRESHOLD_FILES THRESHOLD_CHANGES INCLUDE_REVIEW_COMMENTS BOT_FILTER
}

# Check required dependencies
config::check_deps() {
  local missing=()

  command -v gh &>/dev/null || missing+=("gh")
  command -v claude &>/dev/null || missing+=("claude")
  command -v jq &>/dev/null || missing+=("jq")

  if [[ ${#missing[@]} -gt 0 ]]; then
    echo "Missing required dependencies: ${missing[*]}"
    echo "Please install them first."
    return 1
  fi

  # Check gh auth
  if ! gh auth status &>/dev/null; then
    echo "GitHub CLI not authenticated. Run: gh auth login"
    return 1
  fi

  return 0
}
