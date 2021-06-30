package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/mkfsn/zenhub-tools/internal/browser"
	"github.com/mkfsn/zenhub-tools/internal/zenhub"
)

func main() {
	var membersCSV string
	var zenhubSprintUrl string
	var profileDir string
	flag.StringVar(&membersCSV, "members", "mkfsn,ionicc,pannpers,HowJMay", "for filtering the assignee of the issue")
	flag.StringVar(&zenhubSprintUrl, "zenhub-sprint-url", "", "URL to zenhub sprint burndown report page")
	flag.StringVar(&profileDir, "profile-dir", "Profile 1", "The Chrome profile directory")
	flag.Parse()

	if zenhubSprintUrl == "" {
		fmt.Println("empty zenhub sprint url")
		fmt.Println("usage:")
		flag.PrintDefaults()
		return
	}

	b, err := browser.GetRawSprintIssues(context.Background(), zenhubSprintUrl, profileDir)
	if err != nil {
		fmt.Printf("failed to get raw sprint issues: %s\n", err)
		return
	}

	data, err := zenhub.DecodeFromRawSprintIssues(b)
	if err != nil {
		fmt.Printf("failed to decode sprint issues: %#v\n", err)
		return
	}

	c := zenhub.NewTeamSprint(strings.Split(membersCSV, ","), data[0].Data.Node)

	for user, issues := range c.UserBurndownIssues() {
		fmt.Printf("%s\n", user)
		for _, issue := range issues {
			fmt.Printf(" - (%v) %s %s\n",
				issue.Estimate.Value,
				issue.ClosedAt.In(time.Local).Format(time.RFC3339),
				issue.HTMLURL,
			)
		}
	}

	if err := c.DrawBurndownChart("burndown.html"); err != nil {
		fmt.Printf("failed to draw burndown chart: %s\n", err)
		return
	}
}
