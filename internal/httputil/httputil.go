package httputil

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type HTTPPostRequestFunc func(reqBody interface{}, requestRef *string) ([]byte, error)

func InitHttpClient(timeout time.Duration, maxIdleConn, maxIdleConnPerHost, maxConnPerHost int) *http.Client {
	certPool := x509.NewCertPool()
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            certPool,
				InsecureSkipVerify: true,
			},
			MaxIdleConns:        maxIdleConn,
			MaxIdleConnsPerHost: maxIdleConnPerHost,
			MaxConnsPerHost:     maxConnPerHost,
		},
	}
	return client
}

func NewHttpPostCall(client *http.Client, url string) HTTPPostRequestFunc {
	return func(reqBody interface{}, requestRef *string) ([]byte, error) {

		message, _ := json.Marshal(&reqBody)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(message))
		if err != nil {
			return nil, errors.Wrap(err, "Unable to New http Request")
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Aurora-Secret", "Internal")
		if requestRef != nil {
			req.Header.Add("RequestRef", *requestRef)
		} else {
			req.Header.Add("RequestRef", uuid.NewString())
		}

		res, err := client.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Unable to request %s", url))
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, errors.Wrap(err, "Unable to New http Request ioutil")
		}

		if res.StatusCode != http.StatusOK {
			return body, fmt.Errorf("%s", body)
		}

		return body, nil
	}
}

func DownloadFile(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad response from server: %d", resp.StatusCode)
	}
	fileData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	return fileData, nil
}
