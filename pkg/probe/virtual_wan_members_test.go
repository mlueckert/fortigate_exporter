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

func TestVirtualWANMembersArray(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/virtual-wan/members", "testdata/virtual_wan_members_array.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeVirtualWANMembers, c, r) {
		t.Fatal("probeVirtualWANMembers() returned non-success")
	}

	em := `
			# HELP fortigate_virtual_wan_member_receive_bytes_total Number of bytes received on the interface
			# TYPE fortigate_virtual_wan_member_receive_bytes_total counter
			fortigate_virtual_wan_member_receive_bytes_total{name="wan1"} 7201799146898
			fortigate_virtual_wan_member_receive_bytes_total{name="wan2"} 24
			# HELP fortigate_virtual_wan_member_transmit_bytes_total Number of bytes transmitted on the interface
			# TYPE fortigate_virtual_wan_member_transmit_bytes_total counter
			fortigate_virtual_wan_member_transmit_bytes_total{name="wan1"} 11719487983535
			fortigate_virtual_wan_member_transmit_bytes_total{name="wan2"} 42
			# HELP fortigate_virtual_wan_member_up Whether the link is up (value 1) or not (not taking into account admin status)
			# TYPE fortigate_virtual_wan_member_up gauge
			fortigate_virtual_wan_member_up{name="wan1"} 1
			fortigate_virtual_wan_member_up{name="wan2"} 0
		`
	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestVirtualWANMembersInvalidResults(t *testing.T) {
	for _, tt := range []struct {
		name    string
		results string
	}{
		{"missing interface", `[{"link":"up"}]`},
		{"invalid array entry", `[{"interface":"wan1","tx_bytes":"invalid"}]`},
		{"invalid object entry", `{"wan1":{"tx_bytes":"invalid"}}`},
		{"invalid results type", `"invalid"`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := newFakeClient()
			c.data["api/v2/monitor/virtual-wan/members"] = []preparedResp{
				{d: []byte(`[{"results":` + tt.results + `}]`)},
			}
			if metrics, ok := probeVirtualWANMembers(c, &TargetMetadata{}); ok || len(metrics) != 0 {
				t.Fatalf("probeVirtualWANMembers() = %v, %v; want no metrics and failure", metrics, ok)
			}
		})
	}
}
