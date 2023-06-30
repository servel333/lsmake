package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
)

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
