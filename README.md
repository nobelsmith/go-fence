<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="https://i.imgur.com/6wj0hh6.jpg" alt="Project logo"></a>
</p>

<h3 align="center">Go-Fence</h3>

<div align="center">

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

</div>

---

<p align="center"> Simple go executable to monitor Nginx logs and ban bad actors.
    <br> 
</p>

## 📝 Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Deployment](#deployment)
- [Usage](#usage)
- [Built Using](#built_using)


## 🧐 About <a name = "about"></a>

This project was designed to parse Nginx logs and then take action based upon the information it gathers. Keywords are selected, and then if they are found within request URI in the nginx logs the offending IP is banned at the kernel level using iptables. This program assumes that iptables is installed and you are running a bash shell.

## 🏁 Getting Started <a name = "getting_started"></a>

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See [deployment](#deployment) for notes on how to deploy the project on a live system.

### Prerequisites

Install Nginx

https://nginx.org/en/docs/beginners_guide.html

```
sudo apt-get install nginx
sudo systemctl start nginx
```

Make sure iptables is installed

```
which iptables
```
if you dont see something like '/usr/sbin/iptables' then you may install it

```
sudo apt-get install iptables
```
### Installing

Create a directory to clone the repository into

```
mkdir projects
```

Clone the repository from github
```
git clone git@github.com:nobelsmith/go-fence.git
```

## 🎈 Usage <a name="usage"></a>

### Change Nginx Log Format for Parsing.

Add the following lines to your nginx.conf (likely found at /etc/nginx/nginx.conf)

```
	log_format logger-json-log escape=json '{'
	'"body_bytes_sent":"$body_bytes_sent",'
	'"bytes_sent":"$bytes_sent",'
	'"http_host":"$http_host",'
	'"http_referer":"$http_referer",'
	'"http_user_agent":"$http_user_agent",'
	'"msec":"$msec",'
	'"remote_addr":"$remote_addr",'
	'"request_method":"$request_method",'
	'"request_uri":"$request_uri",'
	'"server_port":"$server_port",'
	'"server_protocol":"$server_protocol",'
	'"ssl_protocol":"$ssl_protocol",'
	'"status":"$status",'
	'"upstream_response_time":"$upstream_response_time",'
	'"upstream_addr":"$upstream_addr",'
	'"upstream_connect_time":"$upstream_connect_time"'
	'}';

	access_log /var/log/nginx/access.log logger-json-log;
```
An example nginx.conf file can be found in the examples folder of the project.

### Setup Protected IPs, Forbiddedn Locations, and Log File within your config

```
forbiddenlocations:
    - wp-admin
    - wp-login.php
nginxlogfile: /var/log/nginx/access.log
protectedips:
    - 192.168.*
    - 172.16.*
    - 10.*
    - 127*
```

## 🚀 Deployment <a name = "deployment"></a>

### Build the project
```
go build -o go-fence
```

### Setup Project on Host Machine

Firstly setup a config file or let go-fence do it for you with the init command.
```
go-fence init
```

One way to use this executable would be to put it make it a systemd service. 
Below is an example systemd service config.
```
[Unit]
Description = go-fence

[Service]
Type           = simple
User           = root
Group          = root
LimitNOFILE    = 4096
Restart        = always
RestartSec     = 5s
StandardOutput = append:/var/log/go-fence.log
StandardError  = append:/var/log/go-fence.log
ExecStart      = /usr/local/bin/go-fence --config /home/nobel/.go-fence.yaml watch

[Install]
WantedBy = multi-user.target
```
## ⛏️ Built Using <a name = "built_using"></a>

- [Golang](https://go.dev/) - Language
- [Nginx](https://nginx.org/en/) - Webserver
- [iptables](https://www.netfilter.org/projects/iptables/index.html) - Packet Filter
