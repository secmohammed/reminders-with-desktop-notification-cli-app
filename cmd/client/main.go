package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/secmohammed/reminders-cli-app/client"
)

var (
	bcakendURI = flag.String("backend", "http://localhost:8000", "Backend API URL to use")
	helpFlag   = flag.Bool("help", false, "Print usage")
)

func main() {
	flag.Parse()
	s := client.NewSwitch(*bcakendURI)
	if *helpFlag || len(os.Args) == 1 {
		s.Help()
		return
	}
	err := s.Switch()
	if err != nil {
		fmt.Printf("cmd switch error: %v\n", err)
		os.Exit(2)
	}
}
