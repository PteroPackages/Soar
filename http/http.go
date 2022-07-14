package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/logger"
)

type Client struct {
	http   *http.Client
	config *config.Config
	auth   *config.Auth
	log    *logger.Logger
}

func New(cfg *config.Config, auth *config.Auth, log *logger.Logger) *Client {
	return &Client{
		http:   &http.Client{},
		config: cfg,
		auth:   auth,
		log:    log,
	}
}

func (c *Client) Request(method, path string) *http.Request {
	var body bytes.Buffer
	req, _ := http.NewRequest(method, c.auth.URL+path, &body)

	req.Header.Set("User-Agent", "Soar Http Client")
	req.Header.Set("Authorization", "Bearer "+c.auth.Key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req
}

type errorInfo struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}

func (e *errorInfo) String() string {
	detail := e.Detail
	if detail == "" {
		detail = "<no details>"
	}

	return fmt.Sprintf("%s (%s): %s", e.Code, e.Status, detail)
}

func (c *Client) Execute(req *http.Request) ([]byte, error) {
	c.log.Ignore().Info("request: %s %s", req.Method, req.URL.Path)
	start := time.Now()

	res, err := http.DefaultClient.Do(req)
	taken := time.Since(start).Microseconds() / 1000
	c.log.Debug("response: %d (%vms)", res.StatusCode, taken)

	if err != nil {
		return nil, err
	}
	c.log.Ignore().Info("response: %d (%dms)", res.StatusCode, taken)

	switch res.StatusCode {
	case http.StatusOK:
		fallthrough

	case http.StatusCreated:
		fallthrough

	case http.StatusAccepted:
		defer res.Body.Close()
		buf, err := io.ReadAll(res.Body)
		return buf, err

	case http.StatusNoContent:
		return nil, nil

	default:
		defer res.Body.Close()
		buf, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("unknown api error: %s", res.Status)
		}

		var data struct {
			Errors []*errorInfo `json:"errors"`
		}
		if err = json.Unmarshal(buf, &data); err != nil {
			return nil, err
		}

		for _, e := range data.Errors {
			c.log.Error(e.String())
		}

		return nil, nil
	}
}
