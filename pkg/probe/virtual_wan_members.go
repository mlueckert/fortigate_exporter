package probe

import (
	"bytes"
	"encoding/json"
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
		Interface string  `json:"interface"`
		Link      string  `json:"link"`
		TxBytes   float64 `json:"tx_bytes"`
		RxBytes   float64 `json:"rx_bytes"`
	}

	type VirtualWanMemberResponse struct {
		HTTPMethod string          `json:"http_method"`
		Results    json.RawMessage `json:"results"`
		VDOM       string          `json:"vdom"`
		Path       string          `json:"path"`
		Name       string          `json:"name"`
		Status     string          `json:"status"`
		Serial     string          `json:"serial"`
		Version    string          `json:"version"`
		Build      int64           `json:"build"`
	}

	var rs []VirtualWanMemberResponse

	if err := c.Get("api/v2/monitor/virtual-wan/members", "vdom=*", &rs); err != nil {
		meta.logf("Error: %v", err)
		return nil, false
	}
	m := []prometheus.Metric{}
	for _, r := range rs {
		var members []ifResult
		if len(r.Results) == 0 {
			continue
		}
		if bytes.HasPrefix(bytes.TrimSpace(r.Results), []byte("[")) {
			if err := json.Unmarshal(r.Results, &members); err != nil {
				meta.logf("Error in probe VirtualWAN/Members: invalid results array: %v", err)
				return nil, false
			}
			for i, member := range members {
				if member.Interface == "" {
					meta.logf("Error in probe VirtualWAN/Members: result at index %d has no interface", i)
					return nil, false
				}
			}
		} else {
			var results map[string]ifResult
			if err := json.Unmarshal(r.Results, &results); err != nil {
				meta.logf("Error in probe VirtualWAN/Members: invalid results object: %v", err)
				return nil, false
			}
			for name, member := range results {
				member.Interface = name
				members = append(members, member)
			}
		}
		for _, ir := range members {
			linkf := 0.0
			if strings.ToLower(ir.Link) == "up" {
				linkf = 1.0
			}
			m = append(m, prometheus.MustNewConstMetric(mLink, prometheus.GaugeValue, linkf, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mTxB, prometheus.CounterValue, ir.TxBytes, ir.Interface))
			m = append(m, prometheus.MustNewConstMetric(mRxB, prometheus.CounterValue, ir.RxBytes, ir.Interface))
		}
	}

	return m, true
}
