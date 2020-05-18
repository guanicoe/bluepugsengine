package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/prometheus/common/log"
)

type jobData struct {
	WorkQueue       chan workRequest
	Result          chan workOutput
	paramPointer    *JobParam
	unscrapedURL    chan string
	unscrapedLen    int
	limit           int
	scrapedSent     []string
	scrapedSentLen  int
	scrapedRecv     []string
	scrapedRecvLen  int
	emailList       []emailSource
	emailListUnique []string
	continueProd    bool
	err             error
}

/*
Dispatcher function. It first creats a go routing for each worker and sends the appropriate
data. In then sents an embedded go routine which listens for available workers.
Once it gets one, it it grabs a work from the work queue and sends it to a worker
*/
func startDispatcher(ctx context.Context, n int, s *jobData) {

	// First, initialize the channel we are going to put the workers' work channels into.
	PugsQueue := make(chan chan workRequest, n)
	// Now, create all of our workers.
	for i := 0; i < n; i++ {
		fmt.Printf("\rStarting worker %v/%v", i+1, n)
		worker := workerNew(ctx, i+1, PugsQueue, s.Result)
		worker.Start()
	}
	fmt.Printf("\n")

	go func(ctx context.Context, s *jobData) {
		for {
			select {
			case <-ctx.Done():
				return
			case work := <-s.WorkQueue:

				go func(ctx context.Context) {
					select {
					case <-ctx.Done():
						return
					default:
						worker := <-PugsQueue
						worker <- work
					}

				}(ctx)

			}
		}
	}(ctx, s)
	return
}

/*
function ran in go routing from producer. This function first sends the initial
target url, it then listens for the processresult function for unscraped urls
and sends them to the workaueue buffered channel
param:
	- firstURL string initial url
	- ctx
	- s pointer to producer currated results
*/
func sendWork(ctx context.Context, firstURL string, s *jobData) {

	sendToPugs := func(l string) {
		s.unscrapedLen--
		work := workRequest{UnscrapedURL: l}
		s.WorkQueue <- work

	}

	sendToPugs(firstURL)
	s.scrapedSent = append(s.scrapedSent, firstURL)
	s.scrapedSentLen++
	for {
		select {
		case <-ctx.Done():
			return
		case l := <-s.unscrapedURL:

			sendToPugs(l)
			s.scrapedSent = append(s.scrapedSent, l)
			s.scrapedSentLen++

		}

	}

}

/*
Function to process result from workers. It first generates a map of email:sourceurl.
It then loops through the found urls and checks if the url was already visited.
if not it increments the unscraped var and sends the url to the channel to be added to the work queue
param:
	 - ctx
	 - r WorkOutput struct with worker result
	 - s pointer to the struct containing the curated results
	 - fv pointer to the original parameters
*/
func processResult(ctx context.Context, r workOutput, s *jobData) {

	for _, emailMap := range r.FoundEmails {
		s.emailList = append(s.emailList, emailMap)
	}
	for _, l := range r.FoundLinks {

		if newValidURL(l, s, s.paramPointer.DomainScope) {
			if s.scrapedSentLen <= s.limit {
				s.unscrapedLen++
				s.unscrapedURL <- l

			}
		}
	}

}

/*
Core function producerm it creats the context, and sends the jobs to workers.
It checks what urls were visited, and creats a struct to keep all the data in one place.
The context is sent to the workers in order to stop them when the work is done.

param:
	- fv pointer
	- return output in JsonOutput and error
*/
func startProducer(param *JobParam) (JsonOutput, error) {

	//Declare context for current job
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(param.TimeOut))

	// local struct contaning all the data. only pointers are sent across
	s := jobData{
		WorkQueue: make(chan workRequest, 1000),
		Result:    make(chan workOutput, 1000),

		paramPointer: param,

		unscrapedURL:    make(chan string),
		unscrapedLen:    1,
		scrapedSent:     []string{},
		scrapedSentLen:  0,
		scrapedRecv:     []string{},
		scrapedRecvLen:  0,
		emailList:       []emailSource{},
		emailListUnique: []string{},
		continueProd:    true,
		err:             nil,
	}

	var wg sync.WaitGroup //Local wait group to wait for the main loop that is sent as a go routing

	startDispatcher(ctx, s.paramPointer.NWorkers, &s) //Calling the dispatcher function that will start the workers and distribute the work

	log.Debug("Sending initial job ", s.paramPointer)
	start := time.Now() //get time as start to give a few seconds wait before timing out if nothing is received from workers

	go sendWork(ctx, s.paramPointer.TargetURL, &s) // starting a goroutiing of the sendWork
	wg.Add(1)

	//main loop as gorouting to listen and clean worker result and send new urls for scraping
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		for {

			//switch condition to verify that we have not hit a wall
			switch {
			case s.scrapedRecvLen >= s.paramPointer.HardLimit:
				// in case we hit hard limit, we exist the loop
				return
			case s.scrapedRecvLen == s.scrapedSentLen && time.Since(start) > time.Duration(20)*time.Second:
				// if we received the same number as we sent and that we waited a bit to make sure the workers are not working, we exit
				return
			default:
				//If we havn't hit a stop condition, we go and listen for workers
				select {
				case <-ctx.Done(): // in case the context is canceled
					return
				case r := <-s.Result: // listening on result channel for workers' output
					s.scrapedRecv = append(s.scrapedRecv, r.InitialURL) //We add turl the worker scraped to our receive slice
					s.scrapedRecvLen++                                  // and we increment the length
					msg := fmt.Sprintf("Unscraped %v | scrapedRecv %v | scrapedSent %v | emails found %v               ", s.unscrapedLen, s.scrapedRecvLen, s.scrapedSentLen, len(s.emailList))
					fmt.Printf("\r%s", msg) // lazy printing of progression on terminal

					processResult(ctx, r, &s) // we send result to the process function

				case <-time.After(2 * time.Second): // we loop every 2 seconds in order not to block on Result
				}
			}
		}
	}(ctx, &wg)
	wg.Wait()
	cancel() //We cancel the context when the go loop returns

	// TODO: check if necessary. Quick for loop to purge the work buffered queue
	for len(s.WorkQueue) > 0 {
		<-s.WorkQueue
	}

	fmt.Println()

	setUniqueEmail(&s) //Sending result pointer to get unique email list

	msg := fmt.Sprintf("\nProducer has finished. Scraped %v urls | Found %v emails | unique %v", s.scrapedRecvLen, len(s.emailList), len(s.emailListUnique))
	log.Info(msg)

	//Prepare the output json struct to send as return
	output := JsonOutput{
		TargetURL:       param.TargetURL,
		HardLimit:       param.HardLimit,
		NmbWorkers:      param.NWorkers,
		DomainScope:     param.DomainScope,
		NmbScraped:      s.scrapedSentLen,
		NmbUniqueEmails: len(s.emailListUnique),
		UniqueEmails:    s.emailListUnique,
		NmbEmails:       len(s.emailList),
		EmailList:       s.emailList,
	}

	return output, s.err //return result and error

}
