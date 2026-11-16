// Tests of fortigate_exporter
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
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/bluecmd/fortigate_exporter/internal/config"
	fortiHTTP "github.com/bluecmd/fortigate_exporter/pkg/http"
	"github.com/google/go-jsonnet"
	"github.com/prometheus/client_golang/prometheus"
)

type preparedResp struct {
	d []byte
	q url.Values
}

type fakeClient struct {
	data map[string][]preparedResp
}

func (c *fakeClient) prepare(path string, jfile string) {
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	vm := jsonnet.MakeVM()
	output, err := vm.EvaluateFile(jfile)
	if err != nil {
		log.Fatalf("Failed to evaluate jsonnet %q: %v", jfile, err)
	}
	c.data[u.Path] = append(c.data[u.Path], preparedResp{
		d: []byte(output),
		q: u.Query(),
	})
}

func (c *fakeClient) Get(path string, query string, obj interface{}) error {
	rs, ok := c.data[path]
	if !ok {
		log.Fatalf("Tried to get unprepared URL %q", path)
	}
	q, err := url.ParseQuery(query)
	if err != nil {
		log.Fatalf("Unable to parse DUT query: %v", err)
	}
alt:
	for _, r := range rs {
		for k, v := range r.q {
			if len(q[k]) == 0 || q[k][0] != v[0] {
				continue alt
			}
		}
		return json.Unmarshal(r.d, obj)
	}
	log.Fatalf("No prepared response matched URL %q, query %q", path, query)
	return nil
}

type Registry interface {
	MustRegister(...prometheus.Collector)
}

type testProbeCollector struct {
	metrics []prometheus.Metric
}

func (p *testProbeCollector) Collect(c chan<- prometheus.Metric) {
	for _, m := range p.metrics {
		c <- m
	}
}

func (p *testProbeCollector) Describe(c chan<- *prometheus.Desc) {
}

func testProbe(pf probeFunc, c fortiHTTP.FortiHTTP, r Registry) bool {
	meta := &TargetMetadata{
		VersionMajor: 7,
		VersionMinor: 0,
	}
	return testProbeWithMetadata(pf, c, meta, r)
}

func testProbeWithMetadata(pf probeFunc, c fortiHTTP.FortiHTTP, meta *TargetMetadata, r Registry) bool {
	m, ok := pf(c, meta)
	if !ok {
		return false
	}
	p := &testProbeCollector{metrics: m}
	r.MustRegister(p)
	return true
}

func newFakeClient() *fakeClient {
	return &fakeClient{data: map[string][]preparedResp{}}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestProbeLogsTargetHostname(t *testing.T) {
	for _, tt := range []struct {
		name        string
		statusCode  int
		statusBody  string
		membersCode int
		membersBody string
		wantLog     string
	}{
		{"connectivity failure", 500, "", 200, "", "API connectivity test failed"},
		{"unsuccessful status", 200, `{"status":"error"}`, 200, "", "API connectivity test returned status"},
		{"invalid version", 200, `{"status":"success","version":"invalid"}`, 200, "", "Failed to parse OS version"},
		{"member request failure", 200, `{"status":"success","version":"v7.4.0"}`, 500, "", "response code was 500"},
		{"invalid members", 200, `{"status":"success","version":"v7.4.0"}`, 200, `[{"results":[{}]}]`, "has no interface"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var logs bytes.Buffer
			original := log.Writer()
			log.SetOutput(&logs)
			defer log.SetOutput(original)

			hc := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				code, body := tt.membersCode, tt.membersBody
				if req.URL.Path == "/api/v2/monitor/system/status" {
					code, body = tt.statusCode, tt.statusBody
				}

				return &http.Response{StatusCode: code, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
			})}
			const target = "https://firewall.example:8443"
			cfg := config.FortiExporterConfig{AuthKeys: config.AuthKeys{
				"profile": {Probes: config.Probes{Include: config.ProbeList{"VirtualWAN/Members"}}},
			}}
			collector := &ProbeCollector{}
			ok, err := collector.Probe(context.Background(), map[string]string{
				"target": target, "token": "test-token", "profile": "profile",
			}, hc, cfg)
			if err != nil || ok {
				t.Fatalf("Probe() = %v, %v; want failure without request error", ok, err)
			}
			if !strings.Contains(logs.String(), `target="firewall.example"`) || !strings.Contains(logs.String(), tt.wantLog) {
				t.Errorf("log = %q; want target hostname and %q", logs.String(), tt.wantLog)
			}
			if strings.Contains(logs.String(), "test-token") {
				t.Errorf("log contains token: %q", logs.String())
			}
		})
	}
}
