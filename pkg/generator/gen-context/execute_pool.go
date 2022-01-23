package gencontext

import "sync"

type ExecutePoolFactory struct {
}

type ExecutePool struct {
	wg      sync.WaitGroup
	results []chan error
}

func (ep *ExecutePool) Execute(work func() error) chan error {
	res := make(chan error)
	go func() {
		err := work()
		res <- err
	}()

	return res
}

func (ep *ExecutePool) EnqueExecution(work func() error) {
	ep.wg.Add(1)

	res := make(chan error)
	ep.results = append(ep.results, res)
	go func() {
		err := work()
		ep.wg.Done()
		res <- err
	}()
}

func (ep *ExecutePool) WaitAll() []error {
	ep.wg.Wait()
	errList := make([]error, 0, len(ep.results))
	for _, resChan := range ep.results {
		res := <-resChan
		if res != nil {
			errList = append(errList, res)
		}
	}

	ep.results = make([]chan error, 0)

	if len(errList) == 0 {
		return nil
	}

	return errList
}

func (ep *ExecutePoolFactory) NewPool() ExecutePool {
	return ExecutePool{}
}
