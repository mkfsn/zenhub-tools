package zenhub

import (
	"sort"
	"time"
)

type TeamSprint struct {
	startTime  time.Time
	endTime    time.Time
	allIssues  issues
	userIssues map[string]issues
}

func NewTeamSprint(members []string, sprint Sprint) TeamSprint {
	allIssues := make(issues, 0)
	userIssues := make(map[string]issues, 0)

	for _, issue := range sprint.SprintIssues.Issues() {
		allIssues = append(allIssues, issue.Issue)

		if len(issue.Issue.Assignees.Nodes) < 1 {
			continue
		}

		assignee := issue.Issue.Assignees.Nodes[0].Login

		if !isMember(members, assignee) {
			continue
		}

		userIssues[assignee] = append(userIssues[assignee], issue.Issue)
	}

	return TeamSprint{
		startTime:  sprint.StartAt,
		endTime:    sprint.EndAt,
		allIssues:  allIssues,
		userIssues: userIssues,
	}
}

type UserBurndownStoryPoints struct {
	byDate map[string][]int
}

func (t TeamSprint) UserBurndownIssues() map[string]issues {
	result := make(map[string]issues)

	for user, issues := range t.userIssues {
		for _, issue := range issues {
			if issue.ClosedAt.IsZero() {
				continue
			}
			result[user] = append(result[user], issue)
		}
	}

	for user, issues := range result {
		sort.SliceStable(issues, func(i, j int) bool {
			return issues[i].ClosedAt.Before(issues[j].ClosedAt)
		})
		result[user] = issues
	}

	return result
}

func (t TeamSprint) issueDates() []string {
	days := int64(t.endTime.Sub(t.startTime).Hours()/24.0) + 1

	result := make([]string, 0, days)

	for ts := t.startTime; ts.Before(t.endTime); ts = ts.Add(time.Hour * 24) {
		result = append(result, toDate(ts))
	}

	return result
}

type issues []Issue

func (i issues) burndownStoryPoints(startTime, endTime time.Time) []float64 {
	var totalPoint float64

	burnPointByDate := make(map[string]float64)

	for _, issue := range i {
		point := issue.Estimate.Value
		totalPoint += point

		if issue.ClosedAt.IsZero() {
			continue
		}

		burnPointByDate[toDate(issue.ClosedAt)] += point
	}

	days := int64(endTime.Sub(startTime).Hours()/24.0) + 1

	result := make([]float64, 0, days)

	for t := startTime; t.Before(endTime); t = t.Add(time.Hour * 24) {
		burnPoint := burnPointByDate[toDate(t)]

		totalPoint -= burnPoint

		result = append(result, totalPoint)
	}

	return result
}

func (i issues) burndownIssues(startTime, endTime time.Time) []int64 {
	var totalIssues int64

	burnIssueByDate := make(map[string]int64)

	for _, issue := range i {
		totalIssues += 1

		if issue.ClosedAt.IsZero() {
			continue
		}

		burnIssueByDate[toDate(issue.ClosedAt)] += 1
	}

	days := int64(endTime.Sub(startTime).Hours()/24.0) + 1

	result := make([]int64, 0, days)

	for t := startTime; t.Before(endTime); t = t.Add(time.Hour * 24) {
		burnPoint := burnIssueByDate[toDate(t)]

		totalIssues -= burnPoint

		result = append(result, totalIssues)
	}

	return result
}

func isMember(members []string, user string) bool {
	for _, member := range members {
		if member == user {
			return true
		}
	}
	return false
}

func toDate(t time.Time) string {
	return t.In(time.Local).Format("2006-01-02")
}
