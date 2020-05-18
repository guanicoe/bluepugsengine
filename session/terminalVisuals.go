package session

import (
	"fmt"

	"github.com/evilsocket/islazy/tui"
)

//ASCIIArt simple visual printing of the logo in terminal. Ran immediatly in main hence exported
func ASCIIArt() {
	asciiArt :=
		`

  $$$$$$$$  $$$$   $$$$  $$$$   $$$$$$     $$$$$$$  $$$$  $$$$    $$$$$$    $$$$$
  $$$$$$$$$ $$$$   $$$$  $$$$  $$$$$$$$    $$$$$$$$ $$$$  $$$$  $$$$$$$$$  $$$$$$$
  $$$$ $$$$ $$$$   $$$$  $$$$ $$$$  $$$$   $$$$ $$$ $$$$  $$$$ $$$$       $$$$
  $$$$$$$$  $$$$   $$$$  $$$$ $$$$$$$$$$   $$$$$$$$ $$$$  $$$$ $$$$ $$$$$$ $$$$$$$
  $$$$ $$$$ $$$$   $$$$  $$$$ $$$$         $$$$$$$  $$$$  $$$$ $$$$  $$$$$     $$$$  $$
  $$$$$$$$$ $$$$$$$ $$$$$$$$   $$$$        $$$       $$$$$$$$  $$$$$$$$$$  $$$$$$$  $$$$
  $$$$$$$$  $$$$$$$  $$$$$$     $$$$$$     $$$        $$$$$$     $$$$$$$    $$$$$    $$

                                              GO - v0.3 guanicoe
  `
	fmt.Println(tui.Wrap(tui.BOLD+tui.BLUE+tui.DIM, asciiArt))
}
