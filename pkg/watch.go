package pkg

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/nxadm/tail"
)

func Watch(logFile string, protectedIPs []string, forbiddenLocations []string, dryRun bool) error {
	// Create a tail
	t, err := tail.TailFile(
		logFile, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		return err
	}

	// check to see if line contains a request to one of the forbiddenLocations phrases and ban the user
	for line := range t.Lines {
		var v LogLine
		err := json.Unmarshal([]byte(line.Text), &v)
		if err != nil {
			return err
		}

		// if protected IP continue to next line
		if MatchPatternSlice(protectedIPs, v.RemoteAddr) {
			continue
			// if the request matches a phrase from the forbiddenLocations append user to slice to be banned
		} else if stringContainsSlice(forbiddenLocations, strings.ToLower(v.RequestURI)) {
			if dryRun {
				log.Println("(dry-run) banned", "[", v.RemoteAddr, "] - ", line.Text)
			} else {
				err := banScriptKitty(v.RemoteAddr)
				if err != nil {
					return err
				}
				log.Println("banned", "[", v.RemoteAddr, "] - ", line.Text)
			}
		}
	}
	if t.Err() != nil {
		return t.Err()
	}
	return nil
}
