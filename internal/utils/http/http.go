package http

import (
	"bytes"
	"fmt"
	"net/http"
)

func MakeRequest(method string, url string, body []byte, headers map[string]string, params map[string]string) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("Error creating request", err.Error())
	}
	for key,value := range headers {
		req.Header.Add(key,value)
	}
	q := req.URL.Query()
	for key,value := range params {
		q.Add(key,value)
	}
	req.URL.RawQuery = q.Encode()
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request", err.Error())
	}
	return response, nil
}

