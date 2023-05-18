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

func main() {
	var (
		port    int
		err     error
		uri     *url.URL
		target  string
		address string
	)
	flag.Parse()
	port = *portFlag
	if port <= 0 {
		port, _ = strconv.Atoi(os.Getenv("HTTP_PORT"))
	}
	if port <= 0 {
		port = 80
	}
	target = *targetFlag
	if target == "" {
		target = os.Getenv("TARGET_URL")
	}
	if uri, err = url.Parse(target); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	reverseProxy := &httputil.ReverseProxy{}
	reverseProxy.Rewrite = func(request *httputil.ProxyRequest) {
		request.SetURL(uri)
		if *debugFlag {
			log.Printf("rewrite url %s -> %s", request.In.URL.String(), request.Out.URL.String())
		}
	}
	address = fmt.Sprintf(":%d", port)
	log.Printf("http listen on %s\n", address)
	if err = http.ListenAndServe(address, reverseProxy); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
