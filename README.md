# cron

A Go Cron parser library based on `github.com/robfig/cron/v3`. It provides more intuitive, semantic custom descriptor support on top of standard Cron expressions, such as `@every`, `@daily`, `@weekly`, and `@monthly`.

## Features

- **Semantic Descriptors**: Supports custom syntax that is easier to read than standard Cron.
- **Flexible Time Settings**: Supports specifying specific time points (24-hour format) in `@daily`, `@weekly`, and `@monthly`.
- **Multiple Time Points**: Allows using comma `,` to separate multiple time points.
- **Last N Days of Month Support**: Unique `L` and `L-n` syntax for targeting the last day or the Nth day from the end of the month.
- **Backward Compatibility**: Fully compatible with `robfig/cron/v3` standard Cron expressions.

## Installation

```bash
go get github.com/hjhsamuel/cron
```

## Quick Start

```go
package main

import (
	"fmt"
	"time"
	"github.com/hjhsamuel/cron"
)

func main() {
	// Parse custom expression
	sch, err := cron.Parse("@monthly L 23:00") // 23:00 on the last day of every month
	if err != nil {
		panic(err)
	}

	next := sch.Next(time.Now())
	fmt.Printf("Next execution time: %v\n", next)
	
	// Default using the seconds field
	cron.New()
}
```

## Supported Formats

### @every
Executes at the specified time interval.

- `@every 10s`: executes every 10 seconds
- `@every 1h30m`: executes every 1 hour and 30 minutes

### @daily
Executes at the specified time(s) every day. Supports comma-separated multiple times.

- `@daily 5:00`: executes at 5:00 every day
- `@daily 0:00,12:00`: executes at 0:00 and 12:00 every day

### @weekly
Executes at the specified day(s) and time(s) every week.

- `@weekly 1 5:00`: executes at 5:00 every Monday
- `@weekly 1-5 9:00`: executes at 9:00 from Monday to Friday
- `@weekly 0,6 10:00`: executes at 10:00 on Sunday and Saturday (both 0 and 7 represent Sunday)

### @monthly
Executes at the specified date(s) and time(s) every month. Supports special `L` syntax.

- `@monthly 1 5:00`: executes at 5:00 on the 1st of every month
- `@monthly 1,15,L 0:00`: executes at 0:00 on the 1st, 15th, and the last day of every month
- `@monthly L-1 23:59`: executes at 23:59 on the second to last day of every month

## Special Symbol Description

| Symbol | Description | Example |
| :--- | :--- | :--- |
| `,` | Conjunction (AND) | `0:00,12:00` (executes at both time points) |
| `-` | Range (BETWEEN) | `1-5` (from day 1 to day 5) |
| `L` | Last day of the month | `@monthly L 5:00` |
| `L-n` | The (n+1)th day from the end of the month | `L-1` means the second to last day, `L-2` means the third to last day |

**Notice**: For the `L-n` syntax, if the total number of days in the current month is less than `n`, it will fallback to executing on the first day (1st) of that month.

## License

[Apache 2.0](LICENSE)
