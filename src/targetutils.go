package main

import (
	"fmt"
	"regexp"
	"strings"
)

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

func extractTarget(line string, r *regexp.Regexp) (string, error) {
	matches := r.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("no target found in line: %s", line)
}
