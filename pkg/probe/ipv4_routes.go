package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type IPv4Route struct {
	Type string `json:"type"`
}

type IPv4Results struct {
	Results []IPv4Route `json:"results"`
	VDOM    string      `json:"vdom"`
	Version string      `json:"version"`
}

type RouteCount struct {
	VDOM string
	Type string
}

func probeIPv4Routes(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		IPv4Routes = prometheus.NewDesc(
			"fortigate_ipv4_routes_total",
			"Count of ipv4 routes by VDOM",
			[]string{"vdom", "type"}, nil,
		)
	)

	var rs []IPv4Results

	if err := c.Get("api/v2/monitor/router/ipv4", "vdom=*", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	rcMap := make(map[RouteCount]int)
	for _, r := range rs {
		for _, route := range r.Results {
			sr := RouteCount{
				VDOM: r.VDOM,
				Type: route.Type,
			}
			rcMap[sr] += 1
		}
	}

	for route, count := range rcMap {
		m = append(m, prometheus.MustNewConstMetric(IPv4Routes, prometheus.GaugeValue, float64(count), route.VDOM, route.Type))
	}

	return m, true
}
