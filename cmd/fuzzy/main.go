package main

import (
	"flag"
	"github.com/simcap/fuzzy"
	"fmt"
	"os"
	"time"
)

var (
	input string
)

func init() {
	flag.StringVar(&input, "s", "", "Example string to be fuzzed")
}

func main() {
	flag.Parse()
	if input == "" {
		flag.Usage()
		os.Exit(1)
	}

	out := fuzzy.Fuzz(input)
	for s := range out {
		fmt.Println(s)
		time.Sleep(1 * time.Second)
	}
}
