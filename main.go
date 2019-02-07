package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"text/template"
	"github.com/BurntSushi/toml"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	config tomlConfig
	tmpl = template.Must(template.ParseGlob("*.tmpl"))
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
type tomlConfig struct {
	DB      database `toml:"database"`
}
type database struct {
	Server   string
	User     string
	Password string
}
type Employee struct {
    Id    int
    Name  string
    City string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
    dbUser := config.DB.User
    dbPass := config.DB.Password
    dbName := "myapp"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+config.DB.Server+":3306)/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

func Index(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM employee ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
	}
    emp := Employee{}
    res := []Employee{}
    for selDB.Next() {
        var id int
        var name, city string
        err = selDB.Scan(&id, &name, &city)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Name = name
        emp.City = city
        res = append(res, emp)
    }
    tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
}
func main() {
	bind := ""
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.StringVar(&bind, "bind", ":8080", "The socket to bind to.")
	flagset.Parse(os.Args[1:])

	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		panic(err)
	}

	r := prometheus.NewRegistry()
	r.MustRegister(httpRequestsTotal)
	r.MustRegister(version)
	r.MustRegister(scale_up_metric)
	r.MustRegister(scale_down_metric)
	scale_up_metric.Set(0)
	scale_down_metric.Set(0)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from example application.This is the version 2"))
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
	http.HandleFunc("/employees", Index)
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(bind, nil))
}
