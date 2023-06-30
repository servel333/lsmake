package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
)

func main() {
	makefilePtr := flag.String("f", "Makefile", "Specify the Makefile to parse")
	helpPtr := flag.Bool("help", false, "Display usage information")
	flag.Parse()

	if *helpPtr {
		displayUsage()
		os.Exit(0)
	}

	makefile := *makefilePtr
	targets, err := listTargets(makefile)
	if err != nil {
		fmt.Printf("Error reading Makefile: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Targets in", makefile+":")
	for _, target := range targets {
		fmt.Println(target)
	}
}

func listTargets(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	targets := make([]string, 0)

	targetPattern := `(^[^()<>~:;,!?"'*|/]+):`
	r := regexp.MustCompile(targetPattern)

	seen := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		matches := r.FindStringSubmatch(line)
		if len(matches) > 1 {
			target := matches[1]
			if target[0] != '.' {
				if !seen[target] {
					targets = append(targets, target)
					seen[target] = true
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.Strings(targets)

	return targets, nil
}

func displayUsage() {
	fmt.Println("Usage: lsmake [options]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}
