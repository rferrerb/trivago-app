package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	version = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "version",
		Help: "Version information about this binary",
		ConstLabels: map[string]string{
			"version": "v0.1.0",
		},
	})
	scale_up_metric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "trivago_app_scale_up",
		Help: "metric to trigger scale events",
	})
	scale_down_metric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "trivago_app_scale_down",
		Help: "metric to trigger scale events",
	})
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"code", "method"})
)

func main() {
	bind := ""
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.StringVar(&bind, "bind", ":8080", "The socket to bind to.")
	flagset.Parse(os.Args[1:])

	r := prometheus.NewRegistry()
	r.MustRegister(httpRequestsTotal)
	r.MustRegister(version)
	r.MustRegister(scale_up_metric)
	r.MustRegister(scale_down_metric)
        scale_up_metric.Set(0)
        scale_down_metric.Set(0)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from example application."))
		w.Write([]byte("Menudo mierdazo")))
	})
	notfound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	scale_up := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Scale up"))
		scale_up_metric.Inc()
	})
	scale_down := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Scale down"))
		scale_down_metric.Inc()
	})
	http.Handle("/", promhttp.InstrumentHandlerCounter(httpRequestsTotal, handler))
	http.Handle("/scale_up", promhttp.InstrumentHandlerCounter(httpRequestsTotal, scale_up))
	http.Handle("/scale_down", promhttp.InstrumentHandlerCounter(httpRequestsTotal, scale_down))
	http.Handle("/err", promhttp.InstrumentHandlerCounter(httpRequestsTotal, notfound))

	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(bind, nil))
}
