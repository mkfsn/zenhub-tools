package zenhub

import (
	"encoding/json"
	"time"
)

type User struct {
	// Typename  string `json:"__typename"`
	AvatarURL string `json:"avatarUrl"`
	ID        string `json:"id"`
	Login     string `json:"login"`
}

type Assignees struct {
	// Typename string `json:"__typename"`
	Nodes []User `json:"nodes"`
}

type Label struct {
	// Typename string `json:"__typename"`
	Color string `json:"color"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

type Labels struct {
	// Typename string  `json:"__typename"`
	Nodes []Label `json:"nodes"`
}

type Milestone struct {
	// Typename string `json:"__typename"`
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Estimate struct {
	// Typename string `json:"__typename"`
	// TODO
	Value float64 `json:"value"`
}

type Repository struct {
	// Typename string `json:"__typename"`
	GhID  int    `json:"ghId"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner User   `json:"owner"`
}

type PipelineIssue struct {
	// Typename           string    `json:"__typename"`
	ID                 string    `json:"id"`
	LatestTransferTime time.Time `json:"latestTransferTime"`
	Pipeline           struct {
		Typename string `json:"__typename"`
		ID       string `json:"id"`
		Name     string `json:"name"`
	} `json:"pipeline"`
	Priority interface{} `json:"priority"`
}

type Epic interface{}

type ParentEpics struct {
	// Typename string `json:"__typename"`
	Nodes []Epic `json:"nodes"`
}

type Issue struct {
	// Typename      string        `json:"__typename"`
	Assignees     Assignees     `json:"assignees"`
	ClosedAt      time.Time     `json:"closedAt"`
	Epic          Epic          `json:"epic"`
	Estimate      Estimate      `json:"estimate"`
	HTMLURL       string        `json:"htmlUrl"`
	ID            string        `json:"id"`
	Labels        Labels        `json:"labels"`
	Milestone     Milestone     `json:"milestone"`
	Number        int           `json:"number"`
	ParentEpics   ParentEpics   `json:"parentEpics"`
	PipelineIssue PipelineIssue `json:"pipelineIssue"`
	PullRequest   bool          `json:"pullRequest"`
	Repository    Repository    `json:"repository"`
	State         string        `json:"state"`
	Title         string        `json:"title"`
	User          User          `json:"user"`
}

type SprintIssue struct {
	// Typename  string    `json:"__typename"`
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Issue     Issue     `json:"issue"`
}

type SprintIssues struct {
	// Typename string        `json:"__typename"`
	Nodes []SprintIssue `json:"nodes"`
}

func (s *SprintIssues) Issues() []SprintIssue {
	result := make([]SprintIssue, 0, len(s.Nodes))

	for _, issue := range s.Nodes {
		if issue.Issue.PullRequest {
			continue
		}
		result = append(result, issue)
	}

	return result
}

type Sprint struct {
	// Typename     string       `json:"__typename"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	SprintIssues SprintIssues `json:"sprintIssues"`
	StartAt      time.Time    `json:"startAt"`
	EndAt        time.Time    `json:"endAt"`
	State        string       `json:"state"`
}

type ResponseData struct {
	Node Sprint `json:"node"`
}

type Response []struct {
	Data ResponseData `json:"data"`
}

func DecodeFromRawSprintIssues(b []byte) (Response, error) {
	var data Response
	err := json.Unmarshal(b, &data)
	return data, err
}
