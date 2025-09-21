package timeattack

import (
	"net/http"
	"time"
)

// DoRequest makes a request against the HMAC verification server vulnerable to a timing attack.
// endpoint specifies the endpoint to attack, file specifies the file data, and signature specifies
// the signature to check the file data against.
func DoRequest(endpoint string, file string, signature string) (time.Duration, bool, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return 0, false, err
	}
	q := req.URL.Query()
	q.Add("file", file)
	q.Add("signature", signature)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	startTime := time.Now()
	resp, err := client.Do(req)
	endTime := time.Now()

	if err != nil {
		return 0, false, err
	}

	return endTime.Sub(startTime), resp.StatusCode == http.StatusOK, nil
}
