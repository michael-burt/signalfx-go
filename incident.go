package signalfx

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/signalfx/signalfx-go/detector"
)

// IncidentAPIURL is the base URL for interacting with alert muting rules.
const IncidentAPIURL = "/v2/incident"

// Get incident with the given id
func (c *Client) GetIncident(ctx context.Context, id string) (*detector.Incident, error) {
	resp, err := c.doRequest(ctx, "GET", IncidentAPIURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalIncident := &detector.Incident{}

	err = json.NewDecoder(resp.Body).Decode(finalIncident)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalIncident, err
}

// Get all incidents
func (c *Client) GetIncidents(ctx context.Context, includeResolved bool, limit int, query string, offset int) ([]*detector.Incident, error) {
	params := url.Values{}
	params.Add("includeResolved", strconv.FormatBool(includeResolved))
	params.Add("limit", strconv.Itoa(limit))
	params.Add("query", query)
	params.Add("offset", strconv.Itoa(offset))
	resp, err := c.doRequest(ctx, "GET", IncidentAPIURL, params, nil)
	if err != nil {
		return nil, err
	}

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	var incidents []*detector.Incident
	err = json.NewDecoder(resp.Body).Decode(&incidents)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return incidents, err
}
