package forward

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Forward struct {
	logger *zerolog.Logger
}

func NewForward(logger *zerolog.Logger) *Forward {
	return &Forward{
		logger: logger,
	}
}

func (f *Forward) ForwardPayload(targetURL string, payload map[string]any) error {
	p, err := json.Marshal(payload)
	if err != nil {
		f.logger.Error().Err(err).Msg("failed to marshal the json")
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		targetURL,
		bytes.NewBuffer(p),
	)

	if err != nil {
		f.logger.Error().Err(err).Msg("failed make new request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 7 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		f.logger.Error().Err(err).Msg("failed send request")
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		f.logger.Error().Err(err).Msg("got non 2xx status code")
		return errors.New("non 2xx status code")
	}

	resp.Body.Close()
	return nil

}
