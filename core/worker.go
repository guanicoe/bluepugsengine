package core

import (
	"context"
)

type worker struct {
	CTX       context.Context
	ID        int
	Work      chan workRequest
	PugsQueue chan chan workRequest
	Result    chan workOutput
}

type workRequest struct {
	UnscrapedURL string
}

type workOutput struct {
	FoundLinks  []string
	FoundEmails []emailSource
	InitialURL  string
}

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func workerNew(ctx context.Context, id int, pugsQueue chan chan workRequest, result chan workOutput) worker {

	worker := worker{
		CTX:       ctx,
		ID:        id,
		Work:      make(chan workRequest),
		PugsQueue: pugsQueue,
		Result:    result,
	}

	return worker
}

// This method "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *worker) Start() {

	go func() {
		for {
			w.PugsQueue <- w.Work
			select {
			case <-w.CTX.Done():
				return
			case work := <-w.Work:
				links, emails := scrap(work.UnscrapedURL)
				r := workOutput{
					FoundLinks:  links,
					FoundEmails: emails,
					InitialURL:  work.UnscrapedURL,
				}
				w.Result <- r
			}
		}
	}()
}
