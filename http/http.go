package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pteropackages/soar/config"
	"github.com/pteropackages/soar/logger"
	"github.com/spf13/pflag"
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

func Request(method, url string, body *bytes.Buffer) *http.Request {
	if body == nil {
		body = &bytes.Buffer{}
	}
	req, _ := http.NewRequest(method, url, body)

	req.Header.Set("User-Agent", "Soar Http Client")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req
}

func (c *Client) Request(method, path string, body *bytes.Buffer) *http.Request {
	if body == nil {
		body = &bytes.Buffer{}
	}
	req, _ := http.NewRequest(method, c.auth.URL+path, body)

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

func (c *Client) ExecuteWithFlags(req *http.Request, flags *pflag.FlagSet) ([]byte, error) {
	query := req.URL.Query()

	if val, _ := flags.GetInt("page"); val > 0 {
		query.Add("page", fmt.Sprint(val))
	}

	if val, _ := flags.GetInt("per-page"); val > 0 {
		if val > 100 {
			val = 100
		}
		query.Add("per_page", fmt.Sprint(val))
	}

	req.URL.RawQuery = query.Encode()
	return c.Execute(req)
}

func (c *Client) Execute(req *http.Request) ([]byte, error) {
	req.URL.Query()
	c.log.Ignore().Info("request: %s %s", req.Method, req.URL.Path)
	c.log.Debug("url: %s", req.URL.String())
	start := time.Now()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	taken := time.Since(start).Microseconds() / 1000
	c.log.Debug("response: %d (%vms)", res.StatusCode, taken)
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

		c.log.Debug("host: %s", res.Request.Host)
		c.log.Debug(string(buf))

		if strings.Contains(c.auth.URL, res.Request.Host) {
			var data struct {
				Errors []*errorInfo `json:"errors"`
			}
			if err = json.Unmarshal(buf, &data); err != nil {
				return nil, err
			}

			c.log.Error("received %d error(s):", len(data.Errors))
			for _, e := range data.Errors {
				c.log.Error(e.String())
			}
		} else {
			var data struct {
				Error string `json:"error"`
			}
			if err = json.Unmarshal(buf, &data); err != nil {
				return nil, err
			}

			c.log.Error("received an error:").Error(data.Error)
		}

		return nil, nil
	}
}
