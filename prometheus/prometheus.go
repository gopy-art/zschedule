package prometheus

import (
	"net"
	"net/http"
	"strings"
	"time"
	"zschedule/cmd"
	logger "zschedule/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ScheduleInfo, TotalExecutedCount *prometheus.GaugeVec
	WorkerIp                         string
)

var Uptime prometheus.Gauge

func PrometheusInit(prometheusHost string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logger.ErrorLogger.Printf("error with get worker ip%v\n", err)
	}
	if len(addrs) > 1 {
		WorkerIp = strings.Split(addrs[1].String(), "/")[0]
	}

	ScheduleInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "zschedule",
			Name:      "info",
			Help:      "all info about zschedules commands",
		},
		[]string{
			"WorkerIp",
			"Name",
			"Command",
			"Interval",
			"Executed",
			"Version",
		}, // labels
	)
	TotalExecutedCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "zschedule",
			Name:      "total_executed",
			Help:      "total executed task throw zschedule",
		},
		[]string{
			"WorkerIp",
			"Version",
		}, // labels
	)
	Uptime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "zschedule_uptime",
		Help: "worker uptime in seconds",
	})

	prometheus.MustRegister(TotalExecutedCount, ScheduleInfo, Uptime)
	http.Handle("/metrics", promhttp.Handler())
	go recordMetrics()
	if err := http.ListenAndServe(prometheusHost, nil); err != nil {
		logger.ErrorLogger.Fatalf("could not run prometheus server: %s", err.Error())
	}
}

func IncreaseTotalCount() {
	TotalExecutedCount.With(prometheus.Labels{
		"WorkerIp": WorkerIp,
		"Version":  cmd.AppVersion,
	}).Add(1)
}

func IncreaseScheduleInfoCount(command, name, interval, executed string) {
	ScheduleInfo.With(prometheus.Labels{
		"WorkerIp": WorkerIp,
		"Name":     name,
		"Command":  command,
		"Interval": interval,
		"Executed": executed,
		"Version":  cmd.AppVersion,
	}).Add(1)
}

func IsValidIpv4(ip string) bool {
	ipaddr := net.ParseIP(ip)
	if ipaddr == nil {
		return false
	}
	return true
}

func recordMetrics() {
	startTime := time.Now()
	for {
		Uptime.Set(time.Since(startTime).Seconds())
		time.Sleep(1 * time.Second)
	}
}
