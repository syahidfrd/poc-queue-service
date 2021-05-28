package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func httpGet(appURL string, params *url.Values, jwtToken string) (map[string]interface{}, int, error) {
	if params == nil {
		params = &url.Values{}
	}

	var fullURL string
	if getParams := params.Encode(); getParams == "" {
		fullURL = appURL
	} else {
		fullURL = fmt.Sprintf("%s?%s", appURL, getParams)
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, -1, err
	}

	if jwtToken != "" {
		token := fmt.Sprintf("Bearer %s", jwtToken)
		req.Header.Set("Authorization", token)
	}

	req.Header.Set("Auth-Token", "6pNPR3qCKApYKvrEbSEN1618729172")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, -1, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}

	var payload map[string]interface{}
	err = json.Unmarshal(respBody, &payload)
	if err != nil {
		return nil, -1, err
	}

	return payload, resp.StatusCode, nil
}

func httpPost(endpoint string, params *url.Values, jwtToken string) (map[string]interface{}, int, error) {
	if params == nil {
		params = &url.Values{}
	}

	contentType := "application/x-www-form-urlencoded"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return nil, -1, err
	}

	if jwtToken != "" {
		token := fmt.Sprintf("Bearer %s", jwtToken)
		req.Header.Set("Authorization", token)
	}

	req.Header.Set("content-type", contentType)
	req.Header.Set("Auth-Token", "6pNPR3qCKApYKvrEbSEN1618729172")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, -1, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}

	var payload map[string]interface{}
	err = json.Unmarshal(respBody, &payload)
	if err != nil {
		return nil, -1, err
	}

	return payload, resp.StatusCode, nil
}