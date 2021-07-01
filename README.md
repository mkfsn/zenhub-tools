# Zenhub Tools


## burndown report

This leverages [go-rod/rod](https://github.com/go-rod/rod) to open the 
[Zenhub Burndown Report page for sprint](https://blog.zenhub.com/tracking-sprint-progress-with-scrum-burndown-charts/)
via Chrome Browser. 

### How to build

```bash
go build ./cmd/burndown-report
```

### How to use

This script requires the following arguments:

- zenhub-sprint-url: The link to [Zenhub Burndown Report page for sprint](https://blog.zenhub.com/tracking-sprint-progress-with-scrum-burndown-charts/)
- (Optional) members: A CSVed GitHub account string for listing only issues belong to given user.

After running this:

```bash
./burndown-report \
  -members mkfsn,ionicc,pannpers,HowJMay \
  -zenhub-sprint-url 'https://app.zenhub.com/workspaces/XXXXXX/reports/burndown?milestoneId=YYYYYY&entity=sprints'
```

you should see a `burndown-chart.html` file which has two charts:

- Story Points
- Issues

![image](https://user-images.githubusercontent.com/667169/123982868-f70d7000-d9f5-11eb-9c80-6b7fcc69ef76.png)
