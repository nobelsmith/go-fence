package main

import (
	"encoding/json"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nxadm/tail"
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
	// Create a tail
	t, err := tail.TailFile(
		logFile, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		panic(err)
	}

	// check to see if line contains a request to one of the honeyPot phrases and ban the user
	for line := range t.Lines {
		var v LogLine
		err := json.Unmarshal([]byte(line.Text), &v)
		if err != nil {
			log.Println(err)
		}

		// if protected IP continue to next line
		if MatchPatternSlice(protectedIPs, v.RemoteAddr) {
			continue
			// if the request matches a phrase from the honeyPot append user to slice to be banned
		} else if stringContainsSlice(honeyPot, strings.ToLower(v.RequestURI)) {
			err := banPooh(v.RemoteAddr)
			if err != nil {
				log.Println(err)
			}
			log.Println("banned", v.RemoteAddr)
		}
	}
}
