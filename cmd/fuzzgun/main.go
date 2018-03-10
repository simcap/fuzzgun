package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"context"
	"os/signal"

	"github.com/simcap/fuzzgun"
)

var (
	input  string
	quoted bool
	rounds int
	tick   time.Duration
)

func init() {
	flag.StringVar(&input, "s", "", "Example string to be fuzzed")
	flag.BoolVar(&quoted, "quoted", false, "Display fizz result quoted (make all chars appear)")
	flag.IntVar(&rounds, "rounds", 0, "Number of fuzzing to perform")
	flag.DurationVar(&tick, "tick", 500*time.Millisecond, "Ticker to display fuzz result")
}

func main() {
	flag.Parse()
	if input == "" {
		flag.Usage()
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		signal.Stop(stop)
		cancel()
	}()

	go func() {
		select {
		case <-stop:
			cancel()
			fmt.Println("\n... emptying buffer (due to ticking)")
		case <-ctx.Done():
		}
	}()

	for s := range fuzzgun.Fuzz(ctx, input, rounds) {
		if quoted {
			fmt.Printf("%q\n", s)
		} else {
			fmt.Println(s)
		}
		time.Sleep(tick)
	}
}
