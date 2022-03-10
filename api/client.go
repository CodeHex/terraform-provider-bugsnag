package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tomnomnom/linkheader"
)

// Client defines the HTTP client used communicate with the Bugsnag data access API. The client
// is restricted to single organization. See https://bugsnagapiv2.docs.apiary.io for more details
type Client struct {
	endPoint  string
	authToken string
	OrgID     string
}

// New creates a new Bugsnag data access API client
func New(token string, endPoint string) (*Client, error) {
	client := &Client{authToken: token, endPoint: endPoint}
	orgID, err := client.GetCurrentOrganization()
	if err != nil {
		return nil, err
	}
	client.OrgID = orgID
	return client, nil
}

func (c *Client) createRequest(verb string, path string, body []byte) (*http.Request, error) {
	if !strings.HasPrefix(path, c.endPoint) {
		path = c.endPoint + "/" + path
	}
	req, err := http.NewRequest(verb, path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Version", "2")
	req.Header.Set("Authorization", " token "+c.authToken)
	if len(body) != 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func (c *Client) callAPI(verb string, path string, body []byte, v interface{}, expCode int) error {
	_, err := c.callPagedAPI(verb, path, body, v, expCode)
	return err
}

func (c *Client) callPagedAPI(verb string, path string, body []byte, v interface{}, expCode int) (*url.URL, error) {
	req, err := c.createRequest(verb, path, body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	var resp *http.Response
	success := false
	for !success {
		resp, err = c.attemptCall(client, req)
		if err != nil {
			return nil, err
		}
		success = resp != nil
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != expCode {
		return nil, fmt.Errorf("invalid response code %s for request %s, %s", http.StatusText(resp.StatusCode), path, string(body))
	}
	if v != nil {
		err := json.Unmarshal(respBody, v)
		if err != nil {
			return nil, err
		}
		links := linkheader.Parse(resp.Header.Get("Link"))
		if len(links) > 0 && links[0].Rel == "next" {
			return url.Parse(links[0].URL)
		}
	}
	return nil, nil
}

func (c *Client) attemptCall(client *http.Client, req *http.Request) (*http.Response, error) {
	var err error
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusTooManyRequests && resp.Header.Get("Retry-After") != "" {
		resp.Body.Close()
		waitSec, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			return nil, err
		}
		time.Sleep(time.Duration(waitSec) * time.Second)
		return nil, nil
	}
	return resp, nil
}
