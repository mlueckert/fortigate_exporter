package probe

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/bluecmd/fortigate_exporter/internal/config"
)

func TestProbeConcurrentAuthentication(t *testing.T) {
	const target = "https://firewall.example"
	cfg := config.FortiExporterConfig{AuthKeys: config.AuthKeys{
		config.Target(target): {Token: "original-token"},
		"members": {Probes: config.Probes{
			Include: config.ProbeList{"VirtualWAN/Members"},
		}},
		"excluded": {Probes: config.Probes{
			Include: config.ProbeList{"VirtualWAN/Members"},
			Exclude: config.ProbeList{"VirtualWAN/Members"},
		}},
	}}

	t.Run("concurrent scrapes", func(t *testing.T) {
		for i := 0; i < 500; i++ {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				t.Parallel()
				token := "test-token-" + strconv.Itoa(i)
				profile, wantRequests := "members", 2
				if i%2 == 0 {
					profile, wantRequests = "excluded", 1
				}
				requests := 0
				hc := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					requests++
					if req.Header.Get("Authorization") != "Bearer "+token {
						t.Error("request used another scrape's token")
					}
					body := `[]`
					switch req.URL.Path {
					case "/api/v2/monitor/system/status":
						body = `{"status":"success","version":"v7.4.0"}`
					case "/api/v2/monitor/virtual-wan/members":
						if profile == "excluded" {
							t.Error("excluded probe was called")
						}
					default:
						t.Errorf("unexpected probe endpoint: %s", req.URL.Path)
					}
					return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}, nil
				})}
				collector := &ProbeCollector{}
				ok, err := collector.Probe(context.Background(), map[string]string{
					"target": target, "token": token, "profile": profile,
				}, hc, cfg)
				if err != nil || !ok {
					t.Fatalf("Probe() = %v, %v; want success", ok, err)
				}
				if requests != wantRequests {
					t.Errorf("requests = %d, want %d", requests, wantRequests)
				}
			})
		}
	})
	if len(cfg.AuthKeys) != 3 || cfg.AuthKeys[config.Target(target)].Token != "original-token" {
		t.Error("scrapes modified the shared authentication map")
	}
}
