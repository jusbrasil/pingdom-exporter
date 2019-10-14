package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jusbrasil/pingdom_exporter/pkg/pingdom"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

var (
	// Injected during the build.
	VERSION string

	token       = os.Getenv("PINGDOM_API_TOKEN")
	waitSeconds int
	port        int

	pingdomUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pingdom_up",
		Help: "Whether the last pingdom scrape was successfull (1: up, 0: down)",
	})

	pingdomCheckStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pingdom_uptime_status",
		Help: "The current status of the check (1: up, 0: down)",
	}, []string{"id", "name", "hostname", "resolution", "paused", "tags"})

	pingdomCheckResponseTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pingdom_uptime_response_time",
		Help: "The response time of last test in milliseconds",
	}, []string{"id", "name", "hostname", "resolution", "paused", "tags"})
)

func init() {
	flag.IntVar(&waitSeconds, "wait", 10, "time (in seconds) between accessing the Pingdom  API")
	flag.IntVar(&port, "port", 9158, "port to listen on")

	prometheus.MustRegister(pingdomUp)
	prometheus.MustRegister(pingdomCheckStatus)
	prometheus.MustRegister(pingdomCheckResponseTime)
}

func retrieveChecksMetrics(client *pingdom.Client) {
	params := map[string]string{
		"include_tags": "true",
	}
	checks, err := client.Checks.List(params)
	if err != nil {
		log.Errorf("Error getting checks: %v", err)
		pingdomUp.Set(0)

		return
	}
	pingdomUp.Set(1)

	for _, check := range checks {
		id := strconv.Itoa(check.ID)

		var status float64
		switch check.Status {
		case "unknown":
			status = 0
		case "paused":
			status = 0
		case "up":
			status = 1
		case "unconfirmed_down":
			status = 0
		case "down":
			status = 0
		default:
			status = 100
		}

		resolution := strconv.Itoa(check.Resolution)

		paused := strconv.FormatBool(check.Paused)
		// Pingdom library doesn't report paused correctly,
		// so calculate it off the status.
		if check.Status == "paused" {
			paused = "true"
		}

		var tagsRaw []string
		for _, tag := range check.Tags {
			tagsRaw = append(tagsRaw, tag.Name)
		}
		tags := strings.Join(tagsRaw, ",")

		pingdomCheckStatus.WithLabelValues(
			id,
			check.Name,
			check.Hostname,
			resolution,
			paused,
			tags,
		).Set(status)

		pingdomCheckResponseTime.WithLabelValues(
			id,
			check.Name,
			check.Hostname,
			resolution,
			paused,
			tags,
		).Set(float64(check.LastResponseTime))
	}
}

func serverRun() {
	var client *pingdom.Client
	flag.Parse()

	if token == "" {
		log.Errorln("Pingdom API token must be provided via the PINGDOM_API_TOKEN environment variable, exiting")
		os.Exit(1)
	}

	client, err := pingdom.NewClientWithConfig(pingdom.ClientConfig{
		Token: token,
	})

	if err != nil {
		log.Errorln("Cannot create Pingdom client, exiting")
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Run metric retrieval loops
	go func() {
		for {
			select {
			case <-time.After(time.Second * time.Duration(waitSeconds)):
				retrieveChecksMetrics(client)

			case <-done:
				log.Infof("Received shutdown signal, exiting")
				os.Exit(0)
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	log.Infof("Pingdom Exporter v%s listening on %d\n", VERSION, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func main() {
	serverRun()
}
