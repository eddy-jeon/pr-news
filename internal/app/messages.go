package app

import "github.com/eddy/pr-news/internal/github"

// Messages for async operations

type ReposLoadedMsg struct {
	Repos []string
	Err   error
}

type PRsFetchedMsg struct {
	PRs []github.PR
	Err error
}

type PRDataCollectedMsg struct {
	Data      string
	Current   int
	Total     int
	StartDate string // 가장 오래된 PR 날짜
	EndDate   string // 가장 최근 PR 날짜
}

type SummaryDoneMsg struct {
	Summary string
	Err     error
}

type ErrMsg struct{ Err error }

// ClearCopyMsg clears the "Copied!" feedback after a delay
type ClearCopyMsg struct{}
