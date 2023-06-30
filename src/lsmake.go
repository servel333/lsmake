package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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
	if !fileExists(makefile) {
		fmt.Printf("Makefile %s does not exist.\n", makefile)
		os.Exit(1)
	}

	targets, err := listTargets(makefile)
	if err != nil {
		fmt.Printf("Error reading Makefile %s: %s\n", makefile, err)
		os.Exit(1)
	}

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
	seen := make(map[string]bool)

	targetPattern := `(^[^()<>~:;,!?"'*|/]+):`
	r := regexp.MustCompile(targetPattern)

	for scanner.Scan() {
		line := scanner.Text()
		if target, err := extractTarget(line, r); err == nil && target != "" {
			if target[0] != '.' && !seen[target] {
				targets = append(targets, target)
				seen[target] = true
			}
		} else if includeFile := extractIncludeFile(line); includeFile != "" {
			includeFile = resolveIncludedFilePath(filename, includeFile)
			includedTargets, err := listTargets(includeFile)
			if err != nil {
				return nil, fmt.Errorf("error processing included file %s: %s", includeFile, err)
			}
			addUniqueTargets(targets, seen, includedTargets)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.Strings(targets)

	return targets, nil
}

func extractTarget(line string, r *regexp.Regexp) (string, error) {
	matches := r.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("no target found in line: %s", line)
}

func extractIncludeFile(line string) string {
	includePrefix := "include"
	if strings.HasPrefix(line, includePrefix) {
		includeFile := strings.TrimSpace(strings.TrimPrefix(line, includePrefix))
		includeFile = strings.Trim(includeFile, "\"'")
		return includeFile
	}
	return ""
}

func resolveIncludedFilePath(baseFile, includeFile string) string {
	baseDir := filepath.Dir(baseFile)
	return filepath.Join(baseDir, includeFile)
}

func addUniqueTargets(targets []string, seen map[string]bool, newTargets []string) {
	for _, target := range newTargets {
		if !seen[target] {
			targets = append(targets, target)
			seen[target] = true
		}
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func displayUsage() {
	fmt.Println("Usage: lsmake [options]")
	fmt.Println("Options:")
	flag.PrintDefaults()
}
