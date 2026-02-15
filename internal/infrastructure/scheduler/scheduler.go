package scheduler

import (
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type Scheduler struct {
	Logger zerolog.Logger
	Cron   *cron.Cron
}

func New(l zerolog.Logger) *Scheduler {
	return &Scheduler{
		Logger: l,
		Cron:   cron.New(cron.WithSeconds()),
	}
}

func (s *Scheduler) Start() {
	s.Cron.Start()
	s.Logger.Info().Msg("Background scheduler started")

	entries := s.Cron.Entries()
	for _, entry := range entries {
		s.Logger.Info().
			Interface("next_run", entry.Next).
			Msg("Job")
	}

}

func (s *Scheduler) Stop() {
	s.Cron.Stop()
	s.Logger.Info().Msg("Background scheduler stopped")
}
