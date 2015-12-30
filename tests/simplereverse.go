package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var target *string

func main() {
	target = flag.String("target", "http://zuolin.com", "target URL for reverse proxy")
	flag.Parse()
	http.HandleFunc("/", report)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func report(w http.ResponseWriter, r *http.Request) {

	uri := *target + r.RequestURI

	fmt.Println(r.Method + ": " + uri)

	body2, err := ioutil.ReadAll(r.Body)
	fatal(err)
	fmt.Printf("Body: %v\n", string(body2))
	bf := bytes.NewBuffer(body2)

	rr, err := http.NewRequest(r.Method, uri, body2)
	fatal(err)
	copyHeader(r.Header, &rr.Header)

	// Create a client and query the target
	var transport http.Transport
	resp, err := transport.RoundTrip(rr)
	fatal(err)

	fmt.Printf("Resp-Headers: %v\n", resp.Header)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fatal(err)

	dH := w.Header()
	copyHeader(resp.Header, &dH)
	dH.Add("Requested-Host", rr.Host)

	w.Write(body)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func copyHeader(source http.Header, dest *http.Header) {
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}