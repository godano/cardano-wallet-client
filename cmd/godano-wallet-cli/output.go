package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func outputResponse(response *http.Response) {
	success := response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices
	if success {
		Log.Debugf("Response status: %v", response.Status)
	} else {
		Log.Infof("Response status: %v", response.Status)
	}

	var body interface{}
	bodyContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Log.Errorf("Failed to read response body: %v", err)
		return
	}

	err = json.Unmarshal(bodyContent, &body)
	if err != nil {
		Log.Errorf("Failed to unmarshal response body: %v", err)
		outputBuffer(bodyContent)
		return
	}

	outputObjectAsJSON(body)
}

func outputObjectAsJSON(obj interface{}) {
	marshalled, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		Log.Errorf("Failed to JSON-marshal object: %v", err)
		return
	}
	outputBuffer(marshalled)
}

func outputBuffer(buf []byte) {
	scanner := bufio.NewScanner(bytes.NewReader(buf))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		Log.Info(scanner.Text())
	}
}
