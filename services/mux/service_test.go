package mux

import (
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/influxdata/kapacitor"
)

func TestServiceBuildIncidentURL(t *testing.T) {
	tests := []struct {
		// Test description.
		name string
		// Receiver fields.
		rhttpClient *http.Client
		rusername   string
		rpassword   string
		rurl        string
		rglobal     bool
		rlogger     *log.Logger
		// Parameters.
		incidentKey string
		// Expected results.
		want    string
		wantErr bool
	}{
		{
			name:        "Single breakdown",
			rurl:        "http://example.com",
			incidentKey: "properties/4/alerts/38/breakdown/country=US,",
			want:        "http://example.com/internal-api/v1/properties/4/alerts/38/incident",
		},
		{
			name:        "Single breakdown with endpoint having trailing slash",
			rurl:        "http://example.com/",
			incidentKey: "properties/4/alerts/38/breakdown/country=US,",
			want:        "http://example.com/internal-api/v1/properties/4/alerts/38/incident",
		},
		{
			name:        "Malformed incident key",
			rurl:        "http://example.com",
			incidentKey: "foobar",
			want:        "",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		s := &Service{
			httpClient: tt.rhttpClient,
			username:   tt.rusername,
			password:   tt.rpassword,
			url:        tt.rurl,
			global:     tt.rglobal,
			logger:     tt.rlogger,
		}
		got, err := s.buildIncidentURL(tt.incidentKey)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Service.buildIncidentURL() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. Service.buildIncidentURL() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestBuildIncident(t *testing.T) {
	utcLocation, _ := time.LoadLocation("UTC")
	tests := []struct {
		// Test description.
		name string
		// Parameters.
		incidentKey string
		level       kapacitor.AlertLevel
		t           time.Time
		// Expected results.
		want    string
		wantErr bool
	}{
		{
			name:        "Open Warning, Single breakdown",
			incidentKey: "properties/4/alerts/38/breakdown/country=US,",
			level:       kapacitor.WarnAlert,
			t:           time.Unix(1468072740, 0).In(utcLocation),
			want:        "{\"incident\":{\"breakdown_key\":\"country=US\",\"breakdowns\":[{\"name\":\"country\",\"value\":\"US\"}],\"severity\":\"warning\",\"started_at\":\"2016-07-09T13:59:00Z\",\"status\":\"open\"}}",
		},
		{
			name:        "Open Warning, Two breakdowns",
			incidentKey: "properties/4/alerts/38/breakdown/country=US,browser=Firefox",
			level:       kapacitor.WarnAlert,
			t:           time.Unix(1468072740, 0).In(utcLocation),
			want:        "{\"incident\":{\"breakdown_key\":\"country=US,browser=Firefox\",\"breakdowns\":[{\"name\":\"country\",\"value\":\"US\"},{\"name\":\"browser\",\"value\":\"Firefox\"}],\"severity\":\"warning\",\"started_at\":\"2016-07-09T13:59:00Z\",\"status\":\"open\"}}",
		},
		{
			name:        "Open Critical Alert, Single breakdown",
			incidentKey: "properties/4/alerts/38/breakdown/country=US,",
			level:       kapacitor.CritAlert,
			t:           time.Unix(1468072740, 0).In(utcLocation),
			want:        "{\"incident\":{\"breakdown_key\":\"country=US\",\"breakdowns\":[{\"name\":\"country\",\"value\":\"US\"}],\"severity\":\"alert\",\"started_at\":\"2016-07-09T13:59:00Z\",\"status\":\"open\"}}",
		},
		{
			name:        "Open Critical Alert, Two breakdowns",
			incidentKey: "properties/4/alerts/38/breakdown/country=US,browser=Firefox",
			level:       kapacitor.CritAlert,
			t:           time.Unix(1468072740, 0).In(utcLocation),
			want:        "{\"incident\":{\"breakdown_key\":\"country=US,browser=Firefox\",\"breakdowns\":[{\"name\":\"country\",\"value\":\"US\"},{\"name\":\"browser\",\"value\":\"Firefox\"}],\"severity\":\"alert\",\"started_at\":\"2016-07-09T13:59:00Z\",\"status\":\"open\"}}",
		},
		{
			name:        "Closed, Single breakdown",
			incidentKey: "properties/4/alerts/38/breakdown/country=US,",
			level:       kapacitor.OKAlert,
			t:           time.Unix(1468072740, 0).In(utcLocation),
			want:        "{\"incident\":{\"breakdown_key\":\"country=US\",\"breakdowns\":[{\"name\":\"country\",\"value\":\"US\"}],\"resolved_at\":\"2016-07-09T13:59:00Z\",\"status\":\"closed\"}}",
		},
		{
			name:        "Closed, Two breakdowns",
			incidentKey: "properties/4/alerts/38/breakdown/country=US,browser=Firefox",
			level:       kapacitor.OKAlert,
			t:           time.Unix(1468072740, 0).In(utcLocation),
			want:        "{\"incident\":{\"breakdown_key\":\"country=US,browser=Firefox\",\"breakdowns\":[{\"name\":\"country\",\"value\":\"US\"},{\"name\":\"browser\",\"value\":\"Firefox\"}],\"resolved_at\":\"2016-07-09T13:59:00Z\",\"status\":\"closed\"}}",
		},
		{
			name:        "Malformed Incident Key",
			incidentKey: "foobar",
			level:       kapacitor.WarnAlert,
			t:           time.Unix(1468072740, 0).In(utcLocation),
			want:        "",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		got, err := buildIncident(tt.incidentKey, tt.level, tt.t)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. buildIncident() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		gotString := string(got)
		if !reflect.DeepEqual(gotString, tt.want) {
			t.Errorf("%q. buildIncident() = %v, want %v", tt.name, gotString, tt.want)
		}
	}
}
