package cron

const (
	// Every
	//
	// Execute every time
	//
	// # Example
	//  @every 10s (execute once every 10 seconds)
	//  @every 10m30s (execute once every 10 minutes and 30 seconds)
	//  @every 1h30m (execute once every 1 hour and 30 minutes)
	Every = "@every"

	// Daily
	//
	// Execute once a day, using the 24-hour format, support AndChar
	//
	// # Example
	//  @daily 5:00 (execute once at 5:00)
	//  @daily 23:00 (execute once at 23:00)
	//  @daily 0:00,12:00 (execute both at 0:00 and 12:00)
	Daily = "@daily"

	// Weekly
	//
	// Execute once a week, support AndChar, BetweenChar
	//
	// # Example
	//  @weekly 1 5:00 (execute once at Monday 5:00)
	//  @weekly 3 0:00,12:00 (execute both at Wednesday 0:00 and 12:00)
	//  @weekly 1-5 5:00 (execute from Monday to Friday at 5:00)
	//  @weekly 3,7 5:00 (execute at 5:00 on Wednesday and Sunday)
	Weekly = "@weekly"

	// Monthly
	//
	// Execute once a month, support AndChar, BetweenChar, LastNDomPrefix
	//
	// # Example
	//  @monthly 1 5:00 (execute once at 5:00 on the 1st day of the month)
	//  @monthly L 5:00 (execute once at 5:00 on the last day of the month)
	//  @monthly 15,L 23:00 (execute both at 23:00 on the 15th and the last day of the month)
	//  @monthly 1-5 5:00 (execute from the 1st to the 5th day of the month at 5:00)
	//  @monthly 5,10,15,20,25,30 5:00,23:00 (execute at 5:00 and 23:00 on the 5th, 10th, 15th, 20th, 25th, and 30th days of the month)
	Monthly = "@monthly"
)

const (
	// DescriptionPrefix
	//
	// The prefix character for descriptions
	DescriptionPrefix = "@"

	// LastNDomPrefix
	//
	// The character for the last N days of the month
	//
	// # Example
	//  L: the last day of the month
	//  L-1: the last second day of the month
	//
	// # Notice
	//  If the days of the month are less than the value of N, the first day of the month will be used
	LastNDomPrefix = "L"

	// AndChar
	//
	// The character for the `BOTH` conditions
	//
	// # Example
	//  0:00,12:00 (execute both at 0:00 and 12:00)
	AndChar = ","

	// BetweenChar
	//
	// The character for multiple `AND` range conditions
	//
	// # Example
	//  0-10 (from 0 to 10, like [0, 10])
	BetweenChar = "-"
)
