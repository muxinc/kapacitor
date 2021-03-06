package mux

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/influxdata/kapacitor"
)

var (
	incidentKeyRegExp = regexp.MustCompile("^properties/(.+)/alerts/(.+)/breakdown/(.+)$")
)

type Service struct {
	httpClient *http.Client
	username   string
	password   string
	url        string
	global     bool
	logger     *log.Logger
}

func NewService(c Config, l *log.Logger) *Service {
	return &Service{
		httpClient: &http.Client{},
		username:   c.Username,
		password:   c.Password,
		url:        c.URL,
		global:     c.Global,
		logger:     l,
	}
}

func (s *Service) Open() error {
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Global() bool {
	return s.global
}

func (s *Service) Alert(incidentKey string, level kapacitor.AlertLevel, t time.Time) error {
	incidentBytes, err := buildIncident(incidentKey, level, t)
	if err != nil {
		return err
	}

	post := bytes.NewReader(incidentBytes)

	var fullURL string
	fullURL, err = s.buildIncidentURL(incidentKey)
	if err != nil {
		return err
	}

	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, fullURL, post)
	req.SetBasicAuth(s.username, s.password)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	var resp *http.Response
	resp, err = s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		type response struct {
			Message string `json:"message"`
		}
		r := &response{Message: fmt.Sprintf("failed to understand Mux response. code: %d content: %s", resp.StatusCode, string(body))}
		b := bytes.NewReader(body)
		dec := json.NewDecoder(b)
		dec.Decode(r)
		return errors.New(r.Message)
	}
	return nil
}

func buildIncident(incidentKey string, level kapacitor.AlertLevel, t time.Time) ([]byte, error) {
	if !incidentKeyRegExp.MatchString(incidentKey) {
		return nil, fmt.Errorf("Incident key did not match regular-expression pattern: key = %s", incidentKey)
	}
	// parse incident key for details
	keyParts := incidentKeyRegExp.FindStringSubmatch(incidentKey)

	parent := make(map[string]map[string]interface{})
	pData := make(map[string]interface{})

	// set breakdowns on incident
	breakdownKey := strings.TrimSuffix(keyParts[3], ",")
	breakdowns := strings.Split(breakdownKey, ",")
	breakdownGroups := make([]map[string]string, 0)
	for _, b := range breakdowns {
		breakdownParts := strings.Split(b, "=")
		breakdownGroups = append(breakdownGroups, map[string]string{"name": breakdownParts[0], "value": breakdownParts[1]})
	}

	parent["incident"] = pData
	pData["breakdown_key"] = breakdownKey
	pData["breakdowns"] = breakdownGroups
	switch level {
	case kapacitor.WarnAlert:
		pData["status"] = "open"
		pData["severity"] = "warning"
		pData["started_at"] = t.Format(time.RFC3339)
	case kapacitor.CritAlert:
		pData["status"] = "open"
		pData["severity"] = "alert"
		pData["started_at"] = t.Format(time.RFC3339)
	case kapacitor.InfoAlert:
		return nil, fmt.Errorf("AlertLevel 'info' is currently ignored by the Mux service")
	default:
		pData["status"] = "closed"
		pData["resolved_at"] = t.Format(time.RFC3339)
	}

	return json.Marshal(parent)
}

func (s *Service) buildIncidentURL(incidentKey string) (string, error) {
	if !incidentKeyRegExp.MatchString(incidentKey) {
		return "", fmt.Errorf("Incident key did not match regular-expression pattern: key = %s", incidentKey)
	}
	// parse incident key for details
	keyParts := incidentKeyRegExp.FindStringSubmatch(incidentKey)

	// incidentKey should look like properties/:property_id/alerts/:alert_id
	fullURL := s.url
	if false == strings.HasSuffix(fullURL, "/") {
		fullURL = fullURL + "/"
	}
	fullURL = fullURL + "internal-api/v1/properties/" + keyParts[1] + "/alerts/" + keyParts[2] + "/incident"
	return fullURL, nil
}
