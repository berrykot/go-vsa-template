package health

import (
	"go-vsa-template/internal/config"
	"go-vsa-template/internal/infrastructure/scheduler"

	"github.com/rs/zerolog"
)

type JobHandler struct {
	logger     zerolog.Logger
	healthCron string
}

func NewJobHandler(l zerolog.Logger, cfg *config.Config) *JobHandler {
	return &JobHandler{
		logger:     l,
		healthCron: cfg.Cron.HealthCron,
	}
}

func (h *JobHandler) Register(s *scheduler.Scheduler) {
	_, err := s.Cron.AddFunc(h.healthCron, h.doJob)
	if err != nil {
		panic(err)
	}
}

func (h *JobHandler) doJob() {
	h.logger.Info().Msg("Executing quarterly report...")
}
