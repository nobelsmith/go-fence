package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

// protectedIPs will not be banned even if they access routes in the honeyPot
// protectedIPs use filepath syntax for wildcard matches
var protectedIPs = []string{"192.168.*", "172.16.*", "10.*", "127*"}

// honeyPot is a slice containing strings to search for inside of the requestUri
var honeyPot = []string{"wp-admin", "wp-includes", ".aspx"}

// logFile is the path to the access.log to be parsed
var logFile = "/var/log/nginx/access.log"

type LogLine struct {
	BodyBytesSent        string `json:"body_bytes_sent"`
	BytesSent            string `json:"bytes_sent"`
	HTTPHost             string `json:"http_host"`
	HTTPReferer          string `json:"http_referer"`
	HTTPUserAgent        string `json:"http_user_agent"`
	Msec                 string `json:"msec"`
	RemoteAddr           string `json:"remote_addr"`
	RequestMethod        string `json:"request_method"`
	RequestURI           string `json:"request_uri"`
	ServerPort           string `json:"server_port"`
	ServerProtocol       string `json:"server_protocol"`
	SslProtocol          string `json:"ssl_protocol"`
	Status               string `json:"status"`
	UpstreamResponseTime string `json:"upstream_response_time"`
	UpstreamAddr         string `json:"upstream_addr"`
	UpstreamConnectTime  string `json:"upstream_connect_time"`
}

// stringContainsSlice returns true if s contains any of the strings in slice
func stringContainsSlice(slice []string, s string) bool {
	for _, v := range slice {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}

// MatchPatternSlice returns true if s matches any of the strings in pattern using filepath wildcard syntax
func MatchPatternSlice(pattern []string, s string) bool {
	for _, p := range pattern {
		if p == "*" {
			return true
		}
		matched, err := filepath.Match(p, s)
		if err != nil {
			log.Println(err)
		}
		if matched {
			return true
		}
	}
	return false
}

// banPooh bans a user using iptables and a bash exec
func banPooh(ip string) error {
	err := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-j", "DROP").Run()
	return err
}

func main() {
	// read access.log
	file, err := os.Open(logFile)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	// parse access.log
	d := json.NewDecoder(file)
	logLines := []LogLine{}
	for {
		var v LogLine
		if err := d.Decode(&v); err == io.EOF {
			break // done decoding file
		} else if err != nil {
			log.Println(err)
		}
		logLines = append(logLines, v)
	}

	// failedRequests := map[string][]LogLine{}
	// poohBears is a slice of ip addresses to be banned
	poohBears := []string{}
	for _, v := range logLines {
		// if protected IP continue to next line
		if MatchPatternSlice(protectedIPs, v.RemoteAddr) {
			continue
			// if the ip address is already flagged continue to next line
		} else if slices.Contains(poohBears, v.RemoteAddr) {
			continue
			// if the request matches a phrase from the honeyPot append user to slice to be banned
		} else if stringContainsSlice(honeyPot, strings.ToLower(v.RequestURI)) {
			poohBears = append(poohBears, v.RemoteAddr)
		} //else if v.Status[:1] != "2" && v.Status[:1] != "3" {
		// 	failedRequests[v.RemoteAddr] = append(failedRequests[v.RemoteAddr], v)
		// }
	}

	for _, v := range poohBears {
		// ban the offending user
		log.Println(v, "banned")
		err := banPooh(v)
		if err != nil {
			log.Println(err)
		}
	}
}

// keys := make([]string, 0, len(failedRequests))
// for k := range failedRequests {
// 	keys = append(keys, k)
// }
// sort.SliceStable(keys, func(i, j int) bool {
// 	return len(failedRequests[keys[i]]) > len(failedRequests[keys[j]])
// })
// for _, key := range keys {
// 	if len(failedRequests[key]) < 10 {
// 		break
// 	}
// 	for _, logline := range failedRequests[key] {
// 		log.Println(logline.RemoteAddr, logline.Status, logline.RequestURI)
// 	}

// }
