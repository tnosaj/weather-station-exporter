package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/tnosaj/weather-station-exporter/lib"
)

func main() {

	args, sites := evaluateInputs()

	setupLogger(args.debug)

	prometheus.MustRegister(lib.NewMetricCollector(
		sites,
		"https://www.bogner-lehner.eu/lwd/",
		args.timeout,
	))

	createHttpRoutes()

	log.Info("Beginning to serve on port :", args.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", args.port), nil))
}

type Args struct {
	debug   bool
	port    int
	timeout int
}

func evaluateInputs() (Args, []string) {
	var args Args
	var sites []string

	flag.BoolVar(&args.debug, "v", false, "Enable verbose debugging output")
	flag.IntVar(&args.port, "p", 8080, "Starts server on this port")
	flag.IntVar(&args.timeout, "t", 10, "Timeout for cloudant api calls")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [flags] command [command args…]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// hardcode for now
	sites = append(sites, "arlingsattel")
	sites = append(sites, "hengstpass")
	sites = append(sites, "hieflerstutzen")
	sites = append(sites, "menaueralm")

	return args, sites
}

func createHttpRoutes() {
	log.Debug("Creating HTTP routes")
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Sprint(w, "OK")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics", http.StatusMovedPermanently)
	})

	http.Handle("/metrics", promhttp.Handler())
}

func setupLogger(debug bool) {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	log.Debug("Configured logger")
}

func GetStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("OK")
}
