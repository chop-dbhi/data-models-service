package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"text/template"
	"time"
)

const DefaultServiceURL = "https://data-models-service.research.chop.edu"

// Set of URL templates for resources.
var resourceURLs = map[string]string{
	"index":          "/",
	"repos":          "/repos",
	"models":         "/models",
	"modelRevisions": "/models/{{.name}}",
	"modelRevision":  "/models/{{.name}}/{{.version}}",
	"schema":         "/schemata/{{.name}}/{{.version}}",
}

var resourceTemplates = template.New("resources")

func init() {
	for name, url := range resourceURLs {
		template.Must(resourceTemplates.New(name).Parse(url))
	}
}

// Client provides methods for accessing resources from a Data Models service.
type Client struct {
	URL     string
	Timeout time.Duration

	url  *url.URL
	http *http.Client
}

func (c *Client) resource(name string, params map[string]interface{}) *url.URL {
	var buf bytes.Buffer

	err := resourceTemplates.ExecuteTemplate(&buf, name, params)

	if err != nil {
		log.Panicf("resource %s params expansion failed %v", name, params)
	}

	path := buf.String()

	// Copy the URL an set an explicit path.
	url := *c.url
	url.Path = path

	return &url
}

// Build a request to be sent.
func (c *Client) request(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	// Accept JSON.
	req.Header.Set("Accept", "application/json")

	req.URL.Query().Set("format", "json")

	return req, nil
}

// Send sends a request ensuring the timeout is set.
func (c *Client) send(r *http.Request) (*http.Response, error) {
	c.http.Timeout = c.Timeout
	resp, err := c.http.Do(r)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp, nil
	}

	return resp, fmt.Errorf("service %s responded with status code %d", c.URL, resp.StatusCode)
}

// Ping sends a GET request to the service endpoint to ensure it
// is available.
func (c *Client) Ping() error {
	url := c.resource("index", nil)

	req, err := c.request("GET", url.String())

	if err != nil {
		return err
	}

	_, err = c.send(req)

	if err != nil {
		return err
	}

	return nil
}

// Models returns all models available by the data service.
func (c *Client) Models() (*Models, error) {
	url := c.resource("models", nil)

	req, err := c.request("GET", url.String())

	if err != nil {
		return nil, err
	}

	resp, err := c.send(req)

	if err != nil {
		return nil, err
	}

	var models *Models

	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, fmt.Errorf("error decoding models")
	}

	return models, nil
}

// ModelRevisions returns all revisions of a model.
func (c *Client) ModelRevisions(name string) (*Models, error) {
	url := c.resource("modelRevisions", map[string]interface{}{
		"name": name,
	})

	req, err := c.request("GET", url.String())

	if err != nil {
		return nil, err
	}

	resp, err := c.send(req)

	if err != nil {
		return nil, err
	}

	var models *Models

	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, fmt.Errorf("error decoding model revisions: %s", err)
	}

	return models, nil
}

// ModelVersion returns all revisions of a model.
func (c *Client) ModelRevision(name, version string) (*Model, error) {
	url := c.resource("modelRevision", map[string]interface{}{
		"name":    name,
		"version": version,
	})

	req, err := c.request("GET", url.String())

	if err != nil {
		return nil, err
	}

	resp, err := c.send(req)

	if err != nil {
		return nil, err
	}

	var model Model

	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&model); err != nil {
		return nil, fmt.Errorf("error decoding model: %s", err)
	}

	return &model, nil
}

// Schema returns the schema information of a model.
func (c *Client) Schema(name, version string) (*Schema, error) {
	url := c.resource("schema", map[string]interface{}{
		"name":    name,
		"version": version,
	})

	req, err := c.request("GET", url.String())

	if err != nil {
		return nil, err
	}

	resp, err := c.send(req)

	if err != nil {
		return nil, err
	}

	var aux map[string]json.RawMessage

	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&aux); err != nil {
		return nil, fmt.Errorf("error decoding schema: %s", err)
	}

	var schema Schema

	if err = json.Unmarshal(aux["schema"], &schema); err != nil {
		return nil, fmt.Errorf("error decoding schema: %s", err)
	}

	return &schema, nil
}

// New initializes a new client to the
func New(service string) (*Client, error) {
	if service == "" {
		service = DefaultServiceURL
	}

	// Problem parsing the service URL.
	purl, err := url.ParseRequestURI(service)

	if err != nil {
		return nil, err
	}

	c := Client{
		URL:     service,
		Timeout: time.Second * 5,
		url:     purl,
		http: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	return &c, nil
}
