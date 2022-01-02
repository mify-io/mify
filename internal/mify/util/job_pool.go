package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/chebykinn/mify/internal/mify/core"
	"github.com/vbauerster/mpb/v7/decor"
)

type JobFunc func(*core.Context) error

type Job struct {
	Name string
	Func JobFunc
}

type JobError struct {
	Name string
	Err  error
}

func (j JobError) Error() string {
	return fmt.Sprintf("job %s error: %s", j.Name, j.Err)
}

type JobPool struct {
	waitGroup   *sync.WaitGroup
	jobChan     chan Job
	stopChan    chan struct{}
	errChan     chan JobError
	jobsQueue   []Job
	runningJobs sync.Map

	progressBar *ProgressBar

	isError AtomicBool
	logDir  string

	ctx *core.Context
}

func (p *JobPool) addJob(j Job) {
	p.runningJobs.Store(j.Name, struct{}{})
}

func (p *JobPool) delJob(j Job) {
	p.waitGroup.Done()
	p.progressBar.Increment()
	p.runningJobs.Delete(j.Name)
}

func (p *JobPool) worker(n int) {
	for {
		var job Job
		select {
		case job = <-p.jobChan:
		case <-p.stopChan:
			p.stopChan <- struct{}{}
			return
		}

		var logFile *os.File
		var err error
		p.addJob(job)

		if !p.isError.Load() {
			path := p.GetJobLogPath(job.Name)
			logFile, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			logger := log.New(logFile, "", 0)
			pCtx := &core.Context{
				Ctx:    p.ctx.Ctx,
				Logger: logger,
			}
			if err == nil {
				err = job.Func(pCtx)
			}
			if err := logFile.Close(); err != nil {
				panic(fmt.Errorf("failed to close file %s %w", path, err))
			}
		}

		if err != nil {
			p.progressBar.Abort()
			jobErr := JobError{
				Name: job.Name,
				Err:  err,
			}
			p.isError.Store(true)
			select {
			case p.errChan <- jobErr:
			default:
			}
		}
		p.delJob(job)
	}
}

func (p *JobPool) GetRunningJobs() []string {
	jobs := make([]string, 0)
	p.runningJobs.Range(func(key interface{}, value interface{}) bool {
		jobs = append(jobs, key.(string))
		return true
	})

	sort.Strings(jobs)
	return jobs
}

func (p *JobPool) updateStatus(s decor.Statistics) string {
	jobs := p.GetRunningJobs()
	return p.progressBar.Spinner() + " running: [" + strings.Join(jobs, ", ") + "] "
}

func NewJobPool(ctx *core.Context, cacheDir string, numWorkers int) (*JobPool, error) {
	job_ch := make(chan Job, numWorkers)

	var wg sync.WaitGroup
	p := &JobPool{
		waitGroup: &wg,
		jobChan:   job_ch,
		stopChan:  make(chan struct{}),
		errChan:   make(chan JobError),
		ctx:       ctx,
	}

	p.progressBar = NewProgressBar(p.updateStatus)

	p.logDir = filepath.Join(cacheDir, "logs")
	err := os.MkdirAll(p.logDir, 0755)
	if err != nil {
		return nil, err
	}

	for i := 0; i < numWorkers; i++ {
		go p.worker(i)
	}

	return p, nil
}

func (p *JobPool) AddJob(j Job) {
	p.waitGroup.Add(1)
	p.progressBar.IncTotal()
	p.jobsQueue = append(p.jobsQueue, j)
}

func (p *JobPool) Run() error {
	p.progressBar.Create(-1)
	for len(p.jobsQueue) > 0 {
		j := p.jobsQueue[0]
		p.jobsQueue = p.jobsQueue[1:]
		p.jobChan <- j
	}

	wait := make(chan struct{})
	go func() {
		p.waitGroup.Wait()
		wait <- struct{}{}
	}()

	select {
	case <-wait:
	case err := <-p.errChan:
		return err
	}

	return nil
}

func (p *JobPool) GetJobLogPath(name string) string {
	return filepath.Join(p.logDir, "job-"+name+".log")
}

func (p *JobPool) ClosePool() {
	if p.isError.Load() {
		return
	}
	p.waitGroup.Wait()
	p.progressBar.Wait()
	p.stopChan <- struct{}{}
}
