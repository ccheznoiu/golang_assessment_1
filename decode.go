package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var client = http.Client{Timeout: time.Minute}

func get(url string, params map[string]int) io.ReadCloser {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range params {
		q := req.URL.Query()
		q.Set(k, strconv.Itoa(v))
		req.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else if code := resp.StatusCode; code != http.StatusOK {
		log.Fatalf("%s returned a %d", url, code)
	}

	return resp.Body
}

func decode(rb io.Reader) *json.Decoder {
	dec := json.NewDecoder(rb)

	if _, err := dec.Token(); err != nil {
		log.Fatal(err)
	}

	return dec
}
