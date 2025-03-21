package job

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"go-hexagonal/util/log"
)

// Job represents a scheduled job
type Job interface {
	// Name returns the job name
	Name() string
	// Run executes the job
	Run(ctx context.Context) error
}

// Scheduler manages scheduled jobs
type Scheduler struct {
	cron     *cron.Cron
	jobs     map[string]Job
	jobSpecs map[string]string
	mu       sync.RWMutex
}

// DefaultJobTimeout is the default timeout for job execution
const DefaultJobTimeout = 5 * time.Minute

// NewScheduler creates a new job scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		cron:     cron.New(cron.WithSeconds()),
		jobs:     make(map[string]Job),
		jobSpecs: make(map[string]string),
	}
}

// AddJob adds a job to the scheduler
func (s *Scheduler) AddJob(spec string, job Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.Name()]; exists {
		return fmt.Errorf("job %s already exists", job.Name())
	}

	_, err := s.cron.AddFunc(spec, func() {
		ctx, cancel := context.WithTimeout(context.Background(), DefaultJobTimeout)
		defer cancel()

		start := time.Now()
		log.Logger.Info("Starting job",
			zap.String("job", job.Name()),
			zap.String("spec", spec),
		)

		if err := job.Run(ctx); err != nil {
			log.Logger.Error("Job failed",
				zap.String("job", job.Name()),
				zap.Error(err),
				zap.Duration("duration", time.Since(start)),
			)
			return
		}

		log.Logger.Info("Job completed",
			zap.String("job", job.Name()),
			zap.Duration("duration", time.Since(start)),
		)
	})

	if err != nil {
		return fmt.Errorf("failed to add job %s: %w", job.Name(), err)
	}

	s.jobs[job.Name()] = job
	s.jobSpecs[job.Name()] = spec
	return nil
}

// RemoveJob removes a job from the scheduler
func (s *Scheduler) RemoveJob(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[name]; !exists {
		return fmt.Errorf("job %s not found", name)
	}

	delete(s.jobs, name)
	delete(s.jobSpecs, name)
	return nil
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Logger.Info("Job scheduler started")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Logger.Info("Job scheduler stopped")
}

// ListJobs returns a list of all registered jobs
func (s *Scheduler) ListJobs() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make(map[string]string)
	for name, spec := range s.jobSpecs {
		jobs[name] = spec
	}
	return jobs
}
