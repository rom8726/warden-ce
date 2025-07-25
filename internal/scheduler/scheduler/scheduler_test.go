package scheduler

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

type TestCronConfig struct {
}

func (*TestCronConfig) Schedule() string {
	return "*/1 * * * * *"
}

type testJob struct {
	name        string
	runCount    *int32
	errToReturn error
	doneCh      chan struct{}
}

func (j *testJob) Name() string {
	return j.name
}

func (j *testJob) Run(ctx context.Context) error {
	atomic.AddInt32(j.runCount, 1)
	if j.doneCh != nil {
		j.doneCh <- struct{}{}
	}
	return j.errToReturn
}

func TestScheduler_RunJob(t *testing.T) {
	runCount := int32(0)
	doneCh := make(chan struct{}, 1)
	job := &testJob{name: "test", runCount: &runCount, doneCh: doneCh}
	testCfg := &TestCronConfig{}

	s := New()
	err := s.Register(job, testCfg)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = s.Start(ctx)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	select {
	case <-doneCh:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("Job was not executed in time")
	}

	if atomic.LoadInt32(&runCount) == 0 {
		t.Error("Job Run was not called")
	}

	s.Stop(context.Background())
}

func TestScheduler_JobErrorLogging(t *testing.T) {
	runCount := int32(0)
	doneCh := make(chan struct{}, 1)
	job := &testJob{name: "errjob", runCount: &runCount, errToReturn: errors.New("fail"), doneCh: doneCh}
	cfg := &TestCronConfig{}

	s := New()
	err := s.Register(job, cfg)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = s.Start(ctx)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	select {
	case <-doneCh:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("Job with error was not executed in time")
	}

	if atomic.LoadInt32(&runCount) == 0 {
		t.Error("Job Run was not called for error job")
	}

	s.Stop(context.Background())
}

func TestScheduler_StopIdempotent(t *testing.T) {
	s := New()
	ctx := context.Background()
	err := s.Start(ctx)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if err := s.Stop(ctx); err != nil {
		t.Fatalf("First Stop failed: %v", err)
	}
	if err := s.Stop(ctx); err != nil {
		t.Fatalf("Second Stop failed: %v", err)
	}
}
