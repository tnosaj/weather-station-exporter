package lib

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type weatherStationCollector struct {
	up      prometheus.Gauge
	baseurl string
	sites   []string
}

type promMetric struct {
	Name      string
	Desc      string
	Value     float64
	Label     []string
	LabelDesc []string
}

func NewMetricCollector(sites []string, url string, timeout int) *weatherStationCollector {
	return &weatherStationCollector{
		up: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "up",
				Help: "Dummy metric.",
			}),
		baseurl: url,
		sites:   sites,
	}
}

func (collector *weatherStationCollector) Describe(ch chan<- *prometheus.Desc) {
	log.Debugf("Describe %s", reflect.TypeOf(ch))
	//Need to register at least one
	ch <- collector.up.Desc()
}

func (collector *weatherStationCollector) collect(ch chan<- prometheus.Metric) error {
	ch <- prometheus.MustNewConstMetric(collector.up.Desc(), prometheus.GaugeValue, 1)
	var errors []string
	for _, site := range collector.sites {
		log.Debugf("Collecting - %s", site)
		cerr := getSiteMetrics(site, collector.baseurl, ch)
		if cerr != nil {
			log.Errorf("Error collecting site - %s", site)
			errors = append(errors, fmt.Sprintf("Error getting cluster based metrics: %s", cerr))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ","))
	}
	return nil
}

func (collector weatherStationCollector) Collect(ch chan<- prometheus.Metric) {
	err := collector.collect(ch)
	if err != nil {
		log.Errorf("Error collecting stats: %s", err)
	}
	return
}
