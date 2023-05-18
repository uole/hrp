package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

var (
	debugFlag  = flag.Bool("debug", false, "Enable debug")
	portFlag   = flag.Int("port", 0, "Listen port")
	targetFlag = flag.String("target", "", "Target rewrite host")
)

func envVariable(name string, defaultVal string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	return defaultVal
}

func main() {
	var (
		n         int
		port      string
		err       error
		targetUri *url.URL
		healthUri *url.URL
		targetUrl string
		healthUrl string
		address   string
	)
	flag.Parse()
	port = envVariable("HTTP_PORT", strconv.Itoa(*portFlag))
	targetUrl = envVariable("TARGET_URL", *targetFlag)
	healthUrl = envVariable("HEALTH_URL", "")
	if targetUri, err = url.Parse(targetUrl); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if healthUrl != "" {
		healthUri, err = url.Parse(healthUrl)
	}
	if n, err = strconv.Atoi(port); err != nil || n <= 0 {
		port = "80"
	}
	log.Printf("reverse proxy url: %s\n", targetUrl)
	reverseProxy := &httputil.ReverseProxy{}
	reverseProxy.Rewrite = func(request *httputil.ProxyRequest) {
		if healthUri != nil && request.In.URL.Path == "/_/health" {
			request.SetURL(healthUri)
			request.Out.URL.Path = "/"
		} else {
			request.SetURL(targetUri)
			if *debugFlag {
				log.Printf("rewrite url %s -> %s", request.In.URL.String(), request.Out.URL.String())
			}
		}
	}
	address = fmt.Sprintf(":%s", port)
	log.Printf("http listen on %s\n", address)
	if err = http.ListenAndServe(address, reverseProxy); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
