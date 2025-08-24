package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/egress"
)

type httpClient struct {
	client *http.Client
	config *models.HttpClient
}

func NewHttpClient(config *models.HttpClient) (egress.HttpClientPorts, error) {
	client, err := initHttpClient(config)
	return &httpClient{
		client: client,
	}, err
}

func (h *httpClient) Execute(url, method string, reqPayload io.Reader, resPayload any) error {
	var (
		reqPayloadReader io.Reader
		err              error
	)

	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		reqPayloadReader, err = getRequestPayload(reqPayload)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, url, reqPayloadReader)
	if err != nil {
		return fmt.Errorf("http execute: failed to create request: %w", err)
	}
	req.Header.Set(constants.ContentType.String(), constants.Json.String())

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("http execute: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resPayload != nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("http execute: failed to read response body: %w", err)
		}

		if err := json.Unmarshal(bodyBytes, resPayload); err != nil {
			return fmt.Errorf("http execute: failed to unmarshal response: %w", err)
		}
	}
	return nil
}
