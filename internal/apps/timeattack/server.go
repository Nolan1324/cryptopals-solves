package timeattack

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

const TestEndpoint = "test"

// a TimeAttackServer holds configuration data for the HMAC verification severe that is vulnerable
// to a timing attack.
type TimeAttackServer struct {
	key             []byte
	compareDuration time.Duration
	server          *http.Server
	verboseLogs     bool
}

// NewServer creates a new HMAC verification server that is vulnerable to a timing attack.
// compareDuration specifies how long the insecure character comparison takes.
func NewServer(address string, key []byte, compareDuration time.Duration, verboseLogs bool) *TimeAttackServer {
	timeAttackServer := TimeAttackServer{key: key, compareDuration: compareDuration, verboseLogs: verboseLogs}

	mux := http.NewServeMux()
	mux.HandleFunc("/test", timeAttackServer.handler)

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}
	timeAttackServer.server = server

	return &timeAttackServer
}

func (s *TimeAttackServer) Sign(data []byte) []byte {
	mac := hmac.New(sha1.New, s.key)
	mac.Write([]byte(data)) // ignore error
	return mac.Sum(nil)
}

// Run starts the server and blocks.
// Returns ErrServerClosed if the server was closed.
func (s *TimeAttackServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *TimeAttackServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *TimeAttackServer) Close() error {
	return s.server.Close()
}

func (s *TimeAttackServer) handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	file := q.Get("file")
	signature := q.Get("signature")

	mac := hmac.New(sha1.New, s.key)
	mac.Write([]byte(file))
	expectedSignatureBytes := mac.Sum(nil)

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if s.verboseLogs {
		log.Printf("\nFile: %v\nGiven signature: %x\nExpected: %x\n", file, signatureBytes, expectedSignatureBytes)
	}

	if insecureCompare(signatureBytes, expectedSignatureBytes, s.compareDuration) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// insecureCompare compares two byte arrays in an insecure manner.
// compareDuration how long each byte comparison takes.
func insecureCompare(buf1 []byte, buf2 []byte, compareDuration time.Duration) bool {
	if len(buf1) != len(buf2) {
		return false
	}
	for i := range buf1 {
		if buf1[i] != buf2[i] {
			return false
		}
		time.Sleep(compareDuration)
	}
	return true
}
