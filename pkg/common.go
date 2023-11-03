package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

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

// banScriptKitty bans a user using iptables and a bash exec
func banScriptKitty(ip string) error {
	err := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-j", "DROP").Run()
	return err
}

func readLogFile(logFile string) ([]LogLine, error) {
	// read access.log
	file, err := os.Open(logFile)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		logLines = append(logLines, v)
	}
	return logLines, nil
}

func ReadConfig(cfgFile string) error {
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			if cfgFile == "" {
				fmt.Println(`no config file found, to create a default config file at $HOME/.go-fence.yaml by running "go-fence init"`)
			} else {
				fmt.Println("unable to find config file at ", cfgFile)
			}
		} else {
			fmt.Println("problem with config file: ", err)
		}
	}
	return err
}
