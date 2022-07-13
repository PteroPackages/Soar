package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

type ErrorInfo struct {
	Status string `json:"string"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

type Error struct {
	msg  string
	info []*ErrorInfo
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Info() []*ErrorInfo {
	return e.info
}

func newError(err error, info []*ErrorInfo) *Error {
	if err == nil {
		return nil
	}

	return &Error{msg: err.Error(), info: info}
}

func (c *Client) Execute(req *http.Request) ([]byte, *Error) {
	c.log.Ignore().Info("request: %s %s", req.Method, req.URL.Path)
	start := time.Now()

	res, err := http.DefaultClient.Do(req)
	taken := time.Since(start).Microseconds() / 1000

	c.log.Debug("response: %d (%vms)", res.StatusCode, taken)
	length := res.Header.Get("Content-Length")
	c.log.Debug("content length: %s", length)

	if err != nil {
		return nil, newError(err, nil)
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
		return buf, newError(err, nil)

	case http.StatusNoContent:
		return nil, nil

	default:
		if length != "" {
			val, err := strconv.Atoi(length)
			if err != nil {
				val = 0
			}

			if val == 0 {
				return nil, newError(fmt.Errorf("unknown api error: %s", res.Status), nil)
			}

			defer res.Body.Close()
			buf, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, newError(err, nil)
			}

			var data struct {
				Errors []*ErrorInfo `json:"errors"`
			}
			if err = json.Unmarshal(buf, &data); err != nil {
				return nil, newError(err, nil)
			}

			return nil, newError(err, data.Errors)
		}

		return nil, newError(fmt.Errorf("unknown api error: %s", res.Status), nil)
	}
}
