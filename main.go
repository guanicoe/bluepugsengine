package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/evilsocket/islazy/tui"
	"github.com/guanicoe/bluepugsengine/session"
)

// Var decleration of the flags for user interaction
var (
	TargetURL   = flag.String("u", "", "Target URL to scrap.")
	TimeOut     = flag.Int("t", 300, "Timeout after wish the scrap will stop. This will still return the result.")
	HardLimit   = flag.Int("l", 2000, "Maximum number of pages to scrap.")
	DomainScope = flag.String("d", "", "Scope of domains to scrap.")
	NWorkers    = flag.Int("w", 100, "Number of pugs to go to work.")
	FileName    = flag.String("o", "", "Name of output file.")
	StartZMQ    = flag.Bool("z", false, "Starting ZMQ server, this is for automation needs")
	PortZMQ     = flag.Int("p", 5155, "Specify the port you want the zmq server to listen on. Only works with -z.")
)

type FlagArguments struct {
	TimeOut     int
	TargetURL   string
	HardLimit   int
	DomainScope string
	NWorkers    int
	FileName    string
	StartZMQ    bool
}

/*
Main function of Blue Pugs. In order to use colour feedback, we first check if tui can work.
We then pull the flags and set them in a struct. A switch condition checks if we are running the program
in local mode or using a ZMQ socket.
*/
func main() {

	if !tui.Effects() {
		fmt.Printf("\n\nWARNING: This terminal does not support colours, view will be very limited.\n\n")
	}

	session.ASCIIArt()

	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	flagValues := FlagArguments{
		TimeOut:     *TimeOut,
		TargetURL:   *TargetURL,
		HardLimit:   *HardLimit,
		DomainScope: *DomainScope,
		NWorkers:    *NWorkers,
		FileName:    *FileName,
		StartZMQ:    *StartZMQ,
	}

	//  Socket to talk to server
	switch {
	case *StartZMQ:
		fmt.Println(tui.Wrap(tui.BACKYELLOW+tui.FOREBLACK, `WARNING: ZMQ flag detected, if you want to run Blue Pugs in the command line, remove the "-z" flag`))
		session.ZmqServer(*PortZMQ)
	case *TargetURL == "":
		log.Fatal(tui.Wrap(tui.BACKRED, "No target urls where given. Use -h or --help for help."))
	default:
		session.RunInTerminal(session.FlagArguments(flagValues))
	}

}
