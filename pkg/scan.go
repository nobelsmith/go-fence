package pkg

import (
	"log"
	"strings"

	"golang.org/x/exp/slices"
)

func Scan(logFile string, protectedIPs []string, forbiddenLocations []string, dryRun bool) error {
	logLines, err := readLogFile(logFile)
	if err != nil {
		return err
	}
	scriptKitties := []string{}
	for _, v := range logLines {
		// if protected IP continue to next line
		if MatchPatternSlice(protectedIPs, v.RemoteAddr) {
			continue
			// if the ip address is already flagged continue to next line
		} else if slices.Contains(scriptKitties, v.RemoteAddr) {
			continue
			// if the request matches a phrase from the honeyPot append user to slice to be banned
		} else if stringContainsSlice(forbiddenLocations, strings.ToLower(v.RequestURI)) {
			scriptKitties = append(scriptKitties, v.RemoteAddr)
		}
	}

	for _, v := range scriptKitties {
		// ban the offending user
		if dryRun {
			log.Println("(dry-run) banned", v)
		} else {
			err := banScriptKitty(v)
			if err != nil {
				return err
			}
			log.Println("banned", v)
		}
	}

	return nil
}
