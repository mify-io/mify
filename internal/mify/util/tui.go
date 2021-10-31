package util

import (
	"time"

	"github.com/briandowns/spinner"
)

func CreateWaitSpinner() *spinner.Spinner {
	return spinner.New(spinner.CharSets[9], 200*time.Millisecond)
}
