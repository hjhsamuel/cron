package cron

import "github.com/robfig/cron/v3"

type (
	Cron    = cron.Cron
	Entry   = cron.Entry
	EntryID = cron.EntryID
	Job     = cron.Job
	FuncJob = cron.FuncJob
)

// New creates a scheduler with a seconds field in cron expressions.
func New() *Cron {
	return cron.New(cron.WithSeconds())
}
