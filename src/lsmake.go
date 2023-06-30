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
		fmt.Printf("Error reading Makefiles: %s\n", err)
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

func listTargets(makefiles []string) (map[string][]string, error) {
	targets := make(map[string][]string)

	noMakefileExists := true
	for _, makefile := range makefiles {
		if !fileExists(makefile) {
			fmt.Printf("Makefile %s does not exist.\n", makefile)
			continue
		}
		noMakefileExists = false

		makefileTargets, err := parseMakefile(makefile)
		if err != nil {
			return nil, fmt.Errorf("error reading Makefile %s: %s", makefile, err)
		}

		targets[makefile] = makefileTargets
	}

	if noMakefileExists {
		return nil, fmt.Errorf("none of the provided Makefiles exist")
	}

	return targets, nil
}

func parseMakefile(makefile string) ([]string, error) {
	file, err := os.Open(makefile)
	if err != nil {
		return nil, fmt.Errorf("failed to open Makefile %s: %w", makefile, err)
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
			if target[0] != '.' {
				addUniqueTarget(&targets, seen, target)
			}
		} else if includeFile := extractIncludeFile(line); includeFile != "" {
			includeFile = resolveIncludedFilePath(makefile, includeFile)
			includedTargets, err := parseMakefile(includeFile)
			if err != nil {
				return nil, fmt.Errorf("error processing included file %s: %s", includeFile, err)
			}
			addUniqueTargets(&targets, seen, includedTargets)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read Makefile %s: %w", makefile, err)
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
	lineLower := strings.ToLower(line)
	if strings.HasPrefix(lineLower, includePrefix) {
		includeFile := strings.TrimSpace(strings.TrimPrefix(lineLower, includePrefix))
		includeFile = strings.Trim(includeFile, "\"'")
		return includeFile
	}
	return ""
}

func resolveIncludedFilePath(baseFile, includeFile string) string {
	baseDir := filepath.Dir(baseFile)
	includedPath, err := filepath.Abs(filepath.Join(baseDir, includeFile))
	if err != nil {
		return filepath.Join(baseDir, includeFile)
	}
	return includedPath
}

func addUniqueTargets(targets *[]string, seen map[string]bool, newTargets []string) {
	for _, target := range newTargets {
		addUniqueTarget(targets, seen, target)
	}
}

func addUniqueTarget(targets *[]string, seen map[string]bool, target string) {
	if !seen[target] {
		*targets = append(*targets, target)
		seen[target] = true
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func displayUsage() {
	fmt.Println("Usage: lsmake [options] [Makefiles...]")
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println("Example: lsmake Makefile1 Makefile2")
}
