package zenhub

import (
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (t TeamSprint) DrawBurndownChart(filename string) error {
	page := components.NewPage()
	page.AddCharts(
		t.storyPointsBurndownChart(),
		t.issuesBurndownChart(),
	)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	return page.Render(io.MultiWriter(f))
}

func (t TeamSprint) storyPointsBurndownChart() *charts.Line {
	line := t.userStoryPointsBurndownChart()
	line.Overlap(t.totalStoryPointsBurndownChart())
	return line
}

func (t TeamSprint) issuesBurndownChart() *charts.Line {
	line := t.userIssuesBurndownChart()
	line.Overlap(t.totalIssuesBurndownChart())
	return line
}

func (t TeamSprint) userStoryPointsBurndownChart() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Story Points",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme: "shine",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
		}),
	)

	line = line.SetXAxis(t.issueDates())

	for user, issues := range t.userIssues {
		points := issues.burndownStoryPoints(t.startTime, t.endTime)
		items := make([]opts.LineData, 0, len(points))

		for _, point := range points {
			items = append(items, opts.LineData{Value: point})
		}

		line = line.AddSeries(user, items, charts.WithLabelOpts(
			opts.Label{
				Show: true,
			}))
	}

	return line
}

func (t TeamSprint) totalStoryPointsBurndownChart() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "area options",
		}),
	)

	points := t.allIssues.burndownStoryPoints(t.startTime, t.endTime)
	items := make([]opts.LineData, 0, len(points))

	for _, point := range points {
		items = append(items, opts.LineData{Value: point})
	}

	line.SetXAxis(t.issueDates()).AddSeries("Total", items).
		SetSeriesOptions(
			charts.WithLabelOpts(
				opts.Label{
					Show: true,
				}),
			charts.WithAreaStyleOpts(
				opts.AreaStyle{
					Opacity: 0.2,
				}),
		)

	return line
}

func (t TeamSprint) userIssuesBurndownChart() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Issues",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme: "shine",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
		}),
	)

	line = line.SetXAxis(t.issueDates())

	for user, issues := range t.userIssues {
		points := issues.burndownIssues(t.startTime, t.endTime)
		items := make([]opts.LineData, 0, len(points))

		for _, point := range points {
			items = append(items, opts.LineData{Value: point})
		}

		line = line.AddSeries(user, items, charts.WithLabelOpts(
			opts.Label{
				Show: true,
			}))
	}

	return line
}

func (t TeamSprint) totalIssuesBurndownChart() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "area options",
		}),
	)

	points := t.allIssues.burndownIssues(t.startTime, t.endTime)
	items := make([]opts.LineData, 0, len(points))

	for _, point := range points {
		items = append(items, opts.LineData{Value: point})
	}

	line.SetXAxis(t.issueDates()).AddSeries("Total", items).
		SetSeriesOptions(
			charts.WithLabelOpts(
				opts.Label{
					Show: true,
				}),
			charts.WithAreaStyleOpts(
				opts.AreaStyle{
					Opacity: 0.2,
				}),
		)

	return line
}
