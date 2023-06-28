package cfgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	baseURL = "https://api.cloudflare.com/client/v4"

	kContentType = "Content-Type"
	kAuthKey     = "X-Auth-Key"
	kAuthEmail   = "X-Auth-Email"

	defaultContentType = "application/json"
)

// do a request with query string
func (c *CloudflareClient) _query(method, endpoint string, params map[string]any) (response []byte, err error) {
	if params == nil {
		params = map[string]any{}
	}

	apiURL := fmt.Sprintf("%s/%s", baseURL, endpoint)

	var req *http.Request
	if req, err = http.NewRequest(method, apiURL, nil); err == nil {
		// parameters
		queries := req.URL.Query()
		for k, v := range params {
			queries.Add(k, fmt.Sprintf("%+v", v))
		}
		req.URL.RawQuery = queries.Encode()

		// authentication headers
		req.Header.Set(kAuthEmail, c.Email)
		req.Header.Set(kAuthKey, c.APIKey)
		req.Header.Set(kContentType, defaultContentType) // set content-type header

		if c.Verbose {
			if dumped, err := httputil.DumpRequest(req, true); err == nil {
				log.Printf("dump request:\n\n%s", string(dumped))
			}
		}

		req.Close = true

		// send request and return response bytes
		var resp *http.Response
		resp, err = c.httpClient.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err == nil {
			if response, err = io.ReadAll(resp.Body); err == nil {
				if c.Verbose {
					log.Printf("API response for %s: '%s'", endpoint, string(response))
				}

				if resp.StatusCode != 200 {
					err = fmt.Errorf("http status %d", resp.StatusCode)
				}

				return response, err
			}
		}
	}

	return nil, err
}

// sends a HTTP GET request
func (c *CloudflareClient) get(endpoint string, params map[string]any) (response []byte, err error) {
	return c._query(http.MethodGet, endpoint, params)
}

// sends a HTTP DELETE request
func (c *CloudflareClient) delete(endpoint string, params map[string]any) (response []byte, err error) {
	return c._query(http.MethodDelete, endpoint, params)
}

// do a request with JSON body
func (c *CloudflareClient) _json(method, endpoint string, params any) (response []byte, err error) {
	if params == nil {
		params = struct{}{}
	}

	apiURL := fmt.Sprintf("%s/%s", baseURL, endpoint)

	var req *http.Request

	// application/json
	var serialized []byte
	if serialized, err = json.Marshal(params); err == nil {
		if req, err = http.NewRequest(method, apiURL, bytes.NewBuffer(serialized)); err != nil {
			return nil, fmt.Errorf("failed to create application/json request: %s", err)
		}

		// authentication headers
		req.Header.Set(kAuthEmail, c.Email)
		req.Header.Set(kAuthKey, c.APIKey)
		req.Header.Set(kContentType, defaultContentType) // set content-type header
	}

	if c.Verbose {
		if dumped, err := httputil.DumpRequest(req, true); err == nil {
			log.Printf("dump request:\n\n%s", string(dumped))
		}
	}
	req.Close = true

	// send request and return response bytes
	var resp *http.Response
	resp, err = c.httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err == nil {
		if response, err = io.ReadAll(resp.Body); err == nil {
			if c.Verbose {
				log.Printf("API response for %s: '%s'", endpoint, string(response))
			}

			if resp.StatusCode != 200 {
				err = fmt.Errorf("http status %d", resp.StatusCode)
			}

			return response, err
		}
	}

	return nil, err
}

// sends a HTTP POST request
func (c *CloudflareClient) post(endpoint string, params any) (response []byte, err error) {
	return c._json(http.MethodPost, endpoint, params)
}

// sends a HTTP PUT request
func (c *CloudflareClient) put(endpoint string, params any) (response []byte, err error) {
	return c._json(http.MethodPut, endpoint, params)
}
