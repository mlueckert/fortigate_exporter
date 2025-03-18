package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestIPv4Routes(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/router/ipv4?vdom=*", "testdata/ipv4-routes.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeIPv4Routes, c, r) {
		t.Errorf("probeIPv4Routes() returned non-success")
	}

	em := `
	# HELP fortigate_ipv4_routes_total Count of ipv4 routes by VDOM
	# TYPE fortigate_ipv4_routes_total gauge
	fortigate_ipv4_routes_total{type="bgp",vdom="globalinet"} 4
	fortigate_ipv4_routes_total{type="connect",vdom="3rdpartyvpn"} 1
	fortigate_ipv4_routes_total{type="connect",vdom="globalinet"} 3
	fortigate_ipv4_routes_total{type="connect",vdom="root"} 3
	fortigate_ipv4_routes_total{type="connect",vdom="wanevalpoc"} 1
	fortigate_ipv4_routes_total{type="connect",vdom="studentenv"} 3
	fortigate_ipv4_routes_total{type="ospf",vdom="globalinet"} 5
	fortigate_ipv4_routes_total{type="static",vdom="3rdpartyvpn"} 3
	fortigate_ipv4_routes_total{type="static",vdom="globalinet"} 4
	fortigate_ipv4_routes_total{type="static",vdom="root"} 1
	fortigate_ipv4_routes_total{type="static",vdom="studentenv"} 4
	fortigate_ipv4_routes_total{type="static",vdom="wanevalpoc"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
