package session

import (
	"fmt"
	"os"
	"time"

	"github.com/evilsocket/islazy/tui"
	"github.com/guanicoe/bluepugsengine/core"
	"github.com/jedib0t/go-pretty/table"
)

func printTable(result core.JsonOutput) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Email"})
	if len(result.UniqueEmails) > 1 {
		for i, m := range result.UniqueEmails {
			t.AppendRows([]table.Row{{i + 1, m}})
		}
	} else {
		t.AppendRows([]table.Row{{"!!! No emails found !!!"}})
	}
	t.AppendSeparator()
	t.AppendFooter(table.Row{"Total", fmt.Sprintf("%v email(s)", result.NmbUniqueEmails), fmt.Sprintf("%v url(s)", result.NmbScraped), time.Millisecond * time.Duration(result.TimeDeltaMS)})
	t.SetStyle(table.StyleColoredBright)
	t.Render()

}

func printParam(param core.JobParam) {
	paramText := fmt.Sprintf(
		`
###############################################################################

           URL:    %s
    Hard limit:    %v
         Scope:    %s
   Number pugs:    %v
	 `,
		param.TargetURL, param.HardLimit, param.DomainScope, param.NWorkers)
	fmt.Println(tui.Wrap(tui.BOLD+tui.YELLOW, paramText))

}

//ASCIIArt simple visual printing of the logo in terminal. Ran immediately in main hence exported
func ASCIIArt() {
	asciiArt :=
		`
$$$$$$$  $$$    $$$  $$$   $$$$$$     $$$$$$  $$$  $$$    $$$$$    $$$$$
$$$$$$$$ $$$    $$$  $$$  $$$$$$$$    $$$$$$$ $$$  $$$  $$$$$$$$  $$$$$$$
$$$ $$$$ $$$    $$$  $$$ $$$$  $$$$   $$$ $$$ $$$  $$$ $$$       $$$$
$$$$$$$  $$$    $$$  $$$ $$$$$$$$$$   $$$$$$$ $$$  $$$ $$$  $$$$  $$$$$$$
$$$ $$$$ $$$    $$$  $$$ $$$$         $$$$$$  $$$  $$$ $$$   $$$$     $$$$  $$
$$$$$$$$ $$$$$$$ $$$$$$   $$$$        $$       $$$$$$   $$$$$$$$  $$$$$$$  $$$$
$$$$$$$  $$$$$$$  $$$$     $$$$$$     $$        $$$$     $$$$$$    $$$$$    $$

                                            GO - v0.3 guanicoe`
	fmt.Println(tui.Wrap(tui.BOLD+tui.BLUE, asciiArt))
}
