# Zenhub Tools


## burndown report

This leverages [go-rod/rod](https://github.com/go-rod/rod) to open the 
[Zenhub Burndown Report page for sprint](https://blog.zenhub.com/tracking-sprint-progress-with-scrum-burndown-charts/)
via Chrome Browser. 

This script need to know which [Chrome Profile](https://chromium.googlesource.com/chromium/src/+/HEAD/docs/user_data_dir.md) 
to use to open the browser so that it can hijack the response data from Zenhub's graphql API (which seems not in GA).

### How to build

```bash
go build ./cmd/burndown-report
```

### How to use

This script requires three arguments:

- members: A CSVed GitHub account string for listing only issues belong to given user.
- zenhub-sprint-url: The link to [Zenhub Burndown Report page for sprint](https://blog.zenhub.com/tracking-sprint-progress-with-scrum-burndown-charts/)
- profile-dir: The Chrome profile directory

After running this:

```bash
./burndown-report \
  -members mkfsn,ionicc,pannpers,HowJMay \
  -zenhub-sprint-url 'https://app.zenhub.com/workspaces/XXXXXX/reports/burndown?milestoneId=YYYYYY&entity=sprints' \
  -profile-dir 'Profile 1'
```

you should see a `burndown-chart.html` file which has two charts:

- Story Points
- Issues

![image](https://user-images.githubusercontent.com/667169/123982868-f70d7000-d9f5-11eb-9c80-6b7fcc69ef76.png)
