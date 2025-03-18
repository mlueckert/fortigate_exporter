package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestSystemStatus(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/status", "testdata/status.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemStatus, c, r) {
		t.Errorf("probeSystemStatus() returned non-success")
	}

	em := `
	# HELP fortigate_version_info System version and build information
	# TYPE fortigate_version_info gauge
	fortigate_version_info{build="1112",model="Fortigate",model_name="FGMODEL1",model_number="162312",serial="FGVMEVZFNTS3OAC8",version="v6.2.4"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
