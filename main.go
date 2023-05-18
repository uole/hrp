package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

var (
	portFlag   = flag.Int("port", 0, "Listen port")
	targetFlag = flag.String("target", "", "Target rewrite host")
)

func main() {
	var (
		port   int
		err    error
		uri    *url.URL
		target string
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
	}
	if err = http.ListenAndServe(":"+strconv.Itoa(port), reverseProxy); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
