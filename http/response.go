package http

import (
	"encoding/json"

	"github.com/pteropackages/soar/config"
)

func HandleItemResponse(buf []byte, cfg *config.Config) ([]byte, error) {
	var model struct {
		O string      `json:"object"`
		A interface{} `json:"attributes"`
	}
	if err := json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	var res []byte
	var err error

	if cfg.Http.ParseBody {
		if cfg.Http.ParseIndent {
			res, err = json.MarshalIndent(model.A, "", "  ")
		} else {
			res, err = json.Marshal(model.A)
		}

		return res, err
	}

	if cfg.Http.ParseIndent {
		res, err = json.MarshalIndent(model, "", "  ")
	} else {
		res, err = json.Marshal(model)
	}

	return res, err
}

func HandleDataResponse(buf []byte, cfg *config.Config) ([]byte, error) {
	var model struct {
		O string `json:"object"`
		D []struct {
			O string      `json:"object"`
			A interface{} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	var res []byte
	var err error

	if cfg.Http.ParseBody {
		inner := make([]interface{}, 0, len(model.D))
		for _, m := range model.D {
			inner = append(inner, m.A)
		}

		if cfg.Http.ParseIndent {
			res, err = json.MarshalIndent(inner, "", "  ")
		} else {
			res, err = json.Marshal(inner)
		}

		return res, err
	}

	if cfg.Http.ParseIndent {
		res, err = json.MarshalIndent(model, "", "  ")
	} else {
		res, err = json.Marshal(model)
	}

	return res, err
}
