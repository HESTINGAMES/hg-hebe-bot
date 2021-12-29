package apiclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var client = &http.Client{}

func DoRequest(verb, url string) ([]byte, error) {
	request, err := http.NewRequest(verb, url, nil)
	if err != nil {
		// Internal error
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		// Internal error
		return nil, err
	}

	if response.StatusCode >= 300 {
		// HTTP error
		return nil, fmt.Errorf("http code is not 2xx")
	}

	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	return body, err
}
