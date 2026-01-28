package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type PR struct {
	Number       int       `json:"number"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	Additions    int       `json:"additions"`
	Deletions    int       `json:"deletions"`
	ChangedFiles int       `json:"changedFiles"`
	MergedAt     time.Time `json:"mergedAt"`
	Author       struct {
		Login string `json:"login"`
	} `json:"author"`
	URL string `json:"url"`
}

const (
	ThresholdFiles   = 10
	ThresholdChanges = 500
)

// ListRepos returns accessible repositories (personal + org).
func ListRepos(limit int) ([]string, error) {
	out, err := exec.Command("gh", "repo", "list",
		"--limit", fmt.Sprintf("%d", limit),
		"--json", "nameWithOwner",
		"-q", ".[].nameWithOwner",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("listing repos: %w", err)
	}

	repos := make(map[string]bool)
	for _, r := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if r != "" {
			repos[r] = true
		}
	}

	// Org repos
	orgOut, _ := exec.Command("gh", "api", "user/orgs", "--jq", ".[].login").Output()
	for _, org := range strings.Split(strings.TrimSpace(string(orgOut)), "\n") {
		if org == "" {
			continue
		}
		oRepos, _ := exec.Command("gh", "repo", "list", org,
			"--limit", fmt.Sprintf("%d", limit),
			"--json", "nameWithOwner",
			"-q", ".[].nameWithOwner",
		).Output()
		for _, r := range strings.Split(strings.TrimSpace(string(oRepos)), "\n") {
			if r != "" {
				repos[r] = true
			}
		}
	}

	result := make([]string, 0, len(repos))
	for r := range repos {
		result = append(result, r)
	}
	return result, nil
}

// ListMergedPRs returns merged PRs in the given date range.
func ListMergedPRs(repo string, days int, baseBranch string) ([]PR, error) {
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	search := fmt.Sprintf("merged:>=%s", since)
	if baseBranch != "" {
		search += " base:" + baseBranch
	}

	out, err := exec.Command("gh", "pr", "list",
		"--repo", repo,
		"--state", "merged",
		"--search", search,
		"--limit", "50",
		"--json", "number,title,body,additions,deletions,changedFiles,mergedAt,author,url",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("listing PRs: %w", err)
	}

	var prs []PR
	if err := json.Unmarshal(out, &prs); err != nil {
		return nil, fmt.Errorf("parsing PRs: %w", err)
	}
	return prs, nil
}

// IsLargePR returns true if the PR exceeds size thresholds.
func IsLargePR(files, changes int) bool {
	return files > ThresholdFiles || changes > ThresholdChanges
}

// GetPRDiff returns the diff for a PR (capped at maxLines).
func GetPRDiff(repo string, number int, maxLines int) (string, error) {
	out, err := exec.Command("gh", "pr", "diff",
		fmt.Sprintf("%d", number),
		"--repo", repo,
	).Output()
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	return strings.Join(lines, "\n"), nil
}

// GetReviewComments returns review comments for a PR (bot-filtered).
func GetReviewComments(repo string, number int) (string, error) {
	out, err := exec.Command("gh", "pr", "view",
		fmt.Sprintf("%d", number),
		"--repo", repo,
		"--json", "comments",
		"--jq", `.comments[] | select(.authorAssociation != "NONE") | "- **\(.author.login)**: \(.body | split("\n")[0])"`,
	).Output()
	if err != nil {
		return "", nil // non-fatal
	}
	return strings.TrimSpace(string(out)), nil
}

// CollectPRData gathers formatted data for a single PR.
func CollectPRData(repo string, pr PR) string {
	var b strings.Builder
	changes := pr.Additions + pr.Deletions

	fmt.Fprintf(&b, "## PR #%d: %s\n", pr.Number, pr.Title)
	fmt.Fprintf(&b, "- Author: %s\n", pr.Author.Login)
	fmt.Fprintf(&b, "- Merged: %s\n", pr.MergedAt.Format("2006-01-02"))
	fmt.Fprintf(&b, "- Stats: +%d -%d (%d files)\n", pr.Additions, pr.Deletions, pr.ChangedFiles)
	fmt.Fprintf(&b, "- URL: %s\n", pr.URL)
	fmt.Fprintf(&b, "\n### Description\n%s\n", pr.Body)

	if !IsLargePR(pr.ChangedFiles, changes) {
		if diff, err := GetPRDiff(repo, pr.Number, 500); err == nil && diff != "" {
			fmt.Fprintf(&b, "\n### Code Changes (excerpt)\n```diff\n%s\n```\n", diff)
		}
	} else {
		b.WriteString("\n> Large PR - showing summary only\n")
	}

	if comments, err := GetReviewComments(repo, pr.Number); err == nil && comments != "" {
		fmt.Fprintf(&b, "\n### Review Comments\n%s\n", comments)
	}

	return b.String()
}
