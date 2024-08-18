package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
)

type BenthosClient struct {
	Host       string
	httpClient *http.Client
}

func NewClient(host string) (*BenthosClient, error) {
	c := &BenthosClient{
		httpClient: &http.Client{},
		Host:       host,
	}
	// TODO: verify connection and authenthication
	return c, nil
}

func (c *BenthosClient) doRequest(req *http.Request, target interface{}, expectedStatus []int) error {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if !slices.Contains(expectedStatus, res.StatusCode) {
		return fmt.Errorf("unexpected status code: %d for request %+v with body %s", res.StatusCode, req)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// return early if target is null
	if target == nil {
		return nil
	}

	if err := json.Unmarshal(body, &target); err != nil {
		return err
	}
	print(fmt.Sprintf("HERE %+v", target))

	return nil
}
