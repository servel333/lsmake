package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	helpPtr := flag.Bool("help", false, "Display usage information. Example: lsmake -help or lsmake Makefile1 Makefile2")
	flag.Parse()

	if *helpPtr {
		displayUsage()
		os.Exit(0)
	}

	makefiles := flag.Args()
	if len(makefiles) == 0 {
		makefiles = []string{"Makefile"}
	}

	targets, err := listTargets(makefiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading Makefiles: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Targets in Makefiles:")
	for makefile, makefileTargets := range targets {
		fmt.Printf("%s:\n", makefile)
		for _, target := range makefileTargets {
			fmt.Println(target)
		}
		fmt.Println()
	}
}

func displayUsage() {
	fmt.Println("Usage: lsmake [options] [Makefiles...]")
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println("Example: lsmake Makefile1 Makefile2")
}
