package util

import "sync/atomic"

type AtomicBool struct{ flag int32 }

func (b *AtomicBool) Store(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32(&(b.flag), int32(i))
}

func (b *AtomicBool) Load() bool {
	return atomic.LoadInt32(&(b.flag)) != 0
}
