package scheduler

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rom8726/di"
)

var _ di.Servicer = (*Scheduler)(nil)

type Scheduler struct {
	cron    *cron.Cron
	jobs    map[cron.EntryID]Job
	jobsMu  sync.Mutex
	started bool
	startMu sync.Mutex
}

func New() *Scheduler {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}

	return &Scheduler{
		cron: cron.New(cron.WithSeconds(), cron.WithLocation(loc)),
		jobs: make(map[cron.EntryID]Job),
	}
}

func (s *Scheduler) Register(job Job, config Config) error {
	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()

	entryID, err := s.cron.AddFunc(config.Schedule(), func() {
		ctx := context.Background()
		if err := job.Run(ctx); err != nil {
			slog.Error("scheduler job error", "job", job.Name(), "err", err)
		}
	})
	if err != nil {
		return err
	}

	s.jobs[entryID] = job

	return nil
}

func (s *Scheduler) Start(context.Context) error {
	s.startMu.Lock()
	defer s.startMu.Unlock()
	if s.started {
		return nil
	}
	s.cron.Start()
	s.started = true

	slog.Info("scheduler started")

	return nil
}

func (s *Scheduler) Stop(ctx context.Context) error {
	s.startMu.Lock()
	defer s.startMu.Unlock()
	if !s.started {
		return nil
	}
	ctxDone := make(chan struct{})
	go func() {
		s.cron.Stop()
		close(ctxDone)
	}()
	select {
	case <-ctxDone:
	case <-ctx.Done():
	}
	s.started = false

	return nil
}
