package http

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
)

// getClient creates an HTTP client with optional TLS configuration
func initHttpClient(cfg *models.HttpClient) (*http.Client, error) {
	client := &http.Client{
		Timeout: cfg.Timeout,
	}

	// If TLS is required, attach a configured Transport
	if cfg.ClientTLSRequired {
		tlsCfg, err := getHttpTransportCerts(cfg.CertPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS config from %s: %w", cfg.CertPath, err)
		}

		client.Transport = &http.Transport{
			TLSClientConfig:     tlsCfg,
			ForceAttemptHTTP2:   true,
			TLSHandshakeTimeout: 10 * time.Second,
		}
	}

	return client, nil
}

// getHTTPTransportCerts loads TLS certificates for secure HTTP transport
func getHttpTransportCerts(certPath string) (*tls.Config, error) {
	caCert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate at %s: %w", certPath, err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("no valid PEM certificates found in %s", certPath)
	}

	return &tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,            // ensure server cert is validated
		MinVersion:         tls.VersionTLS12, // enforce TLS 1.2+
	}, nil
}

// getRequestPayload marshals a body into JSON and returns it as an io.Reader
func getRequestPayload(body any) (io.Reader, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("getRequestPayload: failed to marshal body: %w", err)
	}

	// return io.NopCloser(io.Reader(bytes.NewReader(bodyBytes))), nil
	return bytes.NewReader(bodyBytes), nil
}
