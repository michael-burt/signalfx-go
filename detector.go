package signalfx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/signalfx/signalfx-go/detector"
)

// DetectorAPIURL is the base URL for interacting with detectors.
const DetectorAPIURL = "/v2/detector"

// CreateDetector creates a detector.
func (c *Client) CreateDetector(ctx context.Context, detectorRequest *detector.CreateUpdateDetectorRequest) (*detector.Detector, error) {
	payload, err := json.Marshal(detectorRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "POST", DetectorAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalDetector := &detector.Detector{}

	err = json.NewDecoder(resp.Body).Decode(finalDetector)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalDetector, err
}

// DeleteDetector deletes a detector.
func (c *Client) DeleteDetector(ctx context.Context, id string) error {
	resp, err := c.doRequest(ctx, "DELETE", DetectorAPIURL+"/"+id, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return nil
}

// DisableDetector disables a detector.
func (c *Client) DisableDetector(ctx context.Context, id string, labels []string) error {
	payload, err := json.Marshal(labels)
	if err != nil {
		return err
	}

	resp, err := c.doRequest(ctx, "PUT", DetectorAPIURL+"/"+id+"/disable", nil, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return nil
}

// EnableDetector enables a detector.
func (c *Client) EnableDetector(ctx context.Context, id string, labels []string) error {
	payload, err := json.Marshal(labels)
	if err != nil {
		return err
	}

	resp, err := c.doRequest(ctx, "PUT", DetectorAPIURL+"/"+id+"/enable", nil, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return nil
}

// GetDetector gets a detector.
func (c *Client) GetDetector(ctx context.Context, id string) (*detector.Detector, error) {
	resp, err := c.doRequest(ctx, "GET", DetectorAPIURL+"/"+id, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalDetector := &detector.Detector{}

	err = json.NewDecoder(resp.Body).Decode(finalDetector)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalDetector, err
}

// GetDetectors gets all detectors.
func (c *Client) GetDetectors(ctx context.Context, limit int, name string, offset int) ([]*detector.Detector, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	resp, err := c.doRequest(ctx, "GET", DetectorAPIURL, params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	var allDetectors = struct {
		Count   int                  `json:"count"`
		Results []*detector.Detector `json:"results"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&allDetectors)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return allDetectors.Results, err
}

// UpdateDetector updates a detector.
func (c *Client) UpdateDetector(ctx context.Context, id string, detectorRequest *detector.CreateUpdateDetectorRequest) (*detector.Detector, error) {
	payload, err := json.Marshal(detectorRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest(ctx, "PUT", DetectorAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalDetector := &detector.Detector{}

	err = json.NewDecoder(resp.Body).Decode(finalDetector)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalDetector, err
}

// SearchDetectors searches for detectors, given a query string in `name`.
func (c *Client) SearchDetectors(ctx context.Context, limit int, name string, offset int, tags string) (*detector.SearchResults, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	if tags != "" {
		params.Add("tags", tags)
	}

	resp, err := c.doRequest(ctx, "GET", DetectorAPIURL, params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	finalDetectors := &detector.SearchResults{}

	err = json.NewDecoder(resp.Body).Decode(finalDetectors)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return finalDetectors, err
}

// GetDetectorEvents gets a detector's events.
func (c *Client) GetDetectorEvents(ctx context.Context, id string, from int, to int, offset int, limit int) ([]*detector.Event, error) {
	params := url.Values{}
	params.Add("from", strconv.Itoa(from))
	params.Add("to", strconv.Itoa(to))
	params.Add("offset", strconv.Itoa(offset))
	params.Add("limit", strconv.Itoa(limit))
	resp, err := c.doRequest(ctx, "GET", DetectorAPIURL+"/"+id+"/events", params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	var events []*detector.Event

	err = json.NewDecoder(resp.Body).Decode(&events)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return events, err
}

// GetDetectorIncidents gets a detector's incidents.
func (c *Client) GetDetectorIncidents(ctx context.Context, id string, offset int, limit int) ([]*detector.Incident, error) {
	params := url.Values{}
	params.Add("offset", strconv.Itoa(offset))
	params.Add("limit", strconv.Itoa(limit))
	resp, err := c.doRequest(ctx, "GET", DetectorAPIURL+"/"+id+"/incidents", params, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusOK); err != nil {
		return nil, err
	}

	var incidents []*detector.Incident

	err = json.NewDecoder(resp.Body).Decode(&incidents)
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return incidents, err
}

// ValidateDetector validates a detector.
func (c *Client) ValidateDetector(ctx context.Context, detectorRequest *detector.ValidateDetectorRequestModel) error {
	payload, err := json.Marshal(detectorRequest)
	if err != nil {
		return err
	}

	resp, err := c.doRequest(ctx, "POST", DetectorAPIURL+"/validate", nil, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err = newResponseError(resp, http.StatusNoContent); err != nil {
		return err
	}
	_, _ = io.Copy(ioutil.Discard, resp.Body)

	return nil
}
