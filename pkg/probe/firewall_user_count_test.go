package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestFirewallUserCount(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/user/firewall/count", "testdata/fw-user-count.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeFirewallUserCount, c, r) {
		t.Errorf("probeFirewallUserCount() returned non-success")
	}

	em := `
	# HELP fortigate_fw_user_count Number of authenticated firewall users by VDOM
	# TYPE fortigate_fw_user_count gauge
	fortigate_fw_user_count{vdom="perimeter"} 142
	fortigate_fw_user_count{vdom="root"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
