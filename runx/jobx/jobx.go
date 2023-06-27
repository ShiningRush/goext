package jobx

import (
	"context"
	"time"
)

var (
	DefaultSingleton = NewJobDemon()
)

func RegisterJob(name string, jobType JobType, jobFunc JobFunc) {
	DefaultSingleton.RegisterJob(name, jobType, jobFunc)
}

func Start() {
	DefaultSingleton.Start()
}

func Stop() {
	DefaultSingleton.Stop()
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
	jobs []*JobDesc

	closeChan chan struct{}
}

type JobDesc struct {
	Name string
	Type JobType
	Func JobFunc
}

func (d *JobDesc) Do(ctx context.Context, closeChan chan struct{}) {
	if d.Type.Once != nil {
		if !d.Type.Once.AlwaysStart && d.Type.fired {
			return
		}

		time.AfterFunc(d.Type.Once.Delay, func() {
			d.Func(ctx)
			d.Type.fired = true
		})
		return
	}

	if d.Type.Interval != nil {
		go func() {
			for {
				select {
				case <-closeChan:
					return
				case <-time.After(d.Type.Interval.Interval):
					d.Func(ctx)
				}
			}
		}()
	}
}

type JobFunc func(ctx context.Context)

type JobType struct {
	Once     *OnceJobDesc
	Interval *IntervalJobDesc

	fired bool
}

type IntervalJobDesc struct {
	Interval time.Duration
}

type OnceJobDesc struct {
	Delay       time.Duration
	AlwaysStart bool
}

func (d *JobDemon) RegisterJob(name string, jobType JobType, jobFunc JobFunc) {
	d.jobs = append(d.jobs, &JobDesc{
		Name: name,
		Type: jobType,
		Func: jobFunc,
	})
}

func (d *JobDemon) Start() {
	if d.closeChan != nil {
		return
	}

	d.closeChan = make(chan struct{})
	for _, job := range d.jobs {
		job.Do(context.Background(), d.closeChan)
	}
}

func (d *JobDemon) Stop() {
	if d.closeChan == nil {
		return
	}

	close(d.closeChan)
	d.closeChan = nil
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
