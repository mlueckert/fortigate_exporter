package probe

import (
	"fmt"
	"log"

	"github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/prometheus/client_golang/prometheus"
)

func probeSystemStatus(c http.FortiHTTP, meta *TargetMetadata) ([]prometheus.Metric, bool) {
	var (
		mVersion = prometheus.NewDesc(
			"fortigate_version_info",
			"System version and build information",
			[]string{"serial", "version", "build", "model_name", "model_number", "model"}, nil,
		)
	)

	type systemStatus struct {
		Status  string
		Serial  string
		Version string
		Build   int64
		Result  struct {
			ModelName   string `json:"model_name"`
			ModelNumber string `json:"model_number"`
			Model       string
		} `json:"results"`
	}
	var st systemStatus

	if err := c.Get("api/v2/monitor/system/status", "", &st); err != nil {
		log.Printf("Error: %v", err)
		return nil, false
	}

	m := []prometheus.Metric{
		prometheus.MustNewConstMetric(mVersion, prometheus.GaugeValue, 1.0, st.Serial, st.Version, fmt.Sprintf("%d", st.Build), st.Result.ModelName, st.Result.ModelNumber, st.Result.Model),
	}
	return m, true
}
