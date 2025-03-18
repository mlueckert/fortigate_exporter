package probe

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestHAChecksum(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/ha-checksums", "testdata/ha-checksum.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemHAChecksum, c, r) {
		t.Errorf("probeSystemHAChecksum() returned non-success")
	}

	em := `
	# HELP fortigate_ha_member_has_role Master/Slave information
	# TYPE fortigate_ha_member_has_role gauge
	fortigate_ha_member_has_role{role="manage_master", serial="SERIAL111111111"} 1
	fortigate_ha_member_has_role{role="manage_master", serial="SERIAL222222222"} 0
	fortigate_ha_member_has_role{role="root_master", serial="SERIAL111111111"} 1
	fortigate_ha_member_has_role{role="root_master", serial="SERIAL222222222"} 0
	# HELP fortigate_ha_sync_status HA Sync status
	# TYPE fortigate_ha_sync_status gauge
	fortigate_ha_sync_status{serial="SERIAL111111111"} 0
	fortigate_ha_sync_status{serial="SERIAL222222222"} 0
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}

func TestHAChecksumNotInSync(t *testing.T) {
	c := newFakeClient()
	c.prepare("api/v2/monitor/system/ha-checksums", "testdata/ha-checksum-notinsync.jsonnet")
	r := prometheus.NewPedanticRegistry()
	if !testProbe(probeSystemHAChecksum, c, r) {
		t.Errorf("probeSystemHAChecksum() returned non-success")
	}

	em := `
	# HELP fortigate_ha_member_has_role Master/Slave information
	# TYPE fortigate_ha_member_has_role gauge
	fortigate_ha_member_has_role{role="manage_master", serial="SERIAL111111111"} 1
	fortigate_ha_member_has_role{role="manage_master", serial="SERIAL222222222"} 0
	fortigate_ha_member_has_role{role="root_master", serial="SERIAL111111111"} 1
	fortigate_ha_member_has_role{role="root_master", serial="SERIAL222222222"} 0
	# HELP fortigate_ha_sync_status HA Sync status
	# TYPE fortigate_ha_sync_status gauge
	fortigate_ha_sync_status{serial="SERIAL111111111"} 1
	fortigate_ha_sync_status{serial="SERIAL222222222"} 1
	`

	if err := testutil.GatherAndCompare(r, strings.NewReader(em)); err != nil {
		t.Fatalf("metric compare: err %v", err)
	}
}
