package threading

import "sync"

func DoUnderLock(mtx *sync.Mutex, work func() error) error {
	mtx.Lock()
	defer mtx.Unlock()
	return work()
}
