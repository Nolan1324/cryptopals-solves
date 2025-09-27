package timeattack

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// WaitForServerStartWithTimeout waits for the HTTP server at endpoint to be reachable via GET request.
// It polls the endpoint every interval until it is reachable or the timeout is reached.
// If was reachable, returns nil. Otherwise, returns the last error from the server.
// This function is useful if the HTTP server was started in a seperate process.
func WaitForServerStartWithTimeout(timeout time.Duration, interval time.Duration, endpoint string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), timeout)
	defer cancelCtx()
	return WaitForServerStart(ctx, interval, endpoint)
}

// WaitForServerStart waits for the HTTP server at endpoint to be reachable via GET request.
// It polls the endpoint every interval until it is reachable or ctx signals a cancel.
// If was reachable, returns nil. Otherwise, returns the last error from the server.
// This function is useful if the HTTP server was started in a seperate process.
func WaitForServerStart(ctx context.Context, interval time.Duration, endpoint string) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	check := func() error {
		_, err := http.Get(endpoint)
		return err
	}

	lastErr := check()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout. Last server error: %w", lastErr)
		case <-ticker.C:
			lastErr = check()
			if lastErr == nil {
				return nil
			}
		}
	}
}

// DoRequest makes a request against the HMAC verification server vulnerable to a timing attack.
// endpoint specifies the endpoint to attack, file specifies the file data, and signature specifies
// the signature to check the file data against.
func DoRequest(endpoint string, file string, signature string) (time.Duration, bool, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return 0, false, fmt.Errorf("making request: %w", err)
	}
	q := req.URL.Query()
	q.Add("file", file)
	q.Add("signature", signature)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	startTime := time.Now()
	resp, err := client.Do(req)
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	if err != nil {
		return duration, false, err
	}

	return duration, resp.StatusCode == http.StatusOK, nil
}
