package exporter

import (
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// TODO: customize this
type ServiceStats struct {
	Hostname string `json:"hostname"`
	Name     string `json:"name"`
	Version  string `json:"version"`
}

// FetchStats Fetch base endpoint for internal stats data from the specified service
func FetchStats(uri string) *ServiceStats {
	body, err := get(uri)
	if err != nil {
		return nil
	}

	serviceStats := ServiceStats{}
	if err = json.Unmarshal(body, &serviceStats); err != nil {
		log.Error("Could not parse JSON response from target stats", uri)
		return nil
	}
	return &serviceStats
}

// FetchMetrics Fetch internal metrics from the specified service
func FetchMetrics(uri string) (any, error) {
	var raw any
	body, err := get(uri)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &raw); err != nil {
		log.Error("Could not parse JSON response for target")
		return nil, err
	}
	return raw, nil
}

func get(uri string) ([]byte, error) {
	res, err := http.Get(uri)
	if err != nil {
		log.Errorf("Could not fetch metrics for endpoint of target: %s", uri)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Can't read body of response")
		return nil, err
	}
	return body, nil
}
