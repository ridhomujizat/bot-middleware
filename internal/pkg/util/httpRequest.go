package util

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

// HTTPClient is an interface for making HTTP requests. This allows for easier testing and mocking.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the default HTTP client used by the package.
var Client HTTPClient = &http.Client{Timeout: 10 * time.Second}

// Get sends a GET request to the specified URL with optional headers and returns the response body and status code.
func HttpGet(url string, headers map[string]string) (string, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, HandleAppError(err, "HTTP Get", "New Request", true)
	}
	for key, value := range headers { // Corrected syntax here
		req.Header.Set(key, value)
	}

	resp, err := Client.Do(req)
	if err != nil {
		return "", 0, HandleAppError(err, "HTTP Get", "Do", true)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, HandleAppError(err, "HTTP Get", "ReadAll", true)
	}

	return string(body), resp.StatusCode, nil
}

// Post sends a POST request to the specified URL with the provided body and optional headers.
// It returns the response body and status code.
func HttpPost(url string, body []byte, headers map[string]string) (string, int, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", 0, HandleAppError(err, "HTTP Post", "New Request", true)
	}
	for key, value := range headers { // Corrected syntax here
		req.Header.Set(key, value)
	}

	resp, err := Client.Do(req)
	if err != nil {
		return "", 0, HandleAppError(err, "HTTP Post", "Do", true)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, HandleAppError(err, "HTTP Post", "ReadAll", true)
	}

	return string(responseBody), resp.StatusCode, nil
}

// Put sends a PUT request to the specified URL with the provided body and optional headers.
// It returns the response body and status code.
func HttpPut(url string, body []byte, headers map[string]string) (string, int, error) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return "", 0, HandleAppError(err, "HTTP Put", "New Request", true)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := Client.Do(req)
	if err != nil {
		return "", 0, HandleAppError(err, "HTTP Put", "Do", true)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, HandleAppError(err, "HTTP Put", "ReadAll", true)
	}

	return string(responseBody), resp.StatusCode, nil
}

// Delete sends a DELETE request to the specified URL with optional headers and returns the response body and status code.
func HttpDelete(url string, headers map[string]string) (string, int, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", 0, HandleAppError(err, "HTTP Delete", "New Request", true)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := Client.Do(req)
	if err != nil {
		return "", 0, HandleAppError(err, "HTTP Delete", "Do", true)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, HandleAppError(err, "HTTP Delete", "ReadAll", true)
	}

	return string(body), resp.StatusCode, nil
}
