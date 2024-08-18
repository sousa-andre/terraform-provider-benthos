package client

import (
	"fmt"
	"net/http"
	"strings"
)

type Stream struct {
	Active    bool    `json:"active"`
	Uptime    float64 `json:"uptime"`
	UptimeStr string  `json:"uptime_str"`
}

type DetailedStream struct {
	*Stream
	Configuration map[string]any `json:"config"`
}

func (c *BenthosClient) GetStreams() ([]Stream, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/streams", c.Host), nil)
	if err != nil {
		return nil, err
	}

	var streams []Stream
	if err := c.doRequest(req, &streams, []int{http.StatusOK}); err != nil {
		return nil, err
	}
	return streams, nil
}

func (c *BenthosClient) GetStream(streamId string) (*DetailedStream, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/streams/%s", c.Host, streamId), nil)
	if err != nil {
		return nil, err
	}

	var detailedStream *DetailedStream
	if err := c.doRequest(req, &detailedStream, []int{http.StatusOK}); err != nil {
		return nil, err
	}

	print(fmt.Sprintf("BA %+v", detailedStream))

	return detailedStream, nil
}

func (c *BenthosClient) CreateStream(streamId string, stream string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/streams/%s", c.Host, streamId),
		strings.NewReader(stream))
	if err != nil {
		return err
	}

	if err := c.doRequest(req, nil, []int{http.StatusOK}); err != nil {
		return err
	}
	return nil
}

func (c *BenthosClient) UpdateStream(streamId string, stream string) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/streams/%s", c.Host, streamId),
		strings.NewReader(stream))
	if err != nil {
		return err
	}

	if err := c.doRequest(req, nil, []int{http.StatusOK}); err != nil {
		return err
	}
	return nil
}

func (c *BenthosClient) DeleteStream(streamId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/streams/%s", c.Host, streamId), nil)
	if err != nil {
		return err
	}

	if err := c.doRequest(req, nil, []int{http.StatusOK}); err != nil {
		return err
	}
	return nil
}
