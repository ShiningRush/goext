package jobx

import (
	"context"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	DefaultSingleton = NewJobDemon()
)

func RegisterJobDesc(desc *JobDescriptor) {
	DefaultSingleton.RegisterJobDesc(desc)
}

func RegisterJob(name string, jobType JobType, jobFunc JobFunc) {
	DefaultSingleton.RegisterJob(name, jobType, jobFunc)
}

func Start() {
	DefaultSingleton.Start()
}

func Stop() {
	DefaultSingleton.Stop()
}

func Close() {
	DefaultSingleton.Close()
}

func RemoveJob(name string) {
	DefaultSingleton.RemoveJob(name)
}

func ClearJobs() {
	DefaultSingleton.ClearJobs()
}

func NewJobDemon() *JobDemon {
	return &JobDemon{}
}

type JobDemon struct {
	jobs []*JobDescriptor

	cancel     context.CancelFunc
	scheduleWg sync.WaitGroup
	runWg      sync.WaitGroup
}

type JobDescriptor struct {
	Name string
	Type JobType
	Func JobFunc
}

// Do starts the job and stops scheduling it when closeChan is closed.
//
// Deprecated: register the descriptor with JobDemon and call Start instead.
func (d *JobDescriptor) Do(ctx context.Context, closeChan chan struct{}) {
	ctx, cancel := context.WithCancel(ctx)
	var scheduleWg sync.WaitGroup
	var runWg sync.WaitGroup

	d.do(ctx, &scheduleWg, &runWg)

	go func() {
		select {
		case <-ctx.Done():
		case <-closeChan:
			cancel()
		}
	}()
	go func() {
		scheduleWg.Wait()
		runWg.Wait()
		cancel()
	}()
}

func (d *JobDescriptor) do(ctx context.Context, scheduleWg *sync.WaitGroup, runWg *sync.WaitGroup) {
	run := func() {
		runWg.Add(1)
		defer runWg.Done()
		d.Func(ctx)
	}

	if d.Type.Once != nil {
		if !d.Type.Once.AlwaysStart && d.Type.fired {
			return
		}

		scheduleWg.Add(1)
		go func() {
			defer scheduleWg.Done()

			timer := time.NewTimer(d.Type.Once.Delay)
			defer timer.Stop()

			select {
			case <-ctx.Done():
				return
			case <-timer.C:
			}

			select {
			case <-ctx.Done():
				return
			default:
			}

			run()
			d.Type.fired = true
		}()
		return
	}

	if d.Type.Interval != nil {
		scheduleWg.Add(1)
		go func() {
			defer scheduleWg.Done()

			for {
				timer := time.NewTimer(d.Type.Interval.Interval)
				select {
				case <-ctx.Done():
					timer.Stop()
					return
				case <-timer.C:
				}

				select {
				case <-ctx.Done():
					return
				default:
				}

				run()
			}
		}()
	}

	if d.Type.Cron != nil {
		c := cron.New(
			cron.WithParser(cron.NewParser(cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor)),
			cron.WithLocation(d.Type.Cron.location()),
		)
		if _, err := c.AddFunc(d.Type.Cron.Spec, func() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			run()
		}); err != nil {
			return
		}

		c.Start()
		scheduleWg.Add(1)
		go func() {
			defer scheduleWg.Done()

			<-ctx.Done()
			stopCtx := c.Stop()
			<-stopCtx.Done()
		}()
	}
}

type JobFunc func(ctx context.Context)

type JobType struct {
	Once     *OnceJobDesc
	Interval *IntervalJobDesc
	Cron     *CronJobDesc

	fired bool
}

type IntervalJobDesc struct {
	Interval time.Duration
}

type OnceJobDesc struct {
	Delay       time.Duration
	AlwaysStart bool
}

type CronJobDesc struct {
	Spec     string
	Location *time.Location
}

func (d *CronJobDesc) location() *time.Location {
	if d.Location != nil {
		return d.Location
	}
	return time.Local
}

func (d *JobDemon) RegisterJobDesc(jobDesc *JobDescriptor) {
	d.jobs = append(d.jobs, jobDesc)
}

func (d *JobDemon) RegisterJob(name string, jobType JobType, jobFunc JobFunc) {
	d.jobs = append(d.jobs, &JobDescriptor{
		Name: name,
		Type: jobType,
		Func: jobFunc,
	})
}

func (d *JobDemon) Start() {
	if d.cancel != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	d.cancel = cancel
	for _, job := range d.jobs {
		job.do(ctx, &d.scheduleWg, &d.runWg)
	}
}

func (d *JobDemon) Stop() {
	if d.cancel == nil {
		return
	}

	d.cancel()
	d.cancel = nil
}

func (d *JobDemon) Close() {
	d.Stop()
	d.scheduleWg.Wait()
	d.runWg.Wait()
}

// RemoveJob remove job by name
func (d *JobDemon) RemoveJob(name string) {
	for i, job := range d.jobs {
		if job.Name == name {
			d.jobs = append(d.jobs[:i], d.jobs[i+1:]...)
			return
		}
	}
}

// ClearJobs clear all jobs
func (d *JobDemon) ClearJobs() {
	d.jobs = nil
}
