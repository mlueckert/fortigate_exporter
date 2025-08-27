package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type UserCount struct {
	Results struct {
		IPv4  float64 `json:"ipv4"`
		Total float64 `json:"total"`
	} `json:"results"`
	VDOM string `json:"vdom"`
}

func probeFirewallUserCount(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		FWUserCount = prometheus.NewDesc(
			"fortigate_fw_user_count",
			"Number of authenticated firewall users by VDOM",
			[]string{"vdom"}, nil,
		)
	)

	var res []UserCount
	if err := c.Get("api/v2/monitor/user/firewall/count", "vdom=*", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	for _, v := range res {
		m = append(m, prometheus.MustNewConstMetric(FWUserCount, prometheus.GaugeValue, v.Results.Total, v.VDOM))

	}
	return m, true
}
