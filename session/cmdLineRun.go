package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/guanicoe/bluepugsengine/core"

	"github.com/evilsocket/islazy/tui"
	log "github.com/sirupsen/logrus"
)

func cleanFileName(f *string) {
	if !strings.HasSuffix(*f, ".json") {
		*f = fmt.Sprintf("%s.json", *f)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//RunInTerminal is the function that will handle the work when the program is called from terminal
func RunInTerminal(fv FlagArguments) {

	if len(fv.FileName) > 0 {
		cleanFileName(&fv.FileName)

		if fileExists(fv.FileName) {
			reader := bufio.NewReader(os.Stdin)
			msg := fmt.Sprintf("File %s already exists, overwrite? [y/N]: ", fv.FileName)
			fmt.Print(tui.Wrap(tui.BACKYELLOW+tui.FOREBLACK, msg))
			text, _ := reader.ReadString('\n')
			if text != "y\n" || len(text) == 0 {
				log.Fatal(tui.Wrap(tui.BACKRED, "Quitting, choose another filename"))
			}
		}
	}

	param := core.JobParam{
		TimeOut:     fv.TimeOut,
		TargetURL:   fv.TargetURL,
		HardLimit:   fv.HardLimit,
		DomainScope: fv.DomainScope,
		NWorkers:    fv.NWorkers,
		CheckEmails: fv.CheckEmails,
	}

	result := core.LaunchJob(param)

	if len(fv.FileName) > 0 {

		jsonResult, _ := json.MarshalIndent(result, "", " ")
		err := ioutil.WriteFile(fv.FileName, jsonResult, 0644)

		switch {
		case err == nil:
			msg := fmt.Sprintf(`Data written to file "%s"`, fv.FileName)
			log.Info(tui.Wrap(tui.GREEN, msg))
		default:
			msg := fmt.Sprint("Error when writing to file: ", err)
			log.Panic(tui.Wrap(tui.BACKRED, msg))
		}
	}

}
