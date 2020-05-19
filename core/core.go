package core

import (
	"fmt"
	"time"

	"github.com/evilsocket/islazy/tui"
	log "github.com/sirupsen/logrus"
)

//JobParam contains the job parameters to be sent to the core
type JobParam struct {
	TimeOut     int
	TargetURL   string
	HardLimit   int
	DomainScope string
	NWorkers    int
	CheckEmails bool
}

//JsonOutput is a struct that will contain the output that will be formatted in json before being sent out or saved
type JsonOutput struct {
	TargetURL       string
	HardLimit       int
	NmbWorkers      int
	DomainScope     string
	NmbScraped      int
	NmbUniqueEmails int
	UniqueEmails    []string
	NmbEmails       int
	EmailList       []emailSource
	TimeStarted     time.Time
	TimeFinished    time.Time
	TimeDeltaMS     int64
}

/*
LaunchJob is the main subfunction that actually starts the job. This functions takes a pointer
of the input parameters either given by the zmq server or the command line interpreter.
This function just starts the producer and waits for its results.

param :
- fv pointer
- return output as JsonOutput struct
*/
func LaunchJob(p JobParam) JsonOutput {
	var msg string

	msg = tui.Green("Starting job!")
	log.Info(msg)

	startTime := time.Now()
	output, _ := startProducer(&p)
	endTime := time.Now()
	timeDelta := endTime.Sub(startTime)

	output.TimeStarted = startTime
	output.TimeFinished = endTime
	output.TimeDeltaMS = timeDelta.Milliseconds()

	log.Info(output.UniqueEmails)
	msg = tui.Green(fmt.Sprintf("Finished job at %s - It took %s", endTime, timeDelta))
	log.Info(msg)

	return output
}
