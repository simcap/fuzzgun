package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/simcap/fuzzy"
)

var (
	input  string
	quoted bool
)

func init() {
	flag.StringVar(&input, "s", "", "Example string to be fuzzed")
	flag.BoolVar(&quoted, "quoted", false, "Diplay fizz result quoted (make all chars appear)")
}

func main() {
	flag.Parse()
	if input == "" {
		flag.Usage()
		os.Exit(1)
	}

	out := fuzzy.Fuzz(input)
	for s := range out {
		if quoted {
			fmt.Printf("%q\n", s)
		} else {
			fmt.Println(s)
		}
		time.Sleep(1 * time.Second)
	}
}
