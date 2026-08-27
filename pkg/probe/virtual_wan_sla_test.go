package probe

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestVirtualWANSLA(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/virtual-wan/sla-log", "testdata/virtual-wan-sla-log.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeVirtualWANSLA, c, r) {
		t.Errorf("probeVirtualWANSLA() returned non-success")
	}

	// Verify probe ran and collected metrics
	families, err := r.Gather()
	if err != nil {
		t.Fatalf("failed to gather metrics: %v", err)
	}

	if len(families) == 0 {
		t.Errorf("no metrics collected")
		return
	}

	// Find our metrics
	jitterFound := false
	packetlossFound := false

	for _, family := range families {
		if family.GetName() == "fortigate_virtual_wan_sla_jitter" {
			jitterFound = true
			// Should have 4 tests * 121 log entries = 484 metrics
			if len(family.Metric) < 100 {
				t.Errorf("expected many jitter metrics, got %d", len(family.Metric))
			}
			// Verify first metric has correct labels
			if len(family.Metric[0].Label) < 2 {
				t.Errorf("expected at least 2 labels, got %d", len(family.Metric[0].Label))
			}
			// Verify timestamp is set
			if family.Metric[0].TimestampMs == nil {
				t.Errorf("expected timestamp to be set")
			}
		}
		if family.GetName() == "fortigate_virtual_wan_sla_packetloss_ratio" {
			packetlossFound = true
			// Should have 4 tests * 121 log entries = 484 metrics
			if len(family.Metric) < 100 {
				t.Errorf("expected many packetloss metrics, got %d", len(family.Metric))
			}
		}
	}

	if !jitterFound {
		t.Errorf("jitter metric not found")
	}
	if !packetlossFound {
		t.Errorf("packetloss metric not found")
	}
}
