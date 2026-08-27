// Virtual WAN SLA test metrics
//
// Copyright (C) 2020  Christian Svensson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package probe

import (
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// timedMetric wraps a metric with a specific timestamp in seconds
type timedMetric struct {
	m             prometheus.Metric
	timestampSecs int64
}

func (tm *timedMetric) Desc() *prometheus.Desc {
	return tm.m.Desc()
}

func (tm *timedMetric) Write(out *dto.Metric) error {
	if err := tm.m.Write(out); err != nil {
		return err
	}
	// Convert seconds to milliseconds for Prometheus
	timestampMs := tm.timestampSecs * 1000
	out.TimestampMs = &timestampMs
	return nil
}

func probeVirtualWANSLA(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mJitter = prometheus.NewDesc(
			"fortigate_virtual_wan_sla_jitter",
			"Virtual WAN SLA test jitter in milliseconds",
			[]string{"test_name", "interface"}, nil,
		)
		mPacketloss = prometheus.NewDesc(
			"fortigate_virtual_wan_sla_packetloss_ratio",
			"Virtual WAN SLA test packet loss ratio (0-1)",
			[]string{"test_name", "interface"}, nil,
		)
		mLatency = prometheus.NewDesc(
			"fortigate_virtual_wan_sla_latency",
			"Virtual WAN SLA test latency in milliseconds",
			[]string{"test_name", "interface"}, nil,
		)
		mLinkStatus = prometheus.NewDesc(
			"fortigate_virtual_wan_sla_link_status",
			"Virtual WAN SLA test link status (1 for the current state)",
			[]string{"test_name", "interface", "link"}, nil,
		)
	)

	type slaLog struct {
		Timestamp  int64   `json:"timestamp"`
		Link       string  `json:"link"`
		Latency    float64 `json:"latency"`
		Jitter     float64 `json:"jitter"`
		Packetloss float64 `json:"packetloss"`
	}

	type slaResult struct {
		Name      string   `json:"name"`
		Interface string   `json:"interface"`
		Logs      []slaLog `json:"logs"`
	}

	type slaResponse struct {
		Results []slaResult `json:"results"`
	}

	var r []slaResponse

	if err := c.Get("api/v2/monitor/virtual-wan/sla-log", "vdom=*", &r); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}

	for _, resp := range r {
		for _, test := range resp.Results {
			// Export metrics for all log entries
			for _, logEntry := range test.Logs {
				// Create metric with timestamp from the log entry (in seconds)
				jitterMetric := &timedMetric{
					m:             prometheus.MustNewConstMetric(mJitter, prometheus.GaugeValue, logEntry.Jitter, test.Name, test.Interface),
					timestampSecs: logEntry.Timestamp,
				}
				m = append(m, jitterMetric)

				packetlossMetric := &timedMetric{
					m:             prometheus.MustNewConstMetric(mPacketloss, prometheus.GaugeValue, logEntry.Packetloss/100.0, test.Name, test.Interface),
					timestampSecs: logEntry.Timestamp,
				}
				m = append(m, packetlossMetric)

				latencyMetric := &timedMetric{
					m:             prometheus.MustNewConstMetric(mLatency, prometheus.GaugeValue, logEntry.Latency, test.Name, test.Interface),
					timestampSecs: logEntry.Timestamp,
				}
				m = append(m, latencyMetric)

				// Link status: 1 for the current state
				linkStatusMetric := &timedMetric{
					m:             prometheus.MustNewConstMetric(mLinkStatus, prometheus.GaugeValue, 1, test.Name, test.Interface, logEntry.Link),
					timestampSecs: logEntry.Timestamp,
				}
				m = append(m, linkStatusMetric)
			}
		}
	}

	return m, true
}
