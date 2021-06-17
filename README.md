## Auto Invest Sharesies NZ (Work in progress)

Currently, [Sharesies NZ](http://sharesies.nz) only support auto-invest for managed funds what creates the need for this application to apply dollar-cost averaging for companies on NZ Market Exchange.

### Scheduler
Linux crontab compatible instruction for executing orders

#### Predefined
Entry                  | Description                                | Equivalent To
-----                  | -----------                                | -------------
@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 1 1 *
@monthly               | Run once a month, midnight, first of month | 0 0 1 * *
@weekly                | Run once a week, midnight between Sat/Sun  | 0 0 * * 0
@daily (or @midnight)  | Run once a day, midnight                   | 0 0 * * *
@hourly                | Run once an hour, beginning of hour        | 0 * * * *

More information:
* [crontab.tech](https://crontab.tech/every-monday)
* [robfig/cron](https://pkg.go.dev/github.com/robfig/cron/v3@v3.0.0)


### Configuration 

`config/credentials.yml`
```yml
username: test@test.com
password: password
```


`config/auto_invest.yml`
```yml
scheduler: "0 8 * * MON" # Monday 8am

buy:
  - reference: Delegat Group # Only for log porpose
    id: 0545fbc5-b579-4944-9057-55d01849a493
    amount: 1.00
  - reference: ANZ # Only for log porpose
    id: 860a502e-d07c-435e-9dcc-7d4631a4ee21
    amount: 1.00
```