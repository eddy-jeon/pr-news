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
	Data    string
	Current int
	Total   int
}

type SummaryDoneMsg struct {
	Summary string
	Err     error
}

type ErrMsg struct{ Err error }
