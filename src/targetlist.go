package main

import (
	"fmt"
	"os"
)

func listTargets(makefiles []string) (map[string][]string, error) {
	targets := make(map[string][]string)

	noMakefileExists := true
	for _, makefile := range makefiles {
		if !fileExists(makefile) {
			fmt.Fprintf(os.Stderr, "Makefile %s does not exist.\n", makefile)
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
