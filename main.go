package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var logger = log.New(os.Stdout, "IPchecker ", log.LstdFlags|log.Lshortfile|log.Ltime|log.LUTC)
var port = getEnv("SERVER_PORT", "8080")
var pport = getEnv("PPROF_PORT", "5050")

var html = `<!doctype html>
<html lang="en">
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
        <title> WhatismyIP </title>
    </head>
	<body>
		<div class="navbar-wrapper">
			<div class="container">
				<h2>What is my ip address?</h2>
				<hr style="height:1px;border-width:0;color:gray;background-color:gray">
			</div>
		</div>
		<div class="container">
			<div class="row">
				<div class="col-xs-12 col-sm-6">
					<h3>Your Connection:</h3>
					<table style="font-family: Arial; font-size: 4;">
						<tr>
							<td><b>Host:</b></td>
							<td>%s</td>
						</tr>
						<tr>
							<td><b>Remote Address:</b></td>
							<td>%s</td>
						</tr>
						<tr>
							<td><b>User Agent:</b></td>
							<td>%s</td>
						</tr>
						<tr>
							<td><b>MIME type:</b></td>
							<td>%s</td>
						</tr>
						<tr>
							<td><b>X-Forwarded-For:</b></td>
							<td>%s</td>
						</tr>
					</table>
				</div>
			</div>
		</div>
		<hr style="height:1px;border-width:0;color:gray;background-color:gray">
		<div>
		    <p>
		    Copyright &copy; 2049
			michalswi
		</div>
		<div>	
			<a href="https://github.com/michalswi/whatismyip" target="-blank">
				<img alt="github" src="https://github.com/fluidicon.png" width=40" height="40">
			</a>
        </div>
	</body>
</html>
`

func main() {
	// pprof
	go func() {
		router := mux.NewRouter()
		router.HandleFunc("/debug/pprof/", pprof.Index)
		router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		router.Handle("/debug/pprof/block", pprof.Handler("block"))
		logger.Println("Pprof is ready at port", pport)
		logger.Fatal(http.ListenAndServe(":"+pport, router))
	}()
	// web server
	http.HandleFunc("/", locate)
	http.HandleFunc("/ip", getIP)
	http.HandleFunc("/hz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./robots.txt")
		logger.Println("User-Agent for robots.txt:", r.Header.Get("User-Agent"))
	})
	logger.Println("Server is ready to handle requests at port", port)
	logger.Fatal(http.ListenAndServe(":"+port, nil))
}

func locate(w http.ResponseWriter, r *http.Request) {
	var request []string

	// Request line
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)

	// Header
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	request = append(request, fmt.Sprintf("Remote Address: %s", r.RemoteAddr))
	xforwardedfor := r.Header.Get("X-FORWARDED-FOR")
	if xforwardedfor == "" {
		xforwardedfor = "N/A"
	}
	request = append(request, fmt.Sprintf("X-Forwarded-For: %s", xforwardedfor))

	for name, headers := range r.Header {
		// name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}
	fmt.Fprintf(w, html, r.Host, r.RemoteAddr, r.Header.Get("User-Agent"), r.Header.Get("Accept"), xforwardedfor)

	// display all request's headers
	// logger.Println(strings.Join(request, "\n"))
	// fmt.Println()
	// or
	logger.Println()
	fmt.Println(strings.Join(request, "\n"))
	fmt.Println()
	// or
	// logger.Println()
	// var mySlice interface{} = request
	// fmt.Println(mySlice)

	// port scanner
	// var remoteAddress []string
	// if r.RemoteAddr != "" {
	// 	remoteAddress = strings.Split(r.RemoteAddr, ":")
	// }
	// p := getPorts(remoteAddress[0])
	// logger.Printf("For RemoteIP: %s open ports: %d", remoteAddress[0], p)
}

func getIP(w http.ResponseWriter, r *http.Request) {
	var remoteAddress []string
	if r.RemoteAddr != "" {
		remoteAddress = strings.Split(r.RemoteAddr, ":")
	}
	logger.Printf("IP request from: %s", remoteAddress[0])
	w.Write([]byte(remoteAddress[0]))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// port scanner
func getPorts(ra string) []int {
	validPorts := []int{}
	invalidPorts := []int{}
	for _, port := range []int{21, 22, 3389} {
		addrs := fmt.Sprintf("%s:%d", ra, port)
		conn, err := net.DialTimeout("tcp", addrs, 5*time.Second)
		if err != nil {
			invalidPorts = append(invalidPorts, port)
		} else {
			validPorts = append(validPorts, port)
			conn.Close()
		}
	}
	return validPorts
}
