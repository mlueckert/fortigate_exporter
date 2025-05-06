package probe

import (
	"log"
	"reflect"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

type HAChecksumResults struct {
	IsManageMaster int    `json:"is_manage_master"`
	IsRootMaster   int    `json:"is_root_master"`
	SerialNo       string `json:"serial_no"`
	Checksum       struct {
		All string
	}
}

type HAChecksum struct {
	Results []HAChecksumResults `json:"results"`
}

func probeSystemHAChecksum(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		IsMaster = prometheus.NewDesc(
			"fortigate_ha_member_has_role",
			"Master/Slave information",
			[]string{"role", "serial"}, nil,
		)
		SyncStatus = prometheus.NewDesc(
			"fortigate_ha_sync_status",
			"HA Sync status (0=In sync, 1=Not in sync)",
			[]string{"serial"}, nil,
		)
	)

	var res HAChecksum
	if err := c.Get("api/v2/monitor/system/ha-checksums", "scope=global", &res); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{}
	checksums := []string{}
	for _, ele := range res.Results {
		checksums = append(checksums, ele.Checksum.All)
	}
	areAllValuesEqual := areAllValuesEqual(checksums)
	for _, response := range res.Results {
		m = append(m, prometheus.MustNewConstMetric(IsMaster, prometheus.GaugeValue, float64(response.IsManageMaster), "manage_master", response.SerialNo))
		m = append(m, prometheus.MustNewConstMetric(IsMaster, prometheus.GaugeValue, float64(response.IsRootMaster), "root_master", response.SerialNo))
		m = append(m, prometheus.MustNewConstMetric(SyncStatus, prometheus.GaugeValue, float64(areAllValuesEqual), response.SerialNo))
	}

	return m, true
}

// areAllValuesEqual checks if all elements in an array or slice are equal.
//
// The function uses reflection to handle arrays and slices of various data types.
// It panics if the input is not an array or slice. An empty or single-element
// array/slice is considered to have all equal values, returning 0 in such cases.
//
// Parameters:
//
//	arr interface{}: An array or slice of any data type.
//
// Returns:
//
//	int: 0 if all elements are equal, 1 if any element is different.
//
// Example:
//
//	arr := []int{1, 1, 1, 1}
//	result := areAllValuesEqual(arr) // result will be 0
//
//	arr2 := []string{"a", "a", "b"}
//	result2 := areAllValuesEqual(arr2) // result2 will be 1
//
//	arr3 := []float64{}
//	result3 := areAllValuesEqual(arr3) // result3 will be 0
//
// Panics:
//
//	If the input 'arr' is not an array or slice.
func areAllValuesEqual(arr interface{}) int {
	val := reflect.ValueOf(arr)

	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		panic("Input must be an array or a slice")
	}

	if val.Len() <= 1 {
		return 0 // An empty or single-element array/slice is considered equal
	}

	first := val.Index(0).Interface() // Get the value of the first element

	for i := 1; i < val.Len(); i++ {
		if val.Index(i).Interface() != first {
			return 1 // Found a different value
		}
	}

	return 0 // All values are equal
}
