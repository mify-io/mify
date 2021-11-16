package util

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
	"errors"
	"context"

	"github.com/briandowns/spinner"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

func CreateWaitSpinner() *spinner.Spinner {
	return spinner.New(spinner.CharSets[9], 200*time.Millisecond)
}


type ProgressBar struct {
	bar *mpb.Bar
	progress *mpb.Progress
	statusFunc decor.DecorFunc
	total int64

	spinCycle int
	spinnerChars []string
	mtx sync.Mutex

}

func NewProgressBar(statusFunc decor.DecorFunc) *ProgressBar {
	return &ProgressBar{
		statusFunc: statusFunc,
		spinnerChars: []string{"/", "-", "\\", "|"},
	}
}

func (pb *ProgressBar) Create(total int64) {
	if pb.bar != nil {
		return
	}
	if total >= 0 {
		pb.total = total
	}
	// if pb.total == 0 {
		// return
	// }
	pb.create(total)
}

func (pb *ProgressBar) create(total int64) {
	pb.progress = mpb.New(mpb.WithWidth(64))
	pb.bar = pb.progress.Add(pb.total,
		mpb.NewBarFiller(mpb.BarStyle().Lbound("[").Filler("=").Tip(">").Padding("-").Rbound("]")),
		mpb.PrependDecorators(
			decor.OnComplete(
				decor.Any(func(s decor.Statistics) string {
					return pb.statusFunc(s)
				}),
				"done",
			),
		),
		mpb.AppendDecorators(decor.Any(func(s decor.Statistics) string {
			return fmt.Sprintf("%d/%d", s.Current, s.Total)
		})),
	)
	pb.bar.SetTotal(pb.total, false)
}

func (pb *ProgressBar) Abort() {
	if pb.bar == nil {
		return
	}
	pb.bar.Abort(true)
}

func (pb *ProgressBar) IncTotal() {
	pb.total += 1
	if pb.bar == nil {
		return
	}
	if pb.bar.Completed() {
		pb.create(0)
	}
	pb.bar.SetTotal(pb.total, false)
}

func (pb *ProgressBar) ResetTotal() {
	if pb.bar == nil {
		return
	}
	pb.total = 0
	pb.bar.SetTotal(0, false)
}

func (pb *ProgressBar) Increment() {
	if pb.bar == nil {
		return
	}
	pb.bar.Increment()
}

func (pb *ProgressBar) Wait() {
	if pb.bar == nil {
		return
	}
	pb.bar.SetTotal(pb.total, true)
	pb.progress.Wait()
}

func (pb *ProgressBar) Spinner() string {
	pb.mtx.Lock()
	defer pb.mtx.Unlock()
	char := pb.spinnerChars[pb.spinCycle]
	pb.spinCycle += 1
	if pb.spinCycle == len(pb.spinnerChars) {
		pb.spinCycle = 0
	}
	return char
}

func ShowJobError(pool *JobPool, jerr JobError) {
	if errors.Is(jerr.Err, context.Canceled) {
		return
	}
	fmt.Printf("task %s error: %s\n", jerr.Name, jerr.Err)
	logFile, err := os.Open(pool.GetJobLogPath(jerr.Name))
	if err != nil {
		fmt.Printf("failed to read job %s log: %s", jerr.Name, err)
	}
	fmt.Printf("\nfull log:\n")
	io.Copy(os.Stderr, logFile)
	logFile.Close()
}
