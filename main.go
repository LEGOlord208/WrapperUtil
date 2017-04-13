package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/legolord208/stdutil"
)

const RESTART_TIMEOUT = 10
const RESTART_TIMEOUT_SHORT = 3
const NAME = "wrapperutil"

var COLOR = color.New(color.FgCyan)
var COLOR_ERR = color.New(color.FgRed, color.Bold)

var exit bool

type Packet struct {
	Restart bool
	Exit    bool
}

func main() {
	COLOR.Set()
	stdutil.EventPrePrintError = append(stdutil.EventPrePrintError, func(full string, msg string, err error) bool {
		color.Unset()
		COLOR_ERR.Set()
		return false
	})
	stdutil.EventPostPrintError = append(stdutil.EventPostPrintError, func(full string, msg string, err error) {
		color.Unset()
		COLOR.Set()
	})

	wrapperutil := os.Getenv(NAME)
	if wrapperutil == "true" {
		stdutil.PrintErr("You can't run WrapperUtil inside WrapperUtil...", nil)
		return
	} else if wrapperutil != "" {
		stdutil.PrintErr("Haha, very funny.", nil)
		return
	}

	var restart bool
	var shorter bool
	var timer bool
	var packets bool

	flag.BoolVar(&restart, "r", false, "Auto-restart unless cancelled")
	flag.BoolVar(&shorter, "s", false, "Shorter auto-restart timeout")
	flag.BoolVar(&timer, "t", false, "Measure execution time")
	flag.BoolVar(&packets, "p", false, "Enable program to send command packets to WrapperUtil.")

	flag.Parse()

	args := flag.Args()
	if len(args) <= 0 {
		stdutil.PrintErr("No command given. Run with --help for help", nil)
		return
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		<-c
		exit = true
	}()

	env := append(os.Environ(), NAME+"=true")

loop:
	for {
		var buffer *bytes.Buffer

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Env = env
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr

		if packets {
			buffer = bytes.NewBuffer(nil)
			stream := io.MultiWriter(buffer, os.Stdout)

			cmd.Stdout = stream
		} else {
			cmd.Stdout = os.Stdout
		}

		fmt.Println()
		color.Unset()

		var start time.Time
		if timer {
			start = time.Now()
		}

		err := cmd.Run()
		fmt.Println()
		if err != nil {
			stdutil.PrintErr("Couldn't start", err)
		}

		COLOR.Set()

		if timer {
			end := time.Now()

			fmt.Println("Program finished! Took " + end.Sub(start).String())
		}

		var packet Packet
		if packets {
			lines := strings.Split(buffer.String(), "\n")
			if len(lines) >= 1 {
				line := lines[len(lines)-1]

				if strings.HasPrefix(line, NAME) {
					line = strings.TrimPrefix(line, NAME)

					err = json.Unmarshal([]byte(line), &packet)
					if err != nil {
						stdutil.PrintErr("Program sent "+NAME+" packet, but it was invalid", nil)
					}
				}
			}
		}
		if (restart || packet.Restart) && !packet.Exit {
			timeout := RESTART_TIMEOUT
			if packet.Restart || shorter {
				timeout = RESTART_TIMEOUT_SHORT
			}
			for i := timeout; i >= 0; i-- {
				fmt.Printf("\rRestarting in %d... Cancel with Ctrl+C ", i)
				if i != 0 {
					time.Sleep(time.Second)
				}

				if exit {
					fmt.Println()
					break loop
				}
			}
			fmt.Println()
			continue
		}
		break
	}
	fmt.Println("Exiting!")
	color.Unset()
}
