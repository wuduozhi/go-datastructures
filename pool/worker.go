package pool

import "fmt"

type Worker interface {
	// Process will synchronously perform a job and return the result.
	Process(interface{}) interface{}

	// BlockUntilReady is called before each job is processed and must block the
	// calling goroutine until the Worker is ready to process the next job.
	BlockUntilReady()

	// Interrupt is called when a job is cancelled. The worker is responsible
	// for unblocking the Process implementation.
	Interrupt()

	// Terminate is called when a Worker is removed from the processing pool
	// and is responsible for cleaning up any held resources.
	Terminate()
}

//------------------------------------------------------------------------------
// closureWorker is a minimal Worker implementation that simply wraps a
// func(interface{}) interface{}
type closureWorker struct {
	processor func(interface{}) interface{}
}

func (w *closureWorker) Process(payload interface{}) interface{} {
	return w.processor(payload)
}

func (w *closureWorker) BlockUntilReady() {}
func (w *closureWorker) Interrupt()       {}
func (w *closureWorker) Terminate()       {}

//------------------------------------------------------------------------------
// callbackWorker is a minimal Worker implementation that attempts to cast
// each job into func() and either calls it if successful or returns
// ErrJobNotFunc.
type callbackWorker struct{}

func (w *callbackWorker) Process(payload interface{}) interface{} {
	f, ok := payload.(func())
	if !ok {
		return ErrJobNotFunc
	}
	f()
	return nil
}

func (w *callbackWorker) BlockUntilReady() {}
func (w *callbackWorker) Interrupt()       {}
func (w *callbackWorker) Terminate()       {}

// workRequest is a struct containing context representing a workers intention
// to receive a work payload.
type workRequest struct {
	// jobChan is used to send the payload to this worker.
	jobChan chan<- interface{}

	// retChan is used to read the result from this worker.
	retChan <-chan interface{}

	// interruptFunc can be called to cancel a running job. When called it is no
	// longer necessary to read from retChan.
	interruptFunc func()
}

type workerWrapper struct {
	name          string
	worker        Worker
	interruptChan chan struct{}

	// reqChan is NOT owned by this type, it is used to send requests for work.
	reqChan chan<- workRequest

	// closeChan can be closed in order to cleanly shutdown this worker.
	closeChan chan struct{}

	// closedChan is closed by the run() goroutine when it exits.
	closedChan chan struct{}
}

func (w *workerWrapper) interrupt() {
	close(w.interruptChan)
	w.worker.Interrupt()
}

func (w *workerWrapper) run() {
	jobChan, retChan := make(chan interface{}), make(chan interface{})
	defer func() {
		w.worker.Terminate()
		close(retChan)
		close(w.closedChan)
	}()

	for {
		select {
		case w.reqChan <- workRequest{jobChan: jobChan, retChan: retChan, interruptFunc: w.interrupt}:
			select {
			case payload := <-jobChan:
				result := w.worker.Process(payload)
				select {
				case retChan <- result:
				case <-w.interruptChan:
					w.interruptChan = make(chan struct{})
				}

			case _, _ = <-w.interruptChan:
				w.interruptChan = make(chan struct{})
			}

		case <-w.closeChan:
			return
		}
	}

}

func (w *workerWrapper) stop() {
	close(w.closeChan)
}

func (w *workerWrapper) join() {
	<-w.closedChan
}

var workerNameIndex = 0

func newWorkerWrapper(reqChan chan<- workRequest, worker Worker) *workerWrapper {
	w := workerWrapper{
		name:          fmt.Sprintf("Worker-%v", workerNameIndex),
		worker:        worker,
		interruptChan: make(chan struct{}),
		reqChan:       reqChan,
		closeChan:     make(chan struct{}),
		closedChan:    make(chan struct{}),
	}
	workerNameIndex++

	go w.run()

	return &w
}

// NewCallback creates a new Pool of workers where workers cast the job payload
// into a func() and runs it, or returns ErrNotFunc if the cast failed.
func NewCallback(n int) *Pool {
	return New(n, func() Worker {
		return &callbackWorker{}
	})
}

// NewFunc creates a new Pool of workers where each worker will process using
// the provided func.
func NewFunc(n int, f func(interface{}) interface{}) *Pool {
	return New(n, func() Worker {
		return &closureWorker{
			processor: f,
		}
	})
}
