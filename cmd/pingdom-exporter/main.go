package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/jusbrasil/pingdom-exporter/pkg/pingdom"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	// VERSION will hold the version number injected during the build.
	VERSION string

	token             string
	waitSeconds       int
	port              int
	outageCheckPeriod int
	defaultUptimeSLO  float64

	pingdomUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pingdom_up",
		Help: "Whether the last pingdom scrape was successfull (1: up, 0: down)",
	})

	pingdomOutageCheckPeriod = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pingdom_slo_period_seconds",
		Help: "Outage check period, in seconds",
	})

	pingdomCheckStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pingdom_uptime_status",
		Help: "The current status of the check (1: up, 0: down)",
	}, []string{"id", "name", "hostname", "resolution", "paused", "tags"})

	pingdomCheckResponseTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pingdom_uptime_response_time_seconds",
		Help: "The response time of last test, in seconds",
	}, []string{"id", "name", "hostname", "resolution", "paused", "tags"})

	pingdomOutages = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pingdom_outages_total",
		Help: "Number of outages within the outage check period",
	}, []string{"id", "name", "hostname", "tags"})

	pingdomCheckErrorBudget = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pingdom_uptime_slo_error_budget_total_seconds",
		Help: "Maximum number of allowed downtime, in seconds, according to the uptime SLO",
	}, []string{"id", "name", "hostname", "tags"})

	pingdomCheckAvailableErrorBudget = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pingdom_uptime_slo_error_budget_available_seconds",
		Help: "Number of seconds of downtime we can still have without breaking the uptime SLO",
	}, []string{"id", "name", "hostname", "tags"})

	pingdomDownTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pingdom_down_seconds",
		Help: "Total down time within the outage check period, in seconds",
	}, []string{"id", "name", "hostname", "tags"})

	pingdomUpTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "pingdom_up_seconds",
		Help: "Total up time within the outage check period, in seconds",
	}, []string{"id", "name", "hostname", "tags"})
)

func init() {
	flag.IntVar(&waitSeconds, "wait", 60, "time (in seconds) to wait between each metrics update")
	flag.IntVar(&port, "port", 9158, "port to listen on")
	flag.IntVar(&outageCheckPeriod, "outage-check-period", 7, "time (in days) in which to retrieve outage data from the Pingdom API")
	flag.Float64Var(&defaultUptimeSLO, "default-uptime-slo", 99.0, "default uptime SLO to be used when the check doesn't provide a uptime SLO tag (i.e. uptime_slo_999 to 99.9% uptime SLO)")

	prometheus.MustRegister(pingdomUp)
	prometheus.MustRegister(pingdomCheckStatus)
	prometheus.MustRegister(pingdomCheckResponseTime)
	prometheus.MustRegister(pingdomOutageCheckPeriod)
	prometheus.MustRegister(pingdomOutages)
	prometheus.MustRegister(pingdomDownTime)
	prometheus.MustRegister(pingdomUpTime)
	prometheus.MustRegister(pingdomCheckErrorBudget)
	prometheus.MustRegister(pingdomCheckAvailableErrorBudget)
}

func retrieveMetrics(client *pingdom.Client) {
	checks, err := client.Checks.List(map[string]string{
		"include_tags": "true",
	})

	if err != nil {
		log.Errorf("Error getting checks: %v", err)
		pingdomUp.Set(0)
		return
	}
	pingdomUp.Set(1)

	for _, check := range checks {
		id := strconv.Itoa(check.ID)
		tags := check.TagsString()
		resolution := strconv.Itoa(check.Resolution)

		var status float64
		paused := "false"
		if check.Status == "paused" {
			paused = "true"
		} else if check.Status == "up" {
			status = 1
		}

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
		).Set(float64(check.LastResponseTime) / 1000.0)

		retrieveOutagesForCheck(client, check)
	}
}

func retrieveOutagesForCheck(client *pingdom.Client, check pingdom.CheckResponse) {
	var downCount, upTime, downTime float64

	id := strconv.Itoa(check.ID)
	tags := check.TagsString()

	now := time.Now()

	outageCheckPeriodDuration := time.Hour * time.Duration(24*outageCheckPeriod)
	outageCheckPeriodSecs := float64(outageCheckPeriodDuration / time.Second)

	// Register outage check period as a metric
	pingdomOutageCheckPeriod.Set(outageCheckPeriodSecs)

	// Maximum allowed downtime, in seconds, according to the uptime SLO
	uptimeErrorBudget := outageCheckPeriodSecs * (100.0 - check.UptimeSLOFromTags(defaultUptimeSLO)) / 100.0

	// Retrieve the list of outages within the outage period for the given check
	states, err := client.OutageSummary.List(check.ID, map[string]string{
		"from": strconv.FormatInt(now.Add(-outageCheckPeriodDuration).Unix(), 10),
		"to":   strconv.FormatInt(now.Unix(), 10),
	})

	if err != nil {
		log.Errorf("Error getting outages for check %d: %v", check.ID, err)
		return
	}

	for _, state := range states {
		switch state.Status {
		case "down":
			downCount = downCount + 1
			downTime = downTime + float64(state.ToTime-state.FromTime)
		case "up":
			upTime = upTime + float64(state.ToTime-state.FromTime)

		}
	}

	pingdomOutages.WithLabelValues(
		id,
		check.Name,
		check.Hostname,
		tags,
	).Set(downCount)

	pingdomUpTime.WithLabelValues(
		id,
		check.Name,
		check.Hostname,
		tags,
	).Set(upTime)

	pingdomDownTime.WithLabelValues(
		id,
		check.Name,
		check.Hostname,
		tags,
	).Set(downTime)

	pingdomCheckErrorBudget.WithLabelValues(
		id,
		check.Name,
		check.Hostname,
		tags,
	).Set(uptimeErrorBudget)

	pingdomCheckAvailableErrorBudget.WithLabelValues(
		id,
		check.Name,
		check.Hostname,
		tags,
	).Set(uptimeErrorBudget - downTime)
}

func main() {
	var client *pingdom.Client
	flag.Parse()

	token = os.Getenv("PINGDOM_API_TOKEN")
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

	s := NewServer()
	h := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Run metric retrieval loops
	go func() {
		retrieveMetrics(client)

		for {
			select {
			case <-time.After(time.Second * time.Duration(waitSeconds)):
				retrieveMetrics(client)
			}
		}
	}()

	// Run server
	go func() {
		log.Infof("Pingdom Exporter v%s listening on http://0.0.0.0:%d\n", VERSION, port)
		if err := h.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-done
	log.Infof("Received shutdown signal, exiting")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.Shutdown(ctx)
	log.Infoln("Server gracefully stopped")
}
