package pkg

import (
	"log"
	"sort"
)

func BadRequests(logFile string, protectedIPs []string, forbiddenLocations []string) error {
	logLines, err := readLogFile(logFile)
	if err != nil {
		return err
	}
	log.Println("Processing", len(logLines), "lines")
	failedRequests := map[string][]LogLine{}
	for _, v := range logLines {
		// if protected IP continue to next line
		if MatchPatternSlice(protectedIPs, v.RemoteAddr) {
			continue
			// if the ip address is already flagged continue to next line
		} else if v.Status[:1] != "2" && v.Status[:1] != "3" {
			failedRequests[v.RemoteAddr] = append(failedRequests[v.RemoteAddr], v)
		}
	}

	keys := make([]string, 0, len(failedRequests))
	for k := range failedRequests {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return len(failedRequests[keys[i]]) > len(failedRequests[keys[j]])
	})
	for _, key := range keys {
		if len(failedRequests[key]) < 10 {
			break
		}
		for _, logline := range failedRequests[key] {
			log.Println(logline.RemoteAddr, logline.Status, logline.RequestURI)
		}
	}
	return nil
}
