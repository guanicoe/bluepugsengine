package session

import (
	"encoding/json"
	"fmt"

	"github.com/evilsocket/islazy/tui"
	"github.com/guanicoe/bluepugsengine/core"
	zmq "github.com/pebbe/zmq4"
	log "github.com/prometheus/common/log"
)

/*
Utility function for creating the ZMQ socket
param:
 - port as int
return
	socket pointer
*/
func startServer(port int) *zmq.Socket {
	//  Socket to talk to clients
	responder, _ := zmq.NewSocket(zmq.REP)
	url := fmt.Sprintf("tcp://*:%v", port)
	responder.Bind(url)

	return responder
}

/*
ZmqServer for receiving the zmq packets over port 5515 by default. It first waits for a 1
to be able to acknowledge thatt it is ready to receive a job. It parses the data from
json. This loops until it receives a kill from client.
format for request:
{
	"TargetURL": "http://www.website.com/",
	"HardLimit": 1000,
	"DomainScope": "website",
	"NWorkers": 10,
}
param:
 - port as int
*/
func ZmqServer(port int) {
	responder := startServer(port)
	defer responder.Close()
	kill := true
	for kill {
		//  Wait for next request from client
		recv, _ := responder.Recv(0)

		switch {
		case recv == "1":
			responder.Send("1", 0)
		case recv == "kill":
			kill = false
			responder.Send("stoping", 0)
			break
		default:

			param := flagArguments{}
			if err := json.Unmarshal([]byte(recv), &param); err != nil {

				log.Warn(fmt.Sprintf("Received %s. Don't know what it means", recv))
				responder.Send("not understood", 0)
			} else {
				msg := fmt.Sprintf("Received job demand from client. Parameters ::: url - %s ; Limit - %v ; scope - %s ; #Pugs - %v.", param.TargetURL, param.HardLimit, param.DomainScope, param.NWorkers)
				log.Info(tui.Wrap(tui.BACKLIGHTBLUE, msg))

				p := core.JobParam{
					TimeOut:     param.TimeOut,
					TargetURL:   param.TargetURL,
					HardLimit:   param.HardLimit,
					DomainScope: param.DomainScope,
					NWorkers:    param.NWorkers,
				}

				result := core.LaunchJob(p)

				//  Send reply back to client
				jsonResult, _ := json.Marshal(result)

				log.Info("Sending result back to client.")

				// Sending results over tcp
				responder.Send(string(jsonResult), 0)
			}

			log.Info(tui.Wrap(tui.BACKGREEN, "Job done!"))
		}

	}
}
