package main

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/lookfirst/litecoinpool_exporter/litecoinpool"
	"github.com/prometheus/client_golang/prometheus"
)

//
var (
	namespace    = "lcp"
	idLabelNames = []string{"id"}
)

//
func newGauge(metricName string, help string, constLabels prometheus.Labels) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        metricName,
		Help:        help,
		ConstLabels: constLabels,
	})
}

//
func newGaugeVec(metricName string, help string, constLabels prometheus.Labels, labels []string) prometheus.GaugeVec {
	return *prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   namespace,
			Name:        metricName,
			Help:        help,
			ConstLabels: constLabels,
		},
		labels,
	)
}

type GaugeMap map[string]prometheus.Gauge
type GaugeVecMap map[string]prometheus.GaugeVec

// Collector interface
type Exporter struct {
	client      *litecoinpool.LCPClient
	constLabels prometheus.Labels
	User        GaugeMap
	Workers     map[string]GaugeVecMap
	Pool        GaugeMap
	Network     GaugeMap
	Market      GaugeMap
	sync.Mutex
}

//
func fmtTagName(tag string) string {
	return strings.Split(tag, ",")[0]
}

//
func fmtGaugeName(interfaceName string, tag string) string {
	return fmt.Sprintf("%s_%s", strings.ToLower(interfaceName), fmtTagName(tag))
}

//
func fmtGaugeDescription(interfaceName string, tag string) string {
	return fmt.Sprintf("%s%s", interfaceName, fmtTagName(tag))
}

//
func newGaugeMap(api interface{}, constLabels prometheus.Labels) GaugeMap {
	interfaceName := lookupStructName(api)
	tags := lookupStructTags(api)

	metrics := GaugeMap{}
	for _, tag := range tags {
		gaugeName := fmtGaugeName(interfaceName, tag)
		gauge := newGauge(
			gaugeName,
			fmtGaugeDescription(interfaceName, tag),
			constLabels)
		metrics[gaugeName] = gauge
	}

	return metrics
}

//
func newGaugeVecMap(api interface{}, labels []string, constLabels prometheus.Labels) GaugeVecMap {
	interfaceName := lookupStructName(api)
	tags := lookupStructTags(api)

	metrics := GaugeVecMap{}
	for _, tag := range tags {
		gaugeName := fmtGaugeName(interfaceName, tag)
		metrics[gaugeName] = newGaugeVec(fmtGaugeName(interfaceName, tag),
			fmtGaugeDescription(interfaceName, tag),
			constLabels, labels)
	}

	return metrics
}

//
func NewExporter(apiKey string, timeout time.Duration) *Exporter {
	constLabels := prometheus.Labels{"key": apiKey[:8]}
	return &Exporter{
		client:      litecoinpool.NewClient(apiKey, timeout),
		constLabels: constLabels,
		Workers:     make(map[string]GaugeVecMap),
		User:        newGaugeMap(litecoinpool.User{}, constLabels),
		Pool:        newGaugeMap(litecoinpool.Pool{}, constLabels),
		Network:     newGaugeMap(litecoinpool.Network{}, constLabels),
		Market:      newGaugeMap(litecoinpool.Market{}, constLabels),
	}
}

//
func setFieldsOnGauge(gaugeMap GaugeMap, strct interface{}) {
	interfaceName := lookupStructName(strct)

	structValue := reflect.ValueOf(strct)
	for i := 0; i < structValue.NumField(); i++ {
		value := structValue.Field(i).Interface()
		flt, err := ConvertToFloat(value)
		if err == nil {
			key := structValue.Type().Field(i).Tag.Get("json")
			gauge, ok := gaugeMap[fmtGaugeName(interfaceName, key)]
			if ok {
				gauge.Set(flt)
			}
		}
	}
}

//
func setFieldsOnGaugeVec(gaugeVecMap GaugeVecMap, strct interface{}, labelValue string) {
	interfaceName := lookupStructName(strct)
	structValue := reflect.ValueOf(strct)

	for i := 0; i < structValue.NumField(); i++ {
		value := structValue.Field(i).Interface()

		flt, err := ConvertToFloat(value)
		if err == nil {
			key := structValue.Type().Field(i).Tag.Get("json")
			gaugeName := fmtGaugeName(interfaceName, key)
			gauge, ok := gaugeVecMap[gaugeName]
			if ok {
				gauge.WithLabelValues(labelValue).Set(flt)
			}
		}
	}
}

//
func setFieldsOnGaugeVecs(e *Exporter, poolData *litecoinpool.PoolData) {
	workers := poolData.Workers
	for name, worker := range workers {
		// Workers are generated on each request
		if _, ok := e.Workers[name]; !ok {
			e.Workers[name] = newGaugeVecMap(litecoinpool.Worker{}, idLabelNames, e.constLabels)
		}

		setFieldsOnGaugeVec(e.Workers[name], worker, name)
	}
}

//
func setFieldsOnGauges(e *Exporter, poolData *litecoinpool.PoolData) {
	setFieldsOnGauge(e.User, poolData.User)
	setFieldsOnGauge(e.Pool, poolData.Pool)
	setFieldsOnGauge(e.Network, poolData.Network)
	setFieldsOnGauge(e.Market, poolData.Market)
}

//
func collectGauge(gaugeMap GaugeMap, ch chan<- prometheus.Metric) {
	for _, metric := range gaugeMap {
		metric.Collect(ch)
	}
}

//
func collectGaugeVec(gaugeVecMap map[string]GaugeVecMap, ch chan<- prometheus.Metric) {
	for _, gaugeVec := range gaugeVecMap {
		for _, metric := range gaugeVec {
			metric.Collect(ch)
		}
	}
}

// Outputs the gauge values on the channel
func collectGauges(e *Exporter, ch chan<- prometheus.Metric) {
	collectGaugeVec(e.Workers, ch)
	collectGauge(e.User, ch)
	collectGauge(e.Pool, ch)
	collectGauge(e.Network, ch)
	collectGauge(e.Market, ch)
}

//
func describeGauge(gaugeMap GaugeMap, ch chan<- *prometheus.Desc) {
	for _, metric := range gaugeMap {
		metric.Describe(ch)
	}
}

//
func describeGaugeVec(gaugeVecMap map[string]GaugeVecMap, ch chan<- *prometheus.Desc) {
	for _, gaugeVec := range gaugeVecMap {
		for _, metric := range gaugeVec {
			metric.Describe(ch)
		}
	}
}

//
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	describeGauge(e.User, ch)
	describeGaugeVec(e.Workers, ch)
	describeGauge(e.Pool, ch)
	describeGauge(e.Network, ch)
	describeGauge(e.Market, ch)
}

//
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Prevents multiple concurrent calls
	e.Lock()
	defer e.Unlock()

	poolData, err := e.client.Fetch()
	if err != nil {
		log.Println(err)
		return
	}

	setFieldsOnGauges(e, poolData)
	setFieldsOnGaugeVecs(e, poolData)

	collectGauges(e, ch)
}
