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

	"github.com/jpillora/overseer"
	"github.com/jpillora/overseer/fetcher"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Makefile build
	version = ""

	apiKey  string
	keyFile *string
	timeout *time.Duration
)

const (
	ghUser    = "blockassets"
	ghRepo    = "litecoinpool_exporter"
	twiceADay = 12 * time.Hour
)

func main() {
	apiKey = os.Getenv("LITECOINPOOL_API_KEY")
	keyFile = flag.String("key-file", "./litecoinpool-api-key.txt", "Path to a file that has the API key in it.")
	port := flag.String("port", "5551", "The address to listen on for /metrics HTTP requests.")
	timeout = flag.Duration("timeout", 10*time.Second, "The amount of time to wait for litecoinpool to return.")
	noUpdate := flag.Bool("no-update", false, "Never do any updates. Example: -no-update=true")
	flag.Parse()

	portStr := fmt.Sprintf(":%s", *port)

	if *noUpdate {
		prog(overseer.State{Address: portStr})
	} else {
		overseerRun(portStr, twiceADay)
	}
}

func overseerRun(port string, interval time.Duration) {
	overseer.Run(overseer.Config{
		Program: prog,
		Address: port,
		Debug: true,
		Fetcher: &fetcher.Github{
			User:     ghUser,
			Repo:     ghRepo,
			Interval: interval,
		},
	})
}

func prog(state overseer.State) {
	log.Printf("%s %s on port %s\n", os.Args[0], version, state.Address)

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
	log.Fatal(http.Serve(state.Listener, nil))
}
