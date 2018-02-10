package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Makefile build
	version = ""
)

func main() {
	apiKey := os.Getenv("LITECOINPOOL_API_KEY")
	keyFile := flag.String("key-file", "./litecoinpool-api-key.txt", "Path to a file that has the API key in it.")
	port := flag.String("port", "5551", "The address to listen on for /metrics HTTP requests.")
	timeout := flag.Duration("timeout", 10*time.Second, "The amount of time to wait for litecoinpool to return.")
	flag.Parse()

	if len(apiKey) == 0 {
		dat, err := ioutil.ReadFile(*keyFile)
		if err != nil || len(dat) == 0 {
			log.Println(err)
			log.Fatal("Please specify the API key in LITECOINPOOL_API_KEY or use the -key-file argument.")
		}
		apiKey = strings.TrimSpace(string(dat))
	}

	prometheus.MustRegister(NewExporter(apiKey, *timeout))

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("%s %s", os.Args[0], version)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
