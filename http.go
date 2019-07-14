package main

import (
	"errors"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"time"
)

func fetchMenuImage() (io.ReadCloser, error) {
	client := newHTTPClient(10 * time.Second)

	url, err := fetchMenuImageURL(client)
	if err != nil {
		return nil, fmt.Errorf("failed to determine menu image URL: %v", err)
	}

	response, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch web page: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP server responded with status code %d", response.StatusCode)
	}

	return response.Body, nil
}

func fetchMenuImageURL(client http.Client) (string, error) {
	response, err := client.Get("http://www.schachcafe-hamburg.de/mittagstisch/")
	if err != nil {
		return "", fmt.Errorf("failed to fetch web page: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP server responded with status code %d", response.StatusCode)
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return findMenuImage(content)
}

func findMenuImage(content []byte) (string, error) {
	re := regexp.MustCompile(`property="og:image" content="([^"]+)"`)
	match := re.FindSubmatch(content)
	if match == nil {
		return "", errors.New("could not find og:image tag")
	}

	u := string(match[1])

	return html.UnescapeString(u), nil
}

func newHTTPClient(timeout time.Duration) http.Client {
	client := http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: timeout,
			}).Dial,
			TLSHandshakeTimeout: timeout,
		},
	}

	return client
}
