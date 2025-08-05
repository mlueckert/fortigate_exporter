package probe

import (
	"log"
	"strings"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeVirtualWANMembers(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mLink = prometheus.NewDesc(
			"fortigate_virtual_wan_member_up",
			"Whether the link is up (value 1) or not (not taking into account admin status)",
			[]string{"name"}, nil,
		)
		mTxB = prometheus.NewDesc(
			"fortigate_virtual_wan_member_transmit_bytes_total",
			"Number of bytes transmitted on the interface",
			[]string{"name"}, nil,
		)
		mRxB = prometheus.NewDesc(
			"fortigate_virtual_wan_member_receive_bytes_total",
			"Number of bytes received on the interface",
			[]string{"name"}, nil,
		)
	)

	type ifResult struct {
		Link    string
		TxBytes float64 `json:"tx_bytes"`
		RxBytes float64 `json:"rx_bytes"`
	}

	type VirtualWanMemberResponse struct {
		HTTPMethod string              `json:"http_method"`
		Results    map[string]ifResult `json:"results"`
		VDOM       string              `json:"vdom"`
		Path       string              `json:"path"`
		Name       string              `json:"name"`
		Status     string              `json:"status"`
		Serial     string              `json:"serial"`
		Version    string              `json:"version"`
		Build      int64               `json:"build"`
	}

	var rs []VirtualWanMemberResponse

	if err := c.Get("api/v2/monitor/virtual-wan/members", "vdom=*", &rs); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}
	m := []prometheus.Metric{}
	for _, r := range rs {
		for name, ir := range r.Results {
			linkf := 0.0
			if strings.ToLower(ir.Link) == "up" {
				linkf = 1.0
			}
			m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, linkf, name))
			m = append(m, prometheus.MustNewConstMetric(mTxB, prometheus.CounterValue, ir.TxBytes, name))
			m = append(m, prometheus.MustNewConstMetric(mRxB, prometheus.CounterValue, ir.RxBytes, name))
		}
	}

	return m, true
}
