package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestVirtualWANMembers(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/virtual-wan/members", "testdata/virtual_wan_members.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeVirtualWANMembers, c, r) {
		t.Errorf("probeVirtualWANMembers() returned non-success")
	}

	em := `
		# HELP fortigate_virtual_wan_member_receive_bytes_total Number of bytes received on the interface
		# TYPE fortigate_virtual_wan_member_receive_bytes_total counter
		fortigate_virtual_wan_member_receive_bytes_total{name="EMAC_peri_680"} 1.442313770232e+12
		fortigate_virtual_wan_member_receive_bytes_total{name="H-ABCD_ISP1"} 9.66618808307e+11
		# HELP fortigate_virtual_wan_member_transmit_bytes_total Number of bytes transmitted on the interface
		# TYPE fortigate_virtual_wan_member_transmit_bytes_total counter
		fortigate_virtual_wan_member_transmit_bytes_total{name="EMAC_peri_680"} 6.426696498218e+12
		fortigate_virtual_wan_member_transmit_bytes_total{name="H-ABCD_ISP1"} 5.848688522671e+12
		# HELP fortigate_virtual_wan_member_up Whether the link is up (value 1) or not (not taking into account admin status)
		# TYPE fortigate_virtual_wan_member_up gauge
		fortigate_virtual_wan_member_up{name="EMAC_peri_680"} 1
		fortigate_virtual_wan_member_up{name="H-ABCD_ISP1"} 1
	`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
