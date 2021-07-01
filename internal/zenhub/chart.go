package zenhub

import (
	"io"
	"os"
	"time"

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
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Story Points"}),
		charts.WithInitializationOpts(opts.Initialization{Theme: "shine"}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
	)
	line = line.SetXAxis(t.sprintDates())

	line.Overlap(
		t.expectedBurnRate(t.totalStoryPoints()),
		weekendChart(t.totalStoryPoints()),
		t.userStoryPointsBurndownChart(),
		t.totalStoryPointsBurndownChart(),
	)

	return line
}

func (t TeamSprint) issuesBurndownChart() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Issues"}),
		charts.WithInitializationOpts(opts.Initialization{Theme: "shine"}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
	)
	line = line.SetXAxis(t.sprintDates())

	line.Overlap(
		t.expectedBurnRate(float64(t.totalIssues())),
		weekendChart(float64(t.totalIssues())),
		t.userIssuesBurndownChart(),
		t.totalIssuesBurndownChart(),
	)

	return line
}

func (t TeamSprint) expectedBurnRate(initialValue float64) *charts.Line {
	line := charts.NewLine()

	workingDays := t.workingDays()
	totalStoryPoints := initialValue
	storyPointPerDay := totalStoryPoints / float64(workingDays)

	items := make([]opts.LineData, 0)
	for ts := t.startTime; ts.Before(t.endTime); ts = ts.Add(time.Hour * 24) {
		items = append(items, opts.LineData{Value: totalStoryPoints, Symbol: "none", SymbolSize: 0})
		if ts.Weekday() != time.Sunday && ts.Weekday() != time.Saturday {
			totalStoryPoints -= storyPointPerDay
		}
	}
	items = append(items, opts.LineData{Value: 0})
	line = line.AddSeries("", items,
		charts.WithItemStyleOpts(opts.ItemStyle{Color: "gray", Opacity: 0}),
		charts.WithLineStyleOpts(opts.LineStyle{Type: "dotted", Color: "gray"}),
	)

	return line
}

func weekendChart(maxValue interface{}) *charts.Line {
	line := charts.NewLine()

	line.AddSeries("", []opts.LineData{
		{Value: nil},
		{Value: maxValue, Symbol: "none", SymbolSize: 0},
		{Value: maxValue, Symbol: "none", SymbolSize: 0},
		{Value: maxValue, Symbol: "none", SymbolSize: 0},
		{Value: nil},
		{Value: nil},
		{Value: nil},
		{Value: nil},
		{Value: maxValue, Symbol: "none", SymbolSize: 0},
		{Value: maxValue, Symbol: "none", SymbolSize: 0},
		{Value: maxValue, Symbol: "none", SymbolSize: 0},
		{Value: nil},
		{Value: nil},
	},
		charts.WithItemStyleOpts(opts.ItemStyle{Color: "gray", Opacity: 0}),
		charts.WithLineStyleOpts(opts.LineStyle{Color: "#e5e5e5", Width: 0, Opacity: 0}),
	).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{Step: true}),
			charts.WithAreaStyleOpts(opts.AreaStyle{Opacity: 0.2}),
		)

	return line
}

func (t TeamSprint) userStoryPointsBurndownChart() *charts.Line {
	line := charts.NewLine()

	for user, issues := range t.userIssues {
		points := issues.burndownStoryPoints(t.startTime.Add(-time.Hour*24), t.endTime)
		items := make([]opts.LineData, 0, len(points))

		for _, point := range points {
			items = append(items, opts.LineData{Value: point})
		}

		line = line.AddSeries(user, items, charts.WithLabelOpts(opts.Label{Show: true}))
	}

	return line
}

func (t TeamSprint) totalStoryPointsBurndownChart() *charts.Line {
	line := charts.NewLine()

	points := t.allIssues.burndownStoryPoints(t.startTime.Add(-time.Hour*24), t.endTime)
	items := make([]opts.LineData, 0, len(points))

	for _, point := range points {
		items = append(items, opts.LineData{Value: point})
	}

	line.SetXAxis(t.issueDates()).
		AddSeries("Total", items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{Show: true}),
			charts.WithAreaStyleOpts(opts.AreaStyle{Opacity: 0.2}),
		)

	return line
}

func (t TeamSprint) userIssuesBurndownChart() *charts.Line {
	line := charts.NewLine()

	line = line.SetXAxis(t.issueDates())

	for user, issues := range t.userIssues {
		points := issues.burndownIssues(t.startTime.Add(-time.Hour*24), t.endTime)
		items := make([]opts.LineData, 0, len(points))

		for _, point := range points {
			items = append(items, opts.LineData{Value: point})
		}

		line = line.AddSeries(user, items, charts.WithLabelOpts(opts.Label{Show: true}))
	}

	return line
}

func (t TeamSprint) totalIssuesBurndownChart() *charts.Line {
	line := charts.NewLine()

	points := t.allIssues.burndownIssues(t.startTime.Add(-time.Hour*24), t.endTime)
	items := make([]opts.LineData, 0, len(points))

	for _, point := range points {
		items = append(items, opts.LineData{Value: point})
	}

	line.SetXAxis(t.issueDates()).
		AddSeries("Total", items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{Show: true}),
			charts.WithAreaStyleOpts(opts.AreaStyle{Opacity: 0.2}),
		)

	return line
}
